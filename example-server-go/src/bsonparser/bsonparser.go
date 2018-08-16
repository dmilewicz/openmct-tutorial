package bsonparser

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"time"

	"labix.org/v2/mgo/bson"
)

const (
	DOC_LENGTH_NUM_BYTES = 4
)

type Parser interface {
	Parse(io.Reader)
	Next(interface{}) error
}

type Codec interface {
}

// return &bsonParser{counter: Counter{title: "parser", frameSeconds: 3}}

type bsonParser struct {
	bytesRead    int
	delim        uint
	buflen       uint
	msglen       uint
	counter      Counter
	parseType    reflect.Type
	dataDelivery reflect.Value //*TelemetryBuffer
}

type ParserBuilder interface {
	BufLen(uint) ParserBuilder
	ParseTo(interface{}) ParserBuilder
	Build() Parser
}

func InitBuild() ParserBuilder {
	return new(bsonParser)
}

func (bp *bsonParser) BufLen(buflen uint) ParserBuilder {
	bp.buflen = buflen
	return bp
}

func (bp *bsonParser) ParseTo(i interface{}) ParserBuilder {

	v := reflect.ValueOf(i)
	numChannels := 10

	switch v.Kind() {
	case reflect.Ptr:
		// dereference the element
		bp.parseType = v.Elem().Type()
		bp.dataDelivery = reflect.MakeChan(reflect.ChanOf(reflect.BothDir, v.Type()), 2)
		break
	case reflect.Struct:
		bp.parseType = v.Type()
		bp.dataDelivery = reflect.MakeChan(reflect.ChanOf(reflect.BothDir, reflect.PtrTo(v.Type())), numChannels)
		break
	default:

	}

	return bp
}

func (bp *bsonParser) Build() Parser {
	return *bp
}

func (bp bsonParser) Parse(reader io.Reader) {
	go bp.parse(reader)
}

// Gets the next value read into the buffer
func (bp bsonParser) Next(out interface{}) error {
	// make sure that val is pointer to
	in, _ := bp.dataDelivery.Recv()

	// TYPE CHECKING NEEDED
	reflect.ValueOf(out).Elem().Set(in.Elem())

	return nil
}

func (bp *bsonParser) parse(reader io.Reader) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, bp.buflen, bp.buflen+DOC_LENGTH_NUM_BYTES)

	var token, lastValidByteIdx, partialMsgLength uint
	var err error
	var msgOverflowsBuffer, msgTooLong bool

	go bp.counter.run()

	// Read the incoming connection into the buffer.
	for {
		// read into buffer starting at the next
		bp.bytesRead, err = reader.Read(buf[bp.delim:])
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			return
		}

		token = 0
		lastValidByteIdx = uint(bp.bytesRead) + bp.delim

		for token <= lastValidByteIdx-DOC_LENGTH_NUM_BYTES {
			// token must always point at the beginning of the message
			bp.msglen = uint(binary.LittleEndian.Uint32(buf[token : token+4]))

			// msgLengthInvalid := token > lastValidByteIdx-DOC_LENGTH_NUM_BYTES
			msgOverflowsBuffer = bp.msglen+token > lastValidByteIdx
			msgTooLong = bp.msglen > bp.buflen

			// if message is longer than read in data
			if msgTooLong {
				// set to appropriate length.
				newBuffer := make([]byte, bp.msglen, bp.msglen+DOC_LENGTH_NUM_BYTES) // how to reset this after large message read?
				bp.buflen = bp.msglen
				copy(newBuffer, buf)
				buf = newBuffer
			} else if msgOverflowsBuffer {
				// read next bit of data
				break
			}

			// read the bson document
			v := reflect.New(bp.parseType)
			err = bson.Unmarshal(buf[token:token+bp.msglen], v.Interface())

			if err != nil {
				fmt.Println("ERROR: ", err)
				return
				// deal with it
			} else {
				bp.counter.Add(1)
				bp.dataDelivery.Send(v)
			}

			// set token to the start of the next document
			token += bp.msglen
		}

		// rewind
		partialMsgLength = lastValidByteIdx - token

		// put message at front of buffer
		copy(buf, buf[token:lastValidByteIdx])

		// mark the delimiter
		bp.delim = partialMsgLength
	}

}

type Counter struct {
	title string

	frameSeconds int
	frameCounter int
	countInFrame int
	totalCount   int64
}

// NewCounter()

func (c *Counter) run() {

	for {
		<-time.After(time.Duration(c.frameSeconds) * time.Second)
		c.totalCount += int64(c.countInFrame)
		c.report()
		c.countInFrame = 0
		c.frameCounter++
	}
}

func (c *Counter) Add(count int) {
	c.countInFrame += count
}

func (c *Counter) report() {
	fmt.Println(c.title+" frame rate: ", c.countInFrame/c.frameSeconds, " points/sec")
}
