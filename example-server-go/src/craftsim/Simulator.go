package craftsim

import (
	"server"
	"time"
)

type Simulator struct {
	history        []Datum
	RealtimeData   chan server.Telemetry
	HistoryRequest chan server.DataRequest
	HistoryData    chan []server.Telemetry
	DataIn         chan SpaceCraft
}

type Datum struct {
	Sp        SpaceCraft
	Timestamp int64
}

func NewSim() Simulator {
	return Simulator{[]Datum{}, make(chan server.Telemetry), make(chan server.DataRequest), make(chan []server.Telemetry), make(chan SpaceCraft)}
}

func (s *Simulator) RunSim() {
	sp := NewSpacecraft()
	go s.storeHistory()
	go sp.RunSpacecraft(s.DataIn)

	for {
		dr := <-s.HistoryRequest
		s.HistoryData <- s.getHistory(dr.Start, dr.End)
	}

}

func (s *Simulator) storeHistory() {
	for {
		d := Datum{<-s.DataIn, time.Now().UnixNano() / int64(time.Millisecond)}

		s.RealtimeData <- server.Telemetry{d.Timestamp, d.Sp.PropFuel, "prop.fuel"}

		s.history = append(s.history, d)
	}
}

func (s *Simulator) getHistory(start int64, end int64) []server.Telemetry {
	var telem []server.Telemetry

	for i, _ := range s.history {
		if s.history[i].Timestamp > start && s.history[i].Timestamp < end {
			telem = append(telem, server.Telemetry{s.history[i].Timestamp, s.history[i].Sp.PropFuel, "prop.fuel"})
		}
	}

	return telem
}
