package craftsim

import (
	"server"
)

type Simulator struct {
	history        []server.Telemetry
	RealtimeData   chan server.Datum
	HistoryRequest chan server.DataRequest
	HistoryData    chan []server.Datum
	DataIn         chan SpaceCraft
}

// type Datum struct {
// 	Sp        SpaceCraft
// 	Timestamp int64
// }

func NewSim() Simulator {
	return Simulator{[]server.Telemetry{}, make(chan server.Datum), make(chan server.DataRequest), make(chan []server.Datum), make(chan SpaceCraft)}
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
		// sp :=

		d := server.LoadTelemetry(<-s.DataIn)

		s.RealtimeData <- d.Datum("prop.fuel")

		s.history = append(s.history, d)
	}
}

func (s *Simulator) getHistory(start int64, end int64) []server.Datum {
	var telem []server.Datum

	for i, _ := range s.history {
		if s.history[i].Timestamp > start && s.history[i].Timestamp < end {
			telem = append(telem, s.history[i].Datum("prop.fuel"))
		}
	}

	return telem
}
