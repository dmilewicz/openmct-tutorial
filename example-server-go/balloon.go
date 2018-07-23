package main

import (
	"bsonparser"
	"fmt"
	"net"
	"os"

	"labix.org/v2/mgo/bson"
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

	p := bsonparser.NewBSONParser().SetByteBufferLength(1024).SetChanBufferLength(10)
	go func() {
		var mapbuf bson.M
		for {
			p.Next(&mapbuf)
			fmt.Println("read: ", mapbuf)
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
