package server

import (
	"fmt"
	"strconv"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"labix.org/v2/mgo/bson"
)

type HistoryStore struct {
	db             *leveldb.DB
	DataIn         chan Telemetry
	HistoryRequest chan DataRequest
	HistoryData    chan []Telemetry

	history map[string][]Telemetry
}

type Database interface {
	Get(string) ([]byte, error)
	Put(string, []byte) error
	Delete(string) error
}

func NewHistoryStore() HistoryStore {
	db, err := leveldb.OpenFile("db/leveldb", nil)
	if err != nil {
		panic(err)
	}

	hs := HistoryStore{
		db:             db,
		DataIn:         make(chan Telemetry),
		HistoryRequest: make(chan DataRequest),
		HistoryData:    make(chan []Telemetry),
		history:        make(map[string][]Telemetry),
	}

	go hs.respond()

	return hs
}

func (h *HistoryStore) storeHistory() {
	for {
		d := <-h.DataIn
		name := d.Name + strconv.FormatInt(time.Time(d.Timestamp).UnixNano()/int64(time.Millisecond), 10)

		// get the value bytes
		buf, err := bson.Marshal(d)

		if err != nil {
			fmt.Println("ERROR: ", err)
			continue
		}

		err = h.db.Put([]byte(name), buf, nil)

	}
}

func (h *HistoryStore) respond() {
	go h.storeHistory()

	for dr := range h.HistoryRequest {
		h.HistoryData <- h.getHistory(dr)
	}
}

func (h *HistoryStore) getHistory(dr DataRequest) []Telemetry {
	var telem []Telemetry
	var err error

	r := util.Range{
		Start: []byte(dr.Value + strconv.FormatInt(dr.Start.UnixNano()/int64(time.Millisecond), 10)),
		Limit: []byte(dr.Value + strconv.FormatInt(dr.End.UnixNano()/int64(time.Millisecond), 10)),
	}

	iter := h.db.NewIterator(&r, nil)

	for iter.Next() {
		t := new(Telemetry)
		err = bson.Unmarshal(iter.Value(), t)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		telem = append(telem, *t)
	}

	fmt.Println(len(telem))

	return telem
}
