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

// const (
// 	CONN_HOST = "localhost"
// 	CONN_PORT = "12345"
// 	CONN_TYPE = "tcp"
// )

type telemetryServer struct {
	dispatch Dispatcher
	hserver  HistoryServer
	parser   bsonparser.Parser
	wg       sync.WaitGroup
	close    chan bool
	// parser   bsonparser.Parser

	conn_host string
	conn_port string
	conn_type string
}

type Server interface {
}

func NewServer(port int, dr chan DataRequest, h chan []TelemetryBuffer, hs chan TelemetryBuffer) telemetryServer {
	// p := bsonparser.NewBSONParser().SetByteBufferLength(1024).SetChanBufferLength(10)

	p := bsonparser.InitBuild().BufLen(1024).ParseTo(TelemetryBuffer{}).Build()

	dataChan := make(chan TelemetryBuffer)

	go func() {
		for {
			var t TelemetryBuffer
			err := p.Next(&t)

			if err != nil {
				fmt.Println("err: ", err)
				return
			} else {
				dataChan <- t
			}
		}
	}()

	go ReadData(p)

	s := telemetryServer{
		parser:   p,
		dispatch: NewDispatch(dataChan, hs),
		hserver:  HistoryServer{dr, h},
	}

	fmt.Println("end newserver")

	return s
}

func (s *telemetryServer) HandleWebsocket(ws *websocket.Conn) {
	s.wg.Add(1)
	NewRealtimeServer(make(chan TelemetryCommand), s.dispatch.NewListener(), ws, s.wg)
}

func (s telemetryServer) RunServer() {
	fmt.Println("runserver")

	// define handlers
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/realtime/", websocket.Handler(s.HandleWebsocket))
	http.HandleFunc("/history/", s.hserver.RunServer)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// func (s *telemetryServer) StreamData() {
// 	l, err := net.Listen(s.conn_type, s.conn_host+":"+s.conn_port)
// 	if err != nil {
// 		fmt.Println("Error listening:", err.Error())
// 		os.Exit(1)
// 	}

// 	// p := bsonparser.NewBSONParser().SetByteBufferLength(1024).SetChanBufferLength(10)

// 	// Close the listener when the application closes.
// 	defer l.Close()

// 	for {
// 		// Listen for an incoming connection.
// 		conn, err := l.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting: ", err.Error())
// 			os.Exit(1)
// 		}
// 		// Handle connections in a new goroutine.
// 		p.Parse(conn)
// 	}
// }
