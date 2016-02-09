package main

import (
	"encoding/json"
	// "github.com/larsgk/gotest/auth"
	"github.com/larsgk/gotest/comm"
	"net/http"
	"time"
)

type DaTimeEvent struct {
	Type string
	Time int64
}

type DaListEvent struct {
	Type string
	Data []string
}

func handleNowEvent(w http.ResponseWriter, r *http.Request) {
	daEvent := DaTimeEvent{"Now", time.Now().UnixNano() / 1000000}

	sendJsonEvent(w, daEvent)
}

func handleListCommPortsEvent(w http.ResponseWriter, r *http.Request) {
	daEvent := DaListEvent{Type: "CommPorts"}

	ports, _ := comm.GetSerialPortList()

	for _, port := range ports {
		daEvent.Data = append(daEvent.Data, port.Name)
	}

	sendJsonEvent(w, daEvent)
}

func sendJsonEvent(w http.ResponseWriter, daEvent interface{}) {
	js, err := json.Marshal(daEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "private, no-cache")
	w.Write(js)
}
