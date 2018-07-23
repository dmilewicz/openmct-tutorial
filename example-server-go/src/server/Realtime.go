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
	// TelemIn    <-chan Telemetry
	RequestOut chan<- TelemetryCommand
	Subscribed map[string]bool
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

func NewRealtimeServer(r chan TelemetryCommand, l Listener, ws *websocket.Conn) RealtimeServer {
	// configure server
	rs := RealtimeServer{
		RequestOut: r,
		Subscribed: make(map[string]bool),
		cmdChannel: make(chan TelemetryCommand),
		dataIn:     l,
		ws:         ws,
		close:      make(chan bool),

		RTCodec: websocket.Codec{websocket.JSON.Marshal, commandUnmarshal},
	}

	// rs.wg.Add(1)
	rs.RealtimeSocket()
	// rs.wg.Wait()

	return rs
}

/**
 * Send and receive from socket
 **/
func (rs *RealtimeServer) RealtimeSocket() {
	defer rs.Close()
	// rs.wg.Add(2)
	go rs.Recv(rs.RTCodec, rs.ws)
	rs.Send(rs.RTCodec, rs.ws)
}

// Get subscrition commands from the websocket. Send them to the processing thread.
func (rs *RealtimeServer) Recv(c websocket.Codec, ws *websocket.Conn) {
	var rtc TelemetryCommand
	var err error

	for {
		err = c.Receive(ws, &rtc)
		if err != nil {
			// fmt.Println("recv error:", err)
			// rs.close <- true
			// break
		}
		rs.cmdChannel <- rtc
	}

	// rs.wg.Done()
}

func (rs *RealtimeServer) processCommand(tc TelemetryCommand) {
	switch tc.Cmd {
	case Subscribe:
		rs.Subscribed[tc.ID] = true
	case Unsubscribe:
		delete(rs.Subscribed, tc.ID)
	}
}

/**
 *	Send data through the websocket when available. Process the subscription commands.
 */
func (rs *RealtimeServer) Send(c websocket.Codec, ws *websocket.Conn) {
	var d Telemetry
	var i interface{}
	var rtc TelemetryCommand
	var err error

	for {
		select {
		case i = <-rs.dataIn.Listen():
			d = i.(Telemetry)
			for key := range rs.Subscribed {

				err = c.Send(ws, d.Datum(key))

				// closing system not really working
				if err != nil {
					// fmt.Println("closing here")

					fmt.Println(err)

					rs.Close()
					return
					// rs.close <- true
				}
			}
		case rtc = <-rs.cmdChannel:
			rs.processCommand(rtc)
		case <-rs.close:
			fmt.Println("closing here")
			rs.Close()
			// rs.wg.Done()
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

func (rs *RealtimeServer) Close() {
	rs.ws.Close()
	rs.dataIn.CloseListener()
}
