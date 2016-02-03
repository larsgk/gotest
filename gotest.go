package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Printf("Hello, simple web service :)\n")

	http.HandleFunc("/now", handleNowEvent)
	http.HandleFunc("/commports", handleListCommPortsEvent)
	http.ListenAndServe(":3000", nil)
}
