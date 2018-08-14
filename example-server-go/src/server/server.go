package server

import (
	// "bsonparser"
	"bsonparser"
	"fmt"
	"net/http"
	"sync"
	// "fmt"
	// "encoding/json"
	// "io/ioutil"

	"golang.org/x/net/websocket"
)

type telemetryServer struct {
	dispatch Dispatcher
	hserver  HistoryServer
	dictgen  DictionaryGenerator
	parser   bsonparser.Parser
	wg       sync.WaitGroup
	close    chan bool

	conn_host string
	conn_port string
	conn_type string
}

type Server interface {
}

func NewServer(port int, dr chan DataRequest, h chan []Telemetry, hs chan Telemetry) telemetryServer {
	p := bsonparser.InitBuild().BufLen(1024).ParseTo(Telemetry{}).Build()

	dataChan := make(chan Telemetry)

	// Data Ingestion
	go func() {
		for {
			var t Telemetry
			err := p.Next(&t)

			if err != nil {
				fmt.Println("err: ", err)
				return
			} else {
				dataChan <- t
			}
		}
	}()

	dictChan := make(chan Telemetry)

	s := telemetryServer{
		parser:   p,
		dispatch: NewDispatch(dataChan, hs, dictChan),
		dictgen:  NewDictionaryGenerator(dictChan),
		hserver:  HistoryServer{dr, h},
	}

	return s
}

func (s *telemetryServer) HandleWebsocket(ws *websocket.Conn) {
	s.wg.Add(1)
	NewRealtimeServer(make(chan TelemetryCommand), s.dispatch.NewListener(), ws, s.wg)
}

func (s telemetryServer) RunServer() {

	go ReadData(s.parser)

	// define handlers
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/realtime/", websocket.Handler(s.HandleWebsocket))
	http.HandleFunc("/history/", s.hserver.RunServer)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
