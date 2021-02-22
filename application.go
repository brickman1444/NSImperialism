package main

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brickman1444/NSImperialism/databasemap"
	"github.com/brickman1444/NSImperialism/dynamodbwrapper"
	"github.com/brickman1444/NSImperialism/nationstates_api"
	"github.com/brickman1444/NSImperialism/session"
	"github.com/brickman1444/NSImperialism/strategicmap"
	"github.com/brickman1444/NSImperialism/war"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var globalMaps = strategicmap.MapsDatabase{}
var globalStrategicMap = strategicmap.StaticMap
var globalSessionManager = session.SessionManagerDatabase{}
var globalNationStatesProvider = nationstates_api.NationStatesProviderAPI{}

const SESSION_COOKIE_NAME = "SessionID"
const SESSION_COOKIE_SEPARATOR = ":"

func getLoggedInNationFromCookie(r *http.Request) *nationstates_api.Nation {
	sessionCookie, err := r.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		return nil // Cookie returns ErrNoCookie if the cookie isn't found
	}

	tokens := strings.Split(sessionCookie.Value, SESSION_COOKIE_SEPARATOR)
	if len(tokens) != 2 {
		return nil
	}

	nationName := tokens[0]
	sessionIDString := tokens[1]

	isValid, err := globalSessionManager.IsValidSession(nationName, sessionIDString, time.Now())
	if err != nil {
		return nil
	}

	if !isValid {
		return nil
	}

	nation, err := nationstates_api.GetNationData(nationName)
	if err != nil {
		return nil
	}

	return nation
}

func renderPage(w http.ResponseWriter, bodyTemplateFileName string, data interface{}) {
	bodyTemplate, err := template.ParseFiles(bodyTemplateFileName)
	if err != nil {
		http.Error(w, "Failed parse HTML body", http.StatusInternalServerError)
		return
	}
	headerTemplate, err := template.ParseFiles("header.html")
	if err != nil {
		http.Error(w, "Failed parse HTML header", http.StatusInternalServerError)
		return
	}
	footerTemplate, err := template.ParseFiles("footer.html")
	if err != nil {
		http.Error(w, "Failed parse HTML footer", http.StatusInternalServerError)
		return
	}

	err = headerTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed render HTML header", http.StatusInternalServerError)
		return
	}

	err = bodyTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed render HTML body", http.StatusInternalServerError)
		return
	}

	err = footerTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed render HTML footer", http.StatusInternalServerError)
		return
	}
}

func getParticipatingNations(databaseMap databasemap.DatabaseMap) ([]nationstates_api.Nation, error) {
	uniqueParticipantNationIDs := []string{}
	for _, cell := range databaseMap.Cells {
		if !contains(uniqueParticipantNationIDs, cell.Resident) {
			uniqueParticipantNationIDs = append(uniqueParticipantNationIDs, cell.Resident)
		}
	}

	nations := []nationstates_api.Nation{}
	for _, nationID := range uniqueParticipantNationIDs {

		nation, err := nationstates_api.GetNationData(nationID)
		if err != nil {
			return []nationstates_api.Nation{}, err
		}

		nations = append(nations, *nation)
	}

	return nations, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	loggedInNation := getLoggedInNationFromCookie(r)

	maps, err := dynamodbwrapper.GetAllMaps()
	if err != nil {
		ErrorHandler(w, r, "Failed to get map IDs")
		return
	}

	mapLinkDatas := []MapLinkData{}
	for _, databaseMap := range maps {

		participatingNations, err := getParticipatingNations(databaseMap)
		if err != nil {
			ErrorHandler(w, r, "Failed to get map participants")
			return
		}

		mapLinkDatas = append(mapLinkDatas, MapLinkData{
			MapID:                databaseMap.ID,
			ParticipatingNations: participatingNations,
		})
	}

	page := &Page{LoggedInNation: loggedInNation, Maps: mapLinkDatas}

	renderPage(w, "index.html", page)
}

type MapLinkData struct {
	MapID                string
	ParticipatingNations []nationstates_api.Nation
}

type Page struct {
	Wars           []war.RenderedWar
	Map            strategicmap.RenderedMap
	Year           int
	LoggedInNation *nationstates_api.Nation
	Maps           []MapLinkData
	MapID          string
	Error          string
}

