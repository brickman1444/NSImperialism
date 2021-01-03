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
	"github.com/brickman1444/NSImperialism/grid"
	"github.com/brickman1444/NSImperialism/nationstates_api"
	"github.com/brickman1444/NSImperialism/war"
	"github.com/joho/godotenv"
)

const CELL_TABLE_NAME string = "nsimperialism-cell"

var globalWars []*war.War = []*war.War{}
var globalGrid *grid.Grid = &grid.Grid{}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	indexTemplate := template.Must(template.ParseFiles("index.html"))

	page := &Page{"", nil, globalGrid.Render(globalWars), globalWars}

	indexTemplate.Execute(w, page)
}

type Page struct {
	Query  string
	Nation *nationstates_api.Nation
	Grid   *grid.RenderedGrid
	Wars   []*war.War
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

	page := &Page{searchQuery, nation, globalGrid.Render(globalWars), globalWars}

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
	targetRowIndex, targetColumnIndex, err := globalGrid.GetCoordinates(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defender := globalGrid.Rows[targetRowIndex].Cells[targetColumnIndex].ResidentNation
	if defender == nil {
		http.Error(w, fmt.Sprintf("No nation resides in %s", target), http.StatusBadRequest)
		return
	}

	if attacker.Id == defender.Id {
		http.Error(w, fmt.Sprintf("You can't attack yourself"), http.StatusBadRequest)
		return
	}

	currentWar := war.FindOngoingWarAt(globalWars, targetRowIndex, targetColumnIndex)
	if currentWar != nil {
		http.Error(w, fmt.Sprintf("There is already a war at %s", target), http.StatusBadRequest)
		return
	}

	warName := fmt.Sprintf("The %s War for %s", attacker.Demonym, target)

	if attacker != nil && defender != nil && len(warName) != 0 {
		newWar := war.NewWar(attacker, defender, warName, targetRowIndex, targetColumnIndex)
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
	if len(target) > 2 {
		target = target[0:2]
	}

	err = globalGrid.Colonize(*colonizer, target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func tickHandler(w http.ResponseWriter, r *http.Request) {

	tick(globalGrid, globalWars)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func tick(grid *grid.Grid, wars []*war.War) {
	grid.Year++

	for _, war := range wars {
		didFinish := war.Tick()
		if didFinish {
			grid.Rows[war.TargetRowIndex].Cells[war.TargetColumnIndex].ResidentNation = war.Advantage()
		}
	}
}

type DatabaseCell struct {
	ID string
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
		log.Println("Record: ID:", record.ID)
	}
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
