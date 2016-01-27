package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type EkEvent struct {
	Type string
	Time int64
}

func handlerNowEvent(w http.ResponseWriter, r *http.Request) {
	testEvent := EkEvent{"Now", time.Now().Unix()}

	js, err := json.Marshal(testEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	fmt.Printf("Hello, json\n")
	http.HandleFunc("/now", handlerNowEvent)
	http.ListenAndServe(":3000", nil)
}
