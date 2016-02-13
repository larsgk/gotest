package main

import (
	"encoding/json"
	// "github.com/larsgk/gotest/auth"
	"github.com/gorilla/websocket"
	"github.com/larsgk/gotest/comm"
	"image"
	"image/color"
	"image/png"
	"log"
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

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	m := image.NewNRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{32, 32}})

	// TODO: Reflect connection status or other here...
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			m.SetNRGBA(x, y, color.NRGBA{uint8(x * 8), uint8((x + y) * 8), uint8(y) * 8, 255})

		}
	}
	w.Header().Set("Content-Type", "image/png")
	if err := png.Encode(w, m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var upgrader = websocket.Upgrader{}

func handleWSData(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
