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

	warName := fmt.Sprintf("The %s War for %s", attacker.Demonym, target)

	if attacker != nil && defender != nil && len(warName) != 0 {
		wars = append(wars, &war.War{Attacker: attacker, Defender: defender, Score: 0, Name: warName})
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

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/war", warHandler)
	mux.HandleFunc("/colonize", colonizeHandler)
	mux.HandleFunc("/", indexHandler)

	fileServer := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	http.ListenAndServe(":5000", mux)
}
