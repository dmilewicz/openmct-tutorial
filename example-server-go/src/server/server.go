package server

import (
	// "bsonparser"
	"fmt"
	"net"
	"net/http"
	"os"
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
	wg       sync.WaitGroup
	close    chan bool
	// parser   bsonparser.Parser

	conn_host string
	conn_port string
	conn_type string
}

type Server interface {
}

func NewServer(port int, datain chan Telemetry, dr chan DataRequest, h chan []Datum) telemetryServer {
	// p := bsonparser.NewBSONParser().SetByteBufferLength(1024).SetChanBufferLength(10)

	s := telemetryServer{
		// parser:   p,
		dispatch: NewDispatch(p.GetDataChan()),
		hserver:  HistoryServer{dr, h},
	}

	return s
}

func (s *telemetryServer) HandleWebsocket(ws *websocket.Conn) {
	s.wg.Add(1)
	NewRealtimeServer(make(chan TelemetryCommand), s.dispatch.NewListener(), ws, s.wg)
}

func (s telemetryServer) RunServer() {
	// define handlers
	http.Handle("/", http.FileServer(http.Dir("../")))
	http.Handle("/realtime/", websocket.Handler(s.HandleWebsocket))
	http.HandleFunc("/history/", s.hserver.RunServer)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func (s *telemetryServer) StreamData() {
	l, err := net.Listen(s.conn_type, s.conn_host+":"+s.conn_port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// p := bsonparser.NewBSONParser().SetByteBufferLength(1024).SetChanBufferLength(10)

	// Close the listener when the application closes.
	defer l.Close()

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		p.Parse(conn)
	}
}
