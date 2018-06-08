package main

import (
	"net/http"
	"strings"
	"fmt"
	"server"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))
}

func printRequest(w http.ResponseWriter, r *http.Request) {
	http_request := strings.TrimPrefix(r.URL.Path, "/")
	fmt.Println(http_request)
	w.Write([]byte("Requested: " + http_request))
}


func main() {

	http.Handle("/", http.FileServer(http.Dir("../")))
	http.HandleFunc("/realtime/", printRequest)
	http.HandleFunc("/ping", sayHello)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}