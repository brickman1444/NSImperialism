package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/brickman1444/NSImperialism/grid"
	"github.com/brickman1444/NSImperialism/nationstates_api"
	"github.com/brickman1444/NSImperialism/war"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

	var indexTemplate = template.Must(template.ParseFiles("index.html"))
	indexTemplate.Execute(w, nil)
}

type Page struct {
	Query       string
	Nation      *nationstates_api.Nation
	Belligerent *nationstates_api.Nation
	Grid        *grid.RenderedGrid
	War         *war.War
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
		fmt.Println(err.Error())
		return
	}

	belligerent, err := nationstates_api.GetNationData("testlandia")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	thirdParty, err := nationstates_api.GetNationData("nationstates_department_of_logistics")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	grid := &grid.Grid{}

	grid.Rows[1].Cells[1].ResidentNation = nation
	grid.Rows[2].Cells[1].ResidentNation = nation
	grid.Rows[2].Cells[1].AttackerNation = belligerent
	grid.Rows[3].Cells[2].ResidentNation = thirdParty
	grid.Rows[3].Cells[3].ResidentNation = belligerent

	renderedGrid := grid.Render()

	war := &war.War{Attacker: belligerent, Defender: nation, Score: -10, Name: "The Testlandian Conquest of A2"}

	page := &Page{searchQuery, nation, belligerent, renderedGrid, war}

	err = indexTemplate.Execute(w, page)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/", indexHandler)

	fileServer := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	http.ListenAndServe(":5000", mux)
}
