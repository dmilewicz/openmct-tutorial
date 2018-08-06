package main

import (
	"bsonparser"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "12345"
	CONN_TYPE = "tcp"
)

type TelemetryBuffer struct {
	Name      string      `json:"name" bson:"name"`
	Key       string      `json:"key" bson:"key"`
	Flags     int64       `json:"flags,omitempty" bson:"flags,omitempty"`
	Timestamp time.Time   `json:"timestamp" bson:"timestamp"`
	Raw_Type  int64       `json:"raw_type" bson:"raw_type"`
	Raw_Value interface{} `json:"raw_value" bson:"raw_value"`
	Eng_Type  int64       `json:"eng_type,omitempty" bson:"eng_type,omitempty"`
	Eng_Val   interface{} `json:"eng_val,omitempty" bson:"eng_val,omitempty"`
}

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	p := bsonparser.InitBuild().BufLen(1024).ParseTo(TelemetryBuffer{}).Build()
	go func() {
		// var mapbuf bsonparser.TelemetryBuffer
		i := 0

		ch := make(chan TelemetryBuffer)

		go func() {
			for {
				var t TelemetryBuffer
				err := p.Next(&t)

				if err != nil {
					fmt.Println("errrrrr: ", err)
					return
				} else {
					ch <- t
				}
			}
		}()

		// p.GetDataChan(&ch)
		// var v TelemetryBuffer
		// p.Next(&v)

		for v := range ch {
			if i%50 == 0 {
				fmt.Println("read: ", v)

			}

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
