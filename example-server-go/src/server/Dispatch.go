package server

import (
	"broadcast"
)

type Dispatcher struct {
	TelemIn chan Telemetry

	telemetryIn broadcast.Broadcaster

	hData    chan<- Telemetry
	dictData chan<- Telemetry
}

// type Listener interface {
// 	Read() interface{}
// 	CloseListener()
// }

type Listener struct {
	dataIn chan interface{}
	b      broadcast.Broadcaster
}

// func (l *listener) Read() interface{} {
// 	return <-l.dataIn
// }

func (l *Listener) CloseListener() {
	l.b.Unregister(l.dataIn)
}

func (l *Listener) Listen() <-chan interface{} {
	return l.dataIn
}

func NewDispatch(t chan Telemetry, hData chan Telemetry, dictData chan Telemetry) Dispatcher {
	arbitraryBufLen := 10

	d := Dispatcher{t, broadcast.NewBroadcaster(arbitraryBufLen), hData, dictData}

	go d.run()

	return d
}

func (d *Dispatcher) NewListener() Listener {
	ch := make(chan interface{})

	d.telemetryIn.Register(ch)

	return Listener{ch, d.telemetryIn}
}

func (d *Dispatcher) run() {
	var telem Telemetry

	for telem = range d.TelemIn {
		d.hData <- telem
		d.dictData <- telem
		// fmt.Println(telem)
		d.telemetryIn.Send(telem)
	}
}
