package craftsim

import "server"

type Simulator struct {
	history        []server.Telemetry
	RealtimeData   chan server.Telemetry
	HistoryRequest chan server.DataRequest
	HistoryData    chan []server.Datum
	DataIn         chan SpaceCraft
}

func NewSim() Simulator {
	s := Simulator{[]server.Telemetry{}, make(chan server.Telemetry), make(chan server.DataRequest), make(chan []server.Datum), make(chan SpaceCraft)}
	go s.RunSim()
	return s
}

func (s *Simulator) RunSim() {
	sp := NewSpacecraft()
	go s.storeHistory()
	go sp.RunSpacecraft(s.DataIn)

	for {
		dr := <-s.HistoryRequest
		s.HistoryData <- s.getHistory(dr)
	}
}

func (s *Simulator) storeHistory() {
	for {
		d := server.LoadTelemetry(<-s.DataIn)

		s.RealtimeData <- d

		s.history = append(s.history, d)
	}
}

func (s *Simulator) getHistory(dr server.DataRequest) []server.Datum {
	var telem []server.Datum

	for i, _ := range s.history {
		if s.history[i].Timestamp > dr.Start && s.history[i].Timestamp < dr.End {
			telem = append(telem, s.history[i].Datum(dr.Value))
		}
	}

	return telem
}
