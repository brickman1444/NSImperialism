package main

import (
	"crypto/sha1"
	"encoding/base64"
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
var globalSessionManager = session.NewSessionManager()
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

	isValid := globalSessionManager.IsValidSession(nationName, sessionIDString, time.Now())
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

func indexHandler(w http.ResponseWriter, r *http.Request) {

	loggedInNation := getLoggedInNationFromCookie(r)

	mapIDs, err := dynamodbwrapper.GetAllMapIDs()
	if err != nil {
		http.Error(w, "Failed to get map IDs", http.StatusInternalServerError)
		return
	}

	page := &Page{[]war.War{}, strategicmap.RenderedMap{}, 0, loggedInNation, mapIDs, ""}

	renderPage(w, "index.html", page)
}

type Page struct {
	Wars           []war.War
	Map            strategicmap.RenderedMap
	Year           int
	LoggedInNation *nationstates_api.Nation
	Maps           []string
	MapID          string
}

func warHandler(w http.ResponseWriter, r *http.Request) {

	attacker := getLoggedInNationFromCookie(r)
	if attacker == nil {
		http.Error(w, "You must be logged in to attack", http.StatusBadRequest)
		return
	}

	routeVariables := mux.Vars(r)
	mapID := routeVariables["id"]

	databaseMap, err := globalMaps.GetMap(mapID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get map"), http.StatusInternalServerError)
		return
	}

	target := r.FormValue("target")

	defenderID, err := databaseMap.GetResident(target)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get defender for %s", target), http.StatusInternalServerError)
		return
	}

	if defenderID == "" {
		http.Error(w, fmt.Sprintf("No nation resides in %s", target), http.StatusBadRequest)
		return
	}

	if attacker.Id == defenderID {
		http.Error(w, fmt.Sprintf("You can't attack yourself"), http.StatusBadRequest)
		return
	}

	retrievedWars, err := war.GetWars(databaseMap, globalNationStatesProvider)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to retrieve wars", http.StatusInternalServerError)
		return
	}

	currentWar := war.FindOngoingWarAt(retrievedWars, target)
	if currentWar != nil {
		http.Error(w, fmt.Sprintf("There is already a war at %s", target), http.StatusBadRequest)
		return
	}

	warName := fmt.Sprintf("The %s War for %s", attacker.Demonym, target)

	defender, err := nationstates_api.GetNationData(defenderID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get defender data for %s", defenderID), http.StatusInternalServerError)
		return
	}

	if attacker != nil && len(warName) != 0 {
		newWar := war.NewWar(attacker, defender, warName, target)
		war.PutWars(&databaseMap, []war.War{newWar})
	}

	err = globalMaps.PutMap(databaseMap)
	if err != nil {
		http.Error(w, "Failed to save map", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/maps/"+mapID, http.StatusSeeOther)
}

func tickHandler(w http.ResponseWriter, r *http.Request) {

	routeVariables := mux.Vars(r)
	mapID := routeVariables["id"]

	databaseMap, err := globalMaps.GetMap(mapID)
	if err != nil {
		http.Error(w, "Failed to get map", http.StatusInternalServerError)
		return
	}

	err = tick(&databaseMap, globalNationStatesProvider)
	if err != nil {
		http.Error(w, "Failed to tick map", http.StatusInternalServerError)
		return
	}

	err = globalMaps.PutMap(databaseMap)
	if err != nil {
		http.Error(w, "Failed to save map", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/maps/"+mapID, http.StatusSeeOther)
}

func tick(residentNations *databasemap.DatabaseMap, nationStatesProvider nationstates_api.NationStatesProvider) error {

	residentNations.Year++

	retrievedWars, err := war.GetWars(*residentNations, nationStatesProvider)
	if err != nil {
		return err
	}

	for warIndex := range retrievedWars {
		didFinish := retrievedWars[warIndex].Tick()
		if didFinish {
			residentNations.SetResident(retrievedWars[warIndex].TerritoryName, retrievedWars[warIndex].Advantage().Id)
		}
	}

	war.PutWars(residentNations, retrievedWars)

	return nil
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/uswds-2.10.0/img/flag.svg")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	nationName := nationstates_api.GetCanonicalName(r.FormValue("nation_name"))
	verificationCode := r.FormValue("verification_code")
	if nationName == "" || verificationCode == "" {
		http.Error(w, "Invalid request to login. Try again.", http.StatusBadRequest)
		return
	}

	isVerified, err := nationstates_api.IsCorrectVerificationCode(nationName, verificationCode)
	if err != nil {
		log.Println("Failed to verify nation", nationName, err.Error())
		http.Error(w, "Failed to verify nation", http.StatusInternalServerError)
		return
	}

	log.Println(nationName, "verified:", strconv.FormatBool(isVerified))

	sessionIDBytes := sha1.Sum([]byte(nationName + strconv.Itoa(rand.Int())))
	sessionIDString := base64.URLEncoding.EncodeToString(sessionIDBytes[:]) // [:] converts slice to array

	cookieValue := nationName + SESSION_COOKIE_SEPARATOR + sessionIDString
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{Name: SESSION_COOKIE_NAME, Value: cookieValue, Expires: expire, HttpOnly: true}

	globalSessionManager.AddSession(nationName, sessionIDString, expire)

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getMapHandler(w http.ResponseWriter, r *http.Request) {

	routeVariables := mux.Vars(r)
	mapID := routeVariables["id"]

	databaseMap, err := globalMaps.GetMap(mapID)
	if err != nil {
		http.Error(w, "Failed to retrieve map", http.StatusInternalServerError)
		return
	}

	retrievedWars, err := war.GetWars(databaseMap, globalNationStatesProvider)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to retrieve wars", http.StatusInternalServerError)
		return
	}

	renderedMap, err := strategicmap.Render(globalStrategicMap, databaseMap, retrievedWars)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to render map", http.StatusInternalServerError)
		return
	}

	loggedInNation := getLoggedInNationFromCookie(r)

	page := &Page{retrievedWars, renderedMap, databaseMap.Year, loggedInNation, []string{}, databaseMap.ID}

	renderPage(w, "map.html", page)
}

func postMapHandler(w http.ResponseWriter, r *http.Request) {
	participatingNationNames := strings.Split(r.FormValue("participating_nations"), ",")
	if len(participatingNationNames) == 0 {
		http.Error(w, "List of participating nations was empty", http.StatusBadRequest)
		return
	}

	for _, nationName := range participatingNationNames {
		nation, err := nationstates_api.GetNationData(nationName)
		if nation == nil || err != nil {
			http.Error(w, "Could not find nation "+nationName, http.StatusBadRequest)
			return
		}
	}

	databaseMap, err := strategicmap.MakeNewRandomMap(globalStrategicMap, participatingNationNames)
	if err != nil {
		http.Error(w, "Failed to create map", http.StatusInternalServerError)
		return
	}

	err = dynamodbwrapper.PutMap(databaseMap)
	if err != nil {
		http.Error(w, "Failed to save map", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/maps/"+databaseMap.ID, http.StatusSeeOther)
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
	mux.HandleFunc("/maps/{id}", getMapHandler).Methods("GET")
	mux.HandleFunc("/maps", postMapHandler).Methods("POST")

	http.ListenAndServe(":5000", mux)
}
