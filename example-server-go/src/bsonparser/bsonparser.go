package bsonparser

import (
	"encoding/binary"
	"fmt"
	"io"

	"labix.org/v2/mgo/bson"
)

const (
	DOC_LENGTH_NUM_BYTES = 4
)

type Parser interface {
	Parse(io.Reader)

	Next(map[string]interface{}) error

	SetByteBufferLength(length uint) *bsonParser
	SetChanBufferLength(length uint) *bsonParser

	Close()
}

type bsonParser struct {
	bytesRead    int
	delim        uint
	buflen       uint
	msglen       uint
	dataDelivery chan *bson.M
}

func (bp bsonParser) SetByteBufferLength(length uint) bsonParser {
	bp.buflen = length
	return bp
}

func (bp bsonParser) SetChanBufferLength(length uint) bsonParser {
	bp.dataDelivery = make(chan *bson.M, length)
	return bp
}

func NewBSONParser() bsonParser {

	msgBufferSize := 10

	return bsonParser{
		bytesRead:    0,
		delim:        0,
		buflen:       1024,
		msglen:       0,
		dataDelivery: make(chan *bson.M, msgBufferSize),
	}
}

func (bp *bsonParser) Parse(reader io.Reader) {
	go bp.parse(reader)
}

// func (bp *bsonParser) SetByteStream(reader io.Reader) {
// 	bp.reader = reader
// }

func (bp *bsonParser) Next(mapBuffer *bson.M) error {
	*mapBuffer = *(<-bp.dataDelivery)
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
	// mapadap := make(map[string]interface{})

	var token, lastValidByteIdx, partialMsgLength uint
	var err error
	var msgOverflowsBuffer, msgTooLong bool
	fmt.Println("starting parse")

	// Read the incoming connection into the buffer.
	for {
		// read into buffer starting at the next
		bp.bytesRead, err = reader.Read(buf[bp.delim:])
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
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

			// create new map pointer
			doc_map := new(bson.M)

			// read the bson document
			err = bson.Unmarshal(buf[token:], doc_map)
			if err != nil {
				fmt.Println("ERROR: ", err)
				bp.report()
				fmt.Println("   token:", token)
				return
				// deal with it
			} else {
				bp.dataDelivery <- doc_map
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
