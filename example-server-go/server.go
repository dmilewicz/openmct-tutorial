package main

import (
	"net/http"
	"strings"
	// "fmt"
	// "encoding/json"
	// "io/ioutil"
	"craftsim"
	"server"

	"golang.org/x/net/websocket"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))
}

// func printRequest(w http.ResponseWriter, r *http.Request) {
// 	http_request := strings.TrimPrefix(r.URL.Path, "/")
// 	fmt.Println(http_request)

// 	dictionary, err := ioutil.ReadFile("../dictionary.json")
// 	if err != nil {
// 		fmt.Println("error:", err)
// 	}
// 	w.Write([]byte(dictionary))
// }

func main() {

	// hserver := &server.HistoryServer  { Telem: make(chan interface{}) }

	sim := craftsim.NewSim()

	hserver := &server.HistoryServer{sim.HistoryRequest, sim.HistoryData}
	rserver := server.NewRealtimeServer(sim.RealtimeData)

	go sim.RunSim()

	// go
	go http.Handle("/", http.FileServer(http.Dir("../")))
	go http.Handle("/realtime/", websocket.Handler(rserver.RealtimeSocket))
	go http.HandleFunc("/history/", hserver.RunServer)

	http.HandleFunc("/ping", sayHello)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

	// fmt.Println(sp)
}
