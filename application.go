package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brickman1444/NSImperialism/nationstates_api"
	"github.com/brickman1444/NSImperialism/strategicmap"
	"github.com/brickman1444/NSImperialism/war"
	"github.com/joho/godotenv"
)

const CELL_TABLE_NAME string = "nsimperialism-cell"

var globalWars []*war.War = []*war.War{}
var globalResidentNations = strategicmap.NewResidentsSimpleMap()
var globalStrategicMap = strategicmap.StaticMap
var globalYear = 0

func indexHandler(w http.ResponseWriter, r *http.Request) {

	indexTemplate := template.Must(template.ParseFiles("index.html"))

	page := &Page{"", nil, globalWars, strategicmap.Render(globalStrategicMap, globalResidentNations, globalWars), globalYear}

	indexTemplate.Execute(w, page)
}

type Page struct {
	Query  string
	Nation *nationstates_api.Nation
	Wars   []*war.War
	Map    strategicmap.RenderedMap
	Year   int
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

	page := &Page{searchQuery, nation, globalWars, strategicmap.Render(globalStrategicMap, globalResidentNations, globalWars), globalYear}

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

	defenderID := globalResidentNations.GetResident(target)
	if defenderID == "" {
		http.Error(w, fmt.Sprintf("No nation resides in %s", target), http.StatusBadRequest)
		return
	}

	if attacker.Id == defenderID {
		http.Error(w, fmt.Sprintf("You can't attack yourself"), http.StatusBadRequest)
		return
	}

	currentWar := war.FindOngoingWarAt(globalWars, target)
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
		globalWars = append(globalWars, &newWar)
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

	tick(globalResidentNations, globalWars, &globalYear)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func tick(residentNations strategicmap.ResidentsInterface, wars []*war.War, year *int) {
	(*year)++

	for _, war := range wars {
		didFinish := war.Tick()
		if didFinish {
			residentNations.SetResident(war.TerritoryName, war.Advantage().Id)
		}
	}
}

type DatabaseCell struct {
	ID       string
	Resident string
}

func (cell DatabaseCell) ToString() string {
	return fmt.Sprintf("ID: %v Resident: %v", cell.ID, cell.Resident)
}

func initializeDatabase() {
	databaseContext := context.TODO()

	awsConfig, err := config.LoadDefaultConfig(databaseContext)
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	dynamodbClient := dynamodb.NewFromConfig(awsConfig)

	getItemsResponse, err := dynamodbClient.Scan(databaseContext, &dynamodb.ScanInput{
		TableName: aws.String(CELL_TABLE_NAME),
	})

	if err != nil {
		log.Fatalf("failed to get items, %v", err)
	}

	records := []DatabaseCell{}
	err = attributevalue.UnmarshalListOfMaps(getItemsResponse.Items, &records)
	if err != nil {
		log.Println("failed to unmarshal Items, %w", err)
	}
	for _, record := range records {
		log.Println(record.ToString())
	}

	itemToPut := DatabaseCell{
		ID:       "A",
		Resident: "testlandia",
	}
	itemToPutMap, err := attributevalue.MarshalMap(itemToPut)
	if err != nil {
		log.Fatalf("failed to marshal item, %v", err)
	}

	_, err = dynamodbClient.PutItem(databaseContext, &dynamodb.PutItemInput{
		TableName: aws.String(CELL_TABLE_NAME),
		Item:      itemToPutMap,
	})
	if err != nil {
		log.Fatalf("failed to put item, %v", err)
	}

	getItemOutput, err := dynamodbClient.GetItem(databaseContext, &dynamodb.GetItemInput{
		TableName: aws.String(CELL_TABLE_NAME),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{
				Value: itemToPut.ID,
			},
		},
	})

	updatedItem := DatabaseCell{}
	err = attributevalue.UnmarshalMap(getItemOutput.Item, &updatedItem)
	if err != nil {
		log.Fatalf("failed to unmarshal item, %v", err)
	}

	log.Println(updatedItem.ToString())
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Failed to load .env file:", err.Error())
	}

	initializeDatabase()

	mux := http.NewServeMux()

	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/war", warHandler)
	mux.HandleFunc("/colonize", colonizeHandler)
	mux.HandleFunc("/tick", tickHandler)
	mux.HandleFunc("/", indexHandler)

	fileServer := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	http.ListenAndServe(":5000", mux)
}
