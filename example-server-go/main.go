package main

import (
	"server"
)

func main() {

	// sim := craftsim.NewSim()
	hs := server.NewHistoryStore()

	server.NewServer(8080, hs.HistoryRequest, hs.HistoryData, hs.DataIn).RunServer()

}
