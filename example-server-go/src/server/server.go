package server

import (
	"net/http"
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

	// 	conn_host string
	// 	conn_port string
	// 	conn_type string
}

type Server interface {
}

func NewServer(port int, datain chan Telemetry, dr chan DataRequest, h chan []Datum) {
	s := telemetryServer{
		dispatch: NewDispatch(datain),
		hserver:  HistoryServer{dr, h},
	}

	s.RunServer()
}

func (s *telemetryServer) HandleWebsocket(ws *websocket.Conn) {
	NewRealtimeServer(make(chan TelemetryCommand), s.dispatch.NewListener(), ws)
}

func (s *telemetryServer) RunServer() {
	// define handlers
	http.Handle("/", http.FileServer(http.Dir("../")))
	http.Handle("/realtime/", websocket.Handler(s.HandleWebsocket))
	http.HandleFunc("/history/", s.hserver.RunServer)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// func (s *server) StreamData() {
// 	l, err := net.Listen(conn_type, conn_host+":"+conn_port)
// 	if err != nil {
// 		fmt.Println("Error listening:", err.Error())
// 		os.Exit(1)
// 	}

// 	p := bsonparser.NewBSONParser().SetByteBufferLength(1024).SetChanBufferLength(10)
// 	// go func() {
// 	// 	var mapbuf bson.M
// 	// 	for {
// 	// 		p.Next(&mapbuf)
// 	// 		fmt.Println("read: ", mapbuf)
// 	// 	}
// 	// }()

// 	// Close the listener when the application closes.
// 	defer l.Close()
// 	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
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
