package main

import (
	"bsonparser"
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "12345"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	p := bsonparser.NewBSONParser()
	go func() {
		// var mapbuf bsonparser.TelemetryBuffer
		i := 0

		var ch chan bsonparser.TelemetryBuffer

		p.GetDataChan(ch)

		for {
			var v bsonparser.TelemetryBuffer
			p.Next(&v)

			// S := <-ch

			// if i%50 == 0 {
			fmt.Println("read: ", v)

			// }
			i++
		}
	}()

	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		p.Parse(conn)
	}

}
