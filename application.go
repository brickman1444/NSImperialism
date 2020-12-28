package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

	var indexTemplate = template.Must(template.ParseFiles("index.html"))
	indexTemplate.Execute(w, nil)
}

type Page struct {
	Query  string
	Nation *Nation
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

	nation, err := GetNationData(searchQuery)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	page := &Page{searchQuery, nation}

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
	http.ListenAndServe(":5000", mux)
}
