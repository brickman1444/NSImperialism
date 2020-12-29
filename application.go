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

var wars []*war.War = []*war.War{}
var globalGrid *grid.Grid = &grid.Grid{}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	indexTemplate := template.Must(template.ParseFiles("index.html"))

	page := &Page{"", nil, globalGrid.Render(), wars}

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

	page := &Page{searchQuery, nation, globalGrid.Render(), wars}

	err = indexTemplate.Execute(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func warHandler(w http.ResponseWriter, r *http.Request) {

	var indexTemplate = template.Must(template.ParseFiles("index.html"))

	/*err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}*/

	defender, err := nationstates_api.GetNationData(r.FormValue("defender"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attacker, err := nationstates_api.GetNationData(r.FormValue("attacker"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	target := r.FormValue("target")
	if len(target) > 2 {
		target = target[0:2]
	}

	warName := fmt.Sprintf("The %s War for %s", attacker.Demonym, target)

	if attacker != nil && defender != nil && len(warName) != 0 {
		wars = append(wars, &war.War{Attacker: attacker, Defender: defender, Score: 0, Name: warName})
	}

	page := &Page{"", nil, globalGrid.Render(), wars}

	err = indexTemplate.Execute(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {

	mux := http.NewServeMux()

	mechalus, err := nationstates_api.GetNationData("the_mechalus")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	testlandia, err := nationstates_api.GetNationData("testlandia")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	logistics, err := nationstates_api.GetNationData("nationstates_department_of_logistics")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	globalGrid.Rows[1].Cells[1].ResidentNation = mechalus
	globalGrid.Rows[2].Cells[1].ResidentNation = mechalus
	globalGrid.Rows[2].Cells[1].AttackerNation = testlandia
	globalGrid.Rows[3].Cells[2].ResidentNation = logistics
	globalGrid.Rows[3].Cells[3].ResidentNation = testlandia

	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/war", warHandler)
	mux.HandleFunc("/", indexHandler)

	fileServer := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	http.ListenAndServe(":5000", mux)
}
