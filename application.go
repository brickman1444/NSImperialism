package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brickman1444/NSImperialism/dynamodbwrapper"
	"github.com/brickman1444/NSImperialism/nationstates_api"
	"github.com/brickman1444/NSImperialism/session"
	"github.com/brickman1444/NSImperialism/strategicmap"
	"github.com/brickman1444/NSImperialism/war"
	"github.com/joho/godotenv"
)

var globalWars = war.WarProviderDatabase{}
var globalResidentNations = strategicmap.ResidentsDatabase{}
var globalStrategicMap = strategicmap.StaticMap
var globalYear = 0
var globalSessionManager = session.NewSessionManager()

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

func indexHandler(w http.ResponseWriter, r *http.Request) {

	indexTemplate := template.Must(template.ParseFiles("index.html"))

	retrievedWars, err := globalWars.GetWars()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to retrieve wars", http.StatusInternalServerError)
		return
	}

	renderedMap, err := strategicmap.Render(globalStrategicMap, globalResidentNations, retrievedWars)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to render map", http.StatusInternalServerError)
		return
	}

	loggedInNation := getLoggedInNationFromCookie(r)

	page := &Page{"", nil, retrievedWars, renderedMap, globalYear, loggedInNation}

	indexTemplate.Execute(w, page)
}

type Page struct {
	Query          string
	Nation         *nationstates_api.Nation
	Wars           []war.War
	Map            strategicmap.RenderedMap
	Year           int
	LoggedInNation *nationstates_api.Nation
}

func searchHandler(w http.ResponseWriter, r *http.Request) {

	var indexTemplate = template.Must(template.ParseFiles("index.html"))

	url, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := url.Query()
	searchQuery := params.Get("q")

	nation, err := nationstates_api.GetNationData(searchQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	retrievedWars, err := globalWars.GetWars()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to retrieve wars", http.StatusInternalServerError)
		return
	}

	renderedMap, err := strategicmap.Render(globalStrategicMap, globalResidentNations, retrievedWars)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to render map", http.StatusInternalServerError)
		return
	}

	page := &Page{searchQuery, nation, retrievedWars, renderedMap, globalYear, nil}

	err = indexTemplate.Execute(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func warHandler(w http.ResponseWriter, r *http.Request) {

	attackerName := r.FormValue("attacker")
	attacker, err := nationstates_api.GetNationData(attackerName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if attacker == nil {
		http.Error(w, fmt.Sprintf("%s not found", attackerName), http.StatusNotFound)
		return
	}

	target := r.FormValue("target")

	defenderID, err := globalResidentNations.GetResident(target)
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

	retrievedWars, err := globalWars.GetWars()
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
		globalWars.PutWars([]war.War{newWar})
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func colonizeHandler(w http.ResponseWriter, r *http.Request) {

	colonizer, err := nationstates_api.GetNationData(r.FormValue("colonizer"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	target := r.FormValue("target")

	err = strategicmap.Colonize(&globalResidentNations, globalStrategicMap, *colonizer, target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func tickHandler(w http.ResponseWriter, r *http.Request) {

	err := tick(globalResidentNations, &globalWars, &globalYear)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func tick(residentNations strategicmap.ResidentsInterface, warsProvider war.WarProviderInterface, year *int) error {
	(*year)++

	retrievedWars, err := warsProvider.GetWars()
	if err != nil {
		return err
	}

	for warIndex := range retrievedWars {
		didFinish := retrievedWars[warIndex].Tick()
		if didFinish {
			residentNations.SetResident(retrievedWars[warIndex].TerritoryName, retrievedWars[warIndex].Advantage().Id)
		}
	}

	return warsProvider.PutWars(retrievedWars)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/uswds-2.10.0/img/flag.svg")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	nationName := r.FormValue("nation_name")
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

	if isVerified {
		w.Header().Add("verified", "yes")
	} else {
		w.Header().Add("verified", "no")
	}

	sessionIDBytes := sha1.Sum([]byte(nationName + strconv.Itoa(rand.Int())))
	sessionIDString := base64.URLEncoding.EncodeToString(sessionIDBytes[:]) // [:] converts slice to array

	cookieValue := nationName + SESSION_COOKIE_SEPARATOR + sessionIDString
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{Name: SESSION_COOKIE_NAME, Value: cookieValue, Expires: expire, HttpOnly: true}

	globalSessionManager.AddSession(nationName, sessionIDString, expire)

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DoAllEnvironmentVariablesExist() bool {

	variableNames := []string{
		"AWS_REGION",
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_KEY",
		"CELL_TABLE_NAME",
		"WAR_TABLE_NAME",
	}

	for _, variable := range variableNames {
		if os.Getenv(variable) == "" {
			return false
		}
	}

	return true
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Failed to load .env file:", err.Error())
	}

	dynamodbwrapper.Initialize()

	mux := http.NewServeMux()

	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/war", warHandler)
	mux.HandleFunc("/colonize", colonizeHandler)
	mux.HandleFunc("/tick", tickHandler)
	mux.HandleFunc("/", indexHandler)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.HandleFunc("/favicon.ico", faviconHandler)
	mux.HandleFunc("/login", loginHandler)

	http.ListenAndServe(":5000", mux)
}
