package server

import (
	"fmt"
	"time"
)

type HistoryStore struct {
	DataIn         chan TelemetryBuffer
	HistoryRequest chan DataRequest
	HistoryData    chan []TelemetryBuffer

	history map[string][]TelemetryBuffer
}

func NewHistoryStore() HistoryStore {
	hs := HistoryStore{
		DataIn:         make(chan TelemetryBuffer),
		HistoryRequest: make(chan DataRequest),
		HistoryData:    make(chan []TelemetryBuffer),
		history:        make(map[string][]TelemetryBuffer),
	}

	go hs.respond()

	return hs
}

func (h *HistoryStore) storeHistory() {
	for {
		d := <-h.DataIn

		h.history[d.Name] = append(h.history[d.Name], d)
	}
}

func (h *HistoryStore) respond() {
	go h.storeHistory()

	for {
		dr := <-h.HistoryRequest
		h.HistoryData <- h.getHistory(dr)
	}
}

func (h *HistoryStore) getHistory(dr DataRequest) []TelemetryBuffer {
	var telem []TelemetryBuffer

	for _, v := range h.history[dr.Value] {

		if time.Time(v.Timestamp).After(dr.Start) && time.Time(v.Timestamp).Before(dr.End) {
			telem = append(telem, v)
		}

	}

	fmt.Println(telem)

	return telem
}
