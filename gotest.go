package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/favicon.png", handleFavicon)

	r.HandleFunc("/now", handleNowEvent).
		Methods("POST")
	r.HandleFunc("/commports", handleListCommPortsEvent).
		Methods("POST")

	http.HandleFunc("/data", handleWSData)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", r)

	port := "3000"

	fmt.Printf("This is a simple web service serving on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
