package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/sub", Sub)
	router.HandleFunc("/todos/{todoId}", TodoShow)
	log.Fatal(http.ListenAndServe(":8080", router))

}

func Index(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello, %q", html.EscapeString(request.URL.Path))
}

func Sub(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}