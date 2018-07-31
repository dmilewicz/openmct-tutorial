package main

import (

	// "fmt"
	// "encoding/json"
	// "io/ioutil"
	"craftsim"
	"server"
)

func main() {

	sim := craftsim.NewSim()

	server.NewServer(8080, sim.RealtimeData, sim.HistoryRequest, sim.HistoryData).RunServer()

}
