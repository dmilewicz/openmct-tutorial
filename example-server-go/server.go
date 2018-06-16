package main

import (
	"net/http"
	"strings"
	"fmt"
	// "encoding/json"
	// "io/ioutil"
	"server"
	"craftsim"
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
	server.Hello()
	// craftsim.LoadCraftJSON("../dictionary.json")
	
	rserver := &server.RealtimeServer { Telem: make(chan interface{}) }
	hserver := &server.HistoryServer  { Telem: make(chan interface{}) }

	sp := craftsim.NewSpacecraft()
	go sp.RunSim()



	http.Handle("/", http.FileServer(http.Dir("../")))
	go http.HandleFunc("/realtime/", rserver.RunServer)
	go http.HandleFunc("/history/", hserver.RunServer)
	http.HandleFunc("/ping", sayHello)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

	fmt.Println(sp)
}