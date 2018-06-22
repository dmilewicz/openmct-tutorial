package server

import (
	// "net/http"
	// "net"
	// "bytes"

	// "crypto/sha1"
	// "craftsim"
	// "encoding/base64"
	// "strconv"
	"fmt"
	// "io/ioutil"
	"encoding/json"
	"errors"
	"strings"

	"golang.org/x/net/websocket"
)

type Command int

const (
	Subscribe   Command = 0
	Unsubscribe Command = 1
)

type RealtimeServer struct {
	TelemIn    <-chan Datum
	RequestOut chan<- TelemetryCommand

	RTCodec websocket.Codec
}

type TelemetryCommand struct {
	Cmd Command
	ID  string
}

func NewRealtimeServer(t chan Datum, r chan TelemetryCommand) RealtimeServer {
	return RealtimeServer{t, r, websocket.Codec{websocket.JSON.Marshal, commandUnmarshal}}
}

func describe(i interface{}) {
	fmt.Printf("(%v)\n", i)
}

func commandUnmarshal(data []byte, payloadType byte, v interface{}) (err error) {
	cmdString := string(data)

	fmt.Println("raw data: " + cmdString)

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

func jsonMarshal(v interface{}) (msg []byte, payloadType byte, err error) {
	msg, err = json.Marshal(v)

	// TODO: replace magic number
	return msg, 1, err
}

func (rs *RealtimeServer) RealtimeSocket(ws *websocket.Conn) {

	fmt.Println("rtsocket called again")
	rs.Snd(rs.RTCodec, ws)

	go Rcv(rs.RTCodec, ws)
}

func Rcv(c websocket.Codec, ws *websocket.Conn) {
	var rtc TelemetryCommand

	// for {
	// 	<-time.After(time.Second)
	fmt.Println("in rcv")

	c.Receive(ws, &rtc)
	// }
}

func (rs *RealtimeServer) Snd(c websocket.Codec, ws *websocket.Conn) {
	for {
		t := <-rs.TelemIn
		c.Send(ws, t)
	}
}

// func (rs *RealtimeServer) RealtimeSocket(ws *websocket.Conn) {
// var buf = make([]byte, 1024)
// // n, err := ws.Read(buf)
// if err != nil {
// 	fmt.Println(err)
// 	return
// }

// for {

// 	var rtc RealtimeTelemetry

// 	rs.RTCodec.Receive(ws, rtc)

// 	fmt.Println("We read", rtc.ID)

// 	t := <-rs.Telem
// 	t_json, _ := json.Marshal(t)

// 	ws.Write(t_json)
// 	fmt.Println("looped")
// }

// }

// var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")

// // From Gorilla
// func computeAcceptKey(challengeKey string) string {
// 	h := sha1.New()
// 	h.Write([]byte(challengeKey))
// 	h.Write(keyGUID)
// 	return base64.StdEncoding.EncodeToString(h.Sum(nil))
// }

// /*
//  * Manual handshake completion. Static and universal. No error checking.
//  *
//  */
// func completeHandshake(conn *net.Conn, r *http.Request) {
// 	var buffer bytes.Buffer
// 	buffer.WriteString("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nContent-Type: application/json\r\nSec-WebSocket-Accept: ")
// 	buffer.WriteString(computeAcceptKey(r.Header.Get("Sec-Websocket-Key")))
// 	buffer.WriteString("\r\n\r\n")

// 	(*conn).Write(buffer.Bytes())
// }

// func (rs *RealtimeServer) RunServer(w http.ResponseWriter, r *http.Request) {
// 	hj, ok := w.(http.Hijacker)
// 	if !ok {
// 		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
// 		return
// 	}

// 	conn, brw, err := hj.Hijack()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if brw.Reader.Buffered() > 0 {
// 		conn.Close() // unread data from client
// 		fmt.Println("ERROR")
// 	}

// 	// send handshake completion
// 	completeHandshake(&conn, r)

// 	// res := "{\"timestamp\":1529016469335,\"value\":77,\"id\":\"prop.fuel\"}"

// 	defer conn.Close()

// 	var p []byte
// 	fmt.Println("connecting realtime...")

// 	for {
// 		// <-time.After(time.Second)

// 		// p, _ := ioutil.ReadAll(conn)
// 		p_len, _ := conn.Read(p)
// 		// p_len := len(p)

// 		if p_len > 0 {
// 			fmt.Printf("Message length: " + strconv.Itoa(p_len) + ". Message: %s", p)
// 		}
// 		// else {
// 		// 	fmt.Println("Message length: " + strconv.Itoa(p_len))
// 		// }

// 		// conn.Write([]byte(res))
// 	}

// 	fmt.Println("End realtime")

// }

// // Upgrade upgrades the HTTP server connection to the WebSocket protocol.
// //
// // The responseHeader is included in the response to the client's upgrade
// // request. Use the responseHeader to specify cookies (Set-Cookie) and the
// // application negotiated subprotocol (Sec-WebSocket-Protocol).
// //
// // If the upgrade fails, then Upgrade replies to the client with an HTTP error
// // response.
// func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error) {
// 	const badHandshake = "websocket: the client is not using the websocket protocol: "

// 	if !tokenListContainsValue(r.Header, "Connection", "upgrade") {
// 		return u.returnError(w, r, http.StatusBadRequest, badHandshake+"'upgrade' token not found in 'Connection' header")
// 	}

// 	if !tokenListContainsValue(r.Header, "Upgrade", "websocket") {
// 		return u.returnError(w, r, http.StatusBadRequest, badHandshake+"'websocket' token not found in 'Upgrade' header")
// 	}

// 	if r.Method != "GET" {
// 		return u.returnError(w, r, http.StatusMethodNotAllowed, badHandshake+"request method is not GET")
// 	}

// 	// if !tokenListContainsValue(r.Header, "Sec-Websocket-Version", "13") {
// 	// 	return u.returnError(w, r, http.StatusBadRequest, "websocket: unsupported version: 13 not found in 'Sec-Websocket-Version' header")
// 	// }

// 	// if _, ok := responseHeader["Sec-Websocket-Extensions"]; ok {
// 	// 	return u.returnError(w, r, http.StatusInternalServerError, "websocket: application specific 'Sec-WebSocket-Extensions' headers are unsupported")
// 	// }

// 	checkOrigin := u.CheckOrigin
// 	if checkOrigin == nil {
// 		checkOrigin = checkSameOrigin
// 	}
// 	if !checkOrigin(r) {
// 		return u.returnError(w, r, http.StatusForbidden, "websocket: request origin not allowed by Upgrader.CheckOrigin")
// 	}

// 	challengeKey := r.Header.Get("Sec-Websocket-Key")
// 	if challengeKey == "" {
// 		return u.returnError(w, r, http.StatusBadRequest, "websocket: not a websocket handshake: `Sec-WebSocket-Key' header is missing or blank")
// 	}

// 	subprotocol := u.selectSubprotocol(r, responseHeader)

// 	// Negotiate PMCE
// 	var compress bool
// 	if u.EnableCompression {
// 		for _, ext := range parseExtensions(r.Header) {
// 			if ext[""] != "permessage-deflate" {
// 				continue
// 			}
// 			compress = true
// 			break
// 		}
// 	}

// 	var (
// 		netConn net.Conn
// 		err     error
// 	)

// 	h, ok := w.(http.Hijacker)
// 	if !ok {
// 		return u.returnError(w, r, http.StatusInternalServerError, "websocket: response does not implement http.Hijacker")
// 	}
// 	var brw *bufio.ReadWriter
// 	netConn, brw, err = h.Hijack()
// 	if err != nil {
// 		return u.returnError(w, r, http.StatusInternalServerError, err.Error())
// 	}

// 	if brw.Reader.Buffered() > 0 {
// 		netConn.Close()
// 		return nil, errors.New("websocket: client sent data before handshake is complete")
// 	}

// 	c := newConnBRW(netConn, true, u.ReadBufferSize, u.WriteBufferSize, brw)
// 	c.subprotocol = subprotocol

// 	if compress {
// 		c.newCompressionWriter = compressNoContextTakeover
// 		c.newDecompressionReader = decompressNoContextTakeover
// 	}

// 	p := c.writeBuf[:0]
// 	p = append(p, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "...)
// 	p = append(p, computeAcceptKey(challengeKey)...)
// 	p = append(p, "\r\n"...)
// 	if c.subprotocol != "" {
// 		p = append(p, "Sec-WebSocket-Protocol: "...)
// 		p = append(p, c.subprotocol...)
// 		p = append(p, "\r\n"...)
// 	}
// 	if compress {
// 		p = append(p, "Sec-WebSocket-Extensions: permessage-deflate; server_no_context_takeover; client_no_context_takeover\r\n"...)
// 	}
// 	for k, vs := range responseHeader {
// 		if k == "Sec-Websocket-Protocol" {
// 			continue
// 		}
// 		for _, v := range vs {
// 			p = append(p, k...)
// 			p = append(p, ": "...)
// 			for i := 0; i < len(v); i++ {
// 				b := v[i]
// 				if b <= 31 {
// 					// prevent response splitting.
// 					b = ' '
// 				}
// 				p = append(p, b)
// 			}
// 			p = append(p, "\r\n"...)
// 		}
// 	}
// 	p = append(p, "\r\n"...)

// 	// Clear deadlines set by HTTP server.
// 	netConn.SetDeadline(time.Time{})

// 	if u.HandshakeTimeout > 0 {
// 		netConn.SetWriteDeadline(time.Now().Add(u.HandshakeTimeout))
// 	}
// 	if _, err = netConn.Write(p); err != nil {
// 		netConn.Close()
// 		return nil, err
// 	}
// 	if u.HandshakeTimeout > 0 {
// 		netConn.SetWriteDeadline(time.Time{})
// 	}

// 	return c, nil
// }
