package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/brickman1444/NSImperialism/grid"
	"github.com/brickman1444/NSImperialism/nationstates_api"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

	var indexTemplate = template.Must(template.ParseFiles("index.html"))
	indexTemplate.Execute(w, nil)
}

type Page struct {
	Query       string
	Nation      *nationstates_api.Nation
	Belligerent *nationstates_api.Nation
	ThirdParty  *nationstates_api.Nation
	Grid        *grid.Grid
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

	grid := grid.Get()

	page := &Page{searchQuery, nation, belligerent, thirdParty, grid}

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
