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
	ws         *websocket.Conn
	close      chan bool
	wg         sync.WaitGroup

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
		dataIn:     l,
		ws:         ws,
		close:      make(chan bool),

		RTCodec: websocket.Codec{websocket.JSON.Marshal, commandUnmarshal},
	}

	defer rs.Close(wg)
	rs.wg.Add(2)
	go rs.Recv(rs.RTCodec, rs.ws)
	go rs.Send(rs.RTCodec, rs.ws)
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

func (rs *RealtimeServer) processCommand(tc TelemetryCommand) {
	switch tc.Cmd {
	case Subscribe:
		rs.subscribed[tc.ID] = true
	case Unsubscribe:
		delete(rs.subscribed, tc.ID)
	}
}

// Send data through the websocket when available. Process the subscription commands.
func (rs *RealtimeServer) Send(c websocket.Codec, ws *websocket.Conn) {
	var d Telemetry
	var i interface{}
	var rtc TelemetryCommand
	var err error

	for {
		select {
		case i = <-rs.dataIn.Listen():
			d = i.(Telemetry)

			if _, ok := rs.subscribed[d.Name]; ok {
				// fmt.Println("Sending: ", d)
				err = c.Send(ws, d)
			}

			if err != nil {
				fmt.Println(err)
				rs.close <- true
			}

		case rtc = <-rs.cmdChannel:
			rs.processCommand(rtc)
		case <-rs.close:
			rs.wg.Done()
			return
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
