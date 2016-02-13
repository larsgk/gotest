package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	fmt.Printf("Hello, simple web service :)\n")

	r := mux.NewRouter()

	r.HandleFunc("/favicon.png", handleFavicon)

	r.HandleFunc("/now", handleNowEvent).
		Methods("POST")
	r.HandleFunc("/commports", handleListCommPortsEvent).
		Methods("POST")

	http.HandleFunc("/data", handleWSData)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", r)

	http.ListenAndServe(":3000", nil)
}