func warHandler(w http.ResponseWriter, r *http.Request) {

	attacker := getLoggedInNationFromCookie(r)
	if attacker == nil {
		ErrorHandler(w, r, "You must be logged in to attack")
		return
	}

	routeVariables := mux.Vars(r)
	mapID := routeVariables["id"]

	databaseMap, err := globalMaps.GetMap(mapID)
	if err != nil {
		ErrorHandler(w, r, "Failed to get map")
		return
	}

	target := r.FormValue("target")

	defenderID, err := databaseMap.GetResident(target)
	if err != nil {
		ErrorHandler(w, r, fmt.Sprintf("Failed to get defender for %s", target))
		return
	}

	if defenderID == "" {
		ErrorHandler(w, r, fmt.Sprintf("No nation resides in %s", target))
		return
	}

	if attacker.Id == defenderID {
		ErrorHandler(w, r, fmt.Sprintf("You can't attack yourself"))
		return
	}

	retrievedWars := databaseMap.GetWars()

	currentWar := war.FindOngoingWarAt(retrievedWars, target)
	if currentWar != nil {
		ErrorHandler(w, r, fmt.Sprintf("There is already a war at %s", target))
		return
	}

	occasion := r.FormValue("occasion")
	if len(occasion) == 0 {
		ErrorHandler(w, r, "You didn't choose a valid occasion for war")
		return
	}

	warName := fmt.Sprintf("The %s %s %s", attacker.Demonym, occasion, target)

	defender, err := nationstates_api.GetNationData(defenderID)
	if err != nil {
		ErrorHandler(w, r, fmt.Sprintf("Failed to get defender data for %s", defenderID))
		return
	}

	if attacker != nil && len(warName) != 0 {
		newWar := databasemap.NewWar(attacker.Id, defender.Id, warName, target)
		databaseMap.PutWars([]databasemap.DatabaseWar{newWar})
	}

	err = globalMaps.PutMap(databaseMap)
	if err != nil {
		ErrorHandler(w, r, "Failed to save map")
		return
	}

	http.Redirect(w, r, "/maps/"+mapID, http.StatusSeeOther)
}

func tickHandler(w http.ResponseWriter, r *http.Request) {

	routeVariables := mux.Vars(r)
	mapID := routeVariables["id"]

	databaseMap, err := globalMaps.GetMap(mapID)
	if err != nil {
		ErrorHandler(w, r, "Failed to get map")
		return
	}

	err = tick(&databaseMap, globalNationStatesProvider)
	if err != nil {
		ErrorHandler(w, r, "Failed to tick map")
		return
	}

	err = globalMaps.PutMap(databaseMap)
	if err != nil {
		ErrorHandler(w, r, "Failed to save map")
		return
	}

	http.Redirect(w, r, "/maps/"+mapID, http.StatusSeeOther)
}

func tick(residentNations *databasemap.DatabaseMap, nationStatesProvider nationstates_api.NationStatesProvider) error {

	residentNations.Year++

	databaseWars := residentNations.GetWars()

	for warIndex := range databaseWars {
		didFinish, err := war.Tick(&databaseWars[warIndex], nationStatesProvider)
		if err != nil {
			return err
		}

		if didFinish {

			advantageID := war.WarAdvantage(databaseWars[warIndex])

			if advantageID == nil {
				return errors.New("Nil war winner ID")
			}

			residentNations.SetResident(databaseWars[warIndex].TerritoryName, *advantageID)
		}
	}

	residentNations.PutWars(databaseWars)

	return nil
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/uswds-2.10.0/img/flag.svg")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	nationName := nationstates_api.GetCanonicalName(r.FormValue("nation_name"))
	verificationCode := r.FormValue("verification_code")
	if nationName == "" || verificationCode == "" {
		ErrorHandler(w, r, "Invalid request to login. Try again.")
		return
	}

	isVerified, err := nationstates_api.IsCorrectVerificationCode(nationName, verificationCode)
	if err != nil {
		ErrorHandler(w, r, "Failed to verify nation "+nationName)
		return
	}

	log.Println(nationName, "verified:", strconv.FormatBool(isVerified))

	sessionIDBytes := sha1.Sum([]byte(nationName + strconv.Itoa(rand.Int())))
	sessionIDString := base64.URLEncoding.EncodeToString(sessionIDBytes[:]) // [:] converts slice to array

	cookieValue := nationName + SESSION_COOKIE_SEPARATOR + sessionIDString
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{Name: SESSION_COOKIE_NAME, Value: cookieValue, HttpOnly: true}

	globalSessionManager.AddSession(nationName, sessionIDString, expire)

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	loggedInNation := getLoggedInNationFromCookie(r)
	if loggedInNation == nil {
		ErrorHandler(w, r, "You must be logged in to log out.")
		return
	}

	globalSessionManager.RemoveSession(loggedInNation.Id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getMapHandler(w http.ResponseWriter, r *http.Request) {

	routeVariables := mux.Vars(r)
	mapID := routeVariables["id"]

	databaseMap, err := globalMaps.GetMap(mapID)
	if err != nil {
		ErrorHandler(w, r, "Failed to retrieve map")
		return
	}

	retrievedWars := databaseMap.GetWars()

	renderedMap, err := strategicmap.Render(globalStrategicMap, databaseMap, retrievedWars, globalNationStatesProvider)
	if err != nil {
		ErrorHandler(w, r, "Failed to render map")
		return
	}

	loggedInNation := getLoggedInNationFromCookie(r)

	renderedWars, err := war.RenderWars(retrievedWars, globalNationStatesProvider)
	if err != nil {
		ErrorHandler(w, r, "Failed to render wars")
		return
	}

	page := &Page{Wars: renderedWars, Map: renderedMap, Year: databaseMap.Year, LoggedInNation: loggedInNation, MapID: databaseMap.ID}

	renderPage(w, "map.html", page)
}

func contains(list []string, valueToLookFor string) bool {
	for _, element := range list {
		if element == valueToLookFor {
			return true
		}
	}

	return false
}

func postMapHandler(w http.ResponseWriter, r *http.Request) {
	participatingNationNames := strings.Split(r.FormValue("participating_nations"), ",")
	if len(participatingNationNames) == 0 {
		ErrorHandler(w, r, "List of participating nations was empty.")
		return
	}

	participatingNationNamesCanonical := []string{}
	for _, nationName := range participatingNationNames {
		participatingNationNamesCanonical = append(participatingNationNamesCanonical, nationstates_api.GetCanonicalName(nationName))
	}

	loggedInNation := getLoggedInNationFromCookie(r)
	if loggedInNation == nil {
		ErrorHandler(w, r, "You must be logged in to create a map.")
		return
	}

	if !contains(participatingNationNamesCanonical, loggedInNation.Id) {
		ErrorHandler(w, r, "You must participate in a map to create it. Add your nation to the list of participants and try again.")
		return
	}

	for _, nationName := range participatingNationNamesCanonical {
		nation, err := nationstates_api.GetNationData(nationName)
		if nation == nil || err != nil {
			ErrorHandler(w, r, "Could not find nation '"+nationName+"'. Check for typing or spelling errors such as extra spaces and try again.")
			return
		}
	}

	databaseMap, err := strategicmap.MakeNewRandomMap(globalStrategicMap, participatingNationNamesCanonical)
	if err != nil {
		ErrorHandler(w, r, err.Error())
		return
	}

	err = dynamodbwrapper.PutMap(databaseMap)
	if err != nil {
		ErrorHandler(w, r, "Failed to save map. Try again later.")
		return
	}

	http.Redirect(w, r, "/maps/"+databaseMap.ID, http.StatusSeeOther)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	ErrorHandler(w, r, "Page not found.")
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	ErrorHandler(w, r, "What you're trying to do is not supported.")
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, message string) {
	page := Page{LoggedInNation: getLoggedInNationFromCookie(r), Error: message}
	renderPage(w, "error.html", page)
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Failed to load .env file:", err.Error())
	}

	dynamodbwrapper.Initialize()

	rand.Seed(time.Now().UnixNano())

	mux := mux.NewRouter()

	mux.HandleFunc("/war/{id}", warHandler).Methods("POST")
	mux.HandleFunc("/tick/{id}", tickHandler).Methods("POST")
	mux.HandleFunc("/", indexHandler).Methods("GET")
	mux.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))).Methods("GET")
	mux.HandleFunc("/favicon.ico", faviconHandler).Methods("GET")
	mux.HandleFunc("/login", loginHandler).Methods("POST")
	mux.HandleFunc("/logout", logoutHandler).Methods("POST")
	mux.HandleFunc("/maps/{id}", getMapHandler).Methods("GET")
	mux.HandleFunc("/maps", postMapHandler).Methods("POST")

	mux.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	mux.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedHandler)

	http.ListenAndServe(":5000", mux)
}
