package server

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
)

type Command int

const (
	Subscribe   Command = 0
	Unsubscribe Command = 1
)

type RealtimeServer struct {
	RequestOut chan<- TelemetryCommand
	subscribed map[string]bool
	cmdChannel chan TelemetryCommand
	dataIn     Listener
	dataOut    chan Telemetry
	ws         *websocket.Conn
	close      chan bool
	wg         sync.WaitGroup
	counter    Counter

	RTCodec websocket.Codec
}

type TelemetryCommand struct {
	Cmd Command
	ID  string
}

func NewRealtimeServer(r chan TelemetryCommand, l Listener, ws *websocket.Conn, wg sync.WaitGroup) RealtimeServer {
	// configure server
	rs := RealtimeServer{
		RequestOut: r,
		subscribed: make(map[string]bool),
		cmdChannel: make(chan TelemetryCommand),
		dataOut:    make(chan Telemetry),
		dataIn:     l,
		ws:         ws,
		close:      make(chan bool),
		counter:    Counter{title: "client", frameSeconds: 3},

		RTCodec: websocket.Codec{websocket.JSON.Marshal, commandUnmarshal},
	}

	defer rs.Close(wg)
	rs.wg.Add(2)
	go rs.Recv(rs.RTCodec, rs.ws)
	go rs.Send(rs.RTCodec, rs.ws)
	go rs.counter.run()
	rs.wg.Wait()

	return rs
}

// Send and receive from socket
func (rs *RealtimeServer) RealtimeSocket() {
	rs.wg.Add(2)
	go rs.Recv(rs.RTCodec, rs.ws)
	go rs.Send(rs.RTCodec, rs.ws)
	rs.wg.Wait()
}

// Get subscrition commands from the websocket. Send them to the processing thread.
func (rs *RealtimeServer) Recv(c websocket.Codec, ws *websocket.Conn) {
	var rtc TelemetryCommand
	var err error

	for {
		err = c.Receive(ws, &rtc)
		if err != nil {
			fmt.Println("recv error:", err)
			rs.close <- true
			break
		}
		rs.cmdChannel <- rtc
	}

	rs.wg.Done()
}

// mark telemetry subscribe/unsubscribe
func (rs *RealtimeServer) processCommand(tc TelemetryCommand) {
	switch tc.Cmd {
	case Subscribe:
		rs.subscribed[tc.ID] = true
	case Unsubscribe:
		delete(rs.subscribed, tc.ID)
	}
}

func (rs *RealtimeServer) run() {
	var d Telemetry
	var i interface{}
	var rtc TelemetryCommand

	for {
		select {
		case i = <-rs.dataIn.Listen():
			d = i.(Telemetry)

			if _, ok := rs.subscribed[d.Name]; ok {
				rs.counter.Add(1)
				rs.dataOut <- d
			}
		case rtc = <-rs.cmdChannel:
			rs.processCommand(rtc)
		case <-rs.close:
			rs.wg.Done()
			return
		}
	}
}


// Send data through the websocket when available. Process the subscription commands.
func (rs *RealtimeServer) Send(c websocket.Codec, ws *websocket.Conn) {
	var d Telemetry
	var err error

	for d = range rs.dataOut {
		err = c.Send(ws, d)

		if err != nil {
			fmt.Println(err)
			rs.close <- true
		}
	}
}

func commandUnmarshal(data []byte, payloadType byte, v interface{}) (err error) {
	cmdString := string(data)

	// ASSUMES NO SPACES IN BODY
	cmds := strings.Split(cmdString, " ")
	var cmd Command

	switch cmds[0] {
	case "subscribe":
		cmd = Subscribe
	case "unsubscribe":
		cmd = Unsubscribe
	default:
		return errors.New("Not a valid command")
	}

	switch data := v.(type) {
	case *TelemetryCommand:
		*data = TelemetryCommand{cmd, cmds[1]}
	}

	return nil
}

func (rs *RealtimeServer) Close(wg sync.WaitGroup) {
	wg.Done()
	rs.ws.Close()
	rs.dataIn.CloseListener()
}
