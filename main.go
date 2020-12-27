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

func searchHandler(w http.ResponseWriter, r *http.Request) {

	var indexTemplate = template.Must(template.ParseFiles("index.html"))

	url, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := url.Query()
	searchQuery := params.Get("q")

	_, err = GetStandardData(searchQuery)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf(searchQuery)

	indexTemplate.Execute(w, nil)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":3000", mux)
}
