package bsonparser

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"labix.org/v2/mgo/bson"
)

const (
	DOC_LENGTH_NUM_BYTES = 4
)

type TelemetryBuffer struct {
	Name      string      // `json:"name" bson:"name"`
	Key       string      // `json:"key" bson:"key"`
	Flags     int64       // `json:"flags,omitempty" bson:"flags,omitempty"`
	Timestamp float64     // `json:"timestamp" bson:"timestamp"`
	Raw_Type  int64       // `json:"raw_type" bson:"raw_type"`
	Raw_Value interface{} // `json:"raw_value" bson:"raw_value"`
	Eng_Type  int64       // `json:"eng_type,omitempty" bson:"eng_type,omitempty"`
	Eng_Val   interface{} // `json:"eng_val,omitempty" bson:"eng_val,omitempty"`
}

type Parser interface {
	Parse(io.Reader)

	// setters for building
	// SetByteBufferLength(length uint) Parser
	// SetChanBufferLength(length uint) Parser
	GetDataChan() chan<- interface{}
	Next(interface{}) error

	// Close()
}

type bsonParser struct {
	bytesRead int
	delim     uint
	buflen    uint
	msglen    uint
	parseType reflect.Type

	// dataDelivery chan *bson.M
	dataDelivery reflect.Value //*TelemetryBuffer
}

// func (bp bsonParser) SetByteBufferLength(length uint) Parser {
// 	bp.buflen = length
// 	return bp
// }

// func (bp bsonParser) SetChanBufferLength(length uint) Parser {
// 	bp.dataDelivery = make(chan *Telemetry, length)
// 	return bp
// }

func NewBSONParser() bsonParser {

	msgBufferSize := 10

	return bsonParser{
		bytesRead: 0,
		delim:     0,
		buflen:    1024,
		msglen:    0,
		parseType: reflect.TypeOf(TelemetryBuffer{}),
		// dataDelivery: make(chan *bson.M, msgBufferSize),
		dataDelivery: reflect.MakeChan(reflect.ChanOf(reflect.BothDir, reflect.TypeOf(&TelemetryBuffer{})), msgBufferSize),
		// dataDelivery: make(chan *bson.M, msgBufferSize),

	}
}

func (bp bsonParser) Parse(reader io.Reader) {
	go bp.parse(reader)
}

// func (bp *bsonParser) SetByteStream(reader io.Reader) {
// 	bp.reader = reader
// }

// func (bp bsonParser) Next(mapBuffer *TelemetryBuffer) error {

// 	*mapBuffer = *(<-bp.dataDelivery)

// 	return nil
// }

// Gets the next value read into the buffer
func (bp bsonParser) Next(out interface{}) error {
	// make sure that val is pointer to

	in, _ := bp.dataDelivery.Recv()

	// TYPE CHECKING NEEDED

	reflect.ValueOf(out).Elem().Set(in.Elem())

	return nil
}

// 	// switch s.(type) {
// 	// case bp.parseType:
// 	// 	v = s
// 	// 	return
// 	// }

// }

// func (bp bsonParser) SetMsgType(t Type) {
// 	bp.parseType =
// }

func (bp bsonParser) GetDataChan(ch interface{}) error {
	ch = bp.dataDelivery.Interface()

	return nil
}

func (bp *bsonParser) report() {
	fmt.Println("bsonParser:")
	fmt.Println("   bytesRead: ", bp.bytesRead)
	fmt.Println("   delim: ", bp.delim)
	fmt.Println("   buflen: ", bp.buflen)
	fmt.Println("   msglen: ", bp.msglen)
}

func (bp *bsonParser) parse(reader io.Reader) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, bp.buflen, bp.buflen+DOC_LENGTH_NUM_BYTES)
	// sendbuf
	// mapadap := make(map[string]interface{})

	var token, lastValidByteIdx, partialMsgLength uint
	var err error
	var msgOverflowsBuffer, msgTooLong bool

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
			// msg := make([]byte, bp.msglen)

			// copy(msg, buf[token:token+bp.msglen])

			if err != nil {
				fmt.Println("ERROR: ", err)
				// bp.report()
				return
				// deal with it
			} else {
				fmt.Println(v)
				bp.dataDelivery.Send(v) //reflect.ValueOf(v))
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
