package server

import (
	"broadcast"
	"errors"
)

// ============================================================================
// Telemetry-specific Dispatcher
// ============================================================================

type Dispatcher struct {
	telemIn chan Telemetry

	telemetryIn broadcast.Broadcaster

	hData    chan<- Telemetry
	dictData chan<- Telemetry

	switcher map[string]*broadcast.Broadcaster
}

func NewDispatch(t chan Telemetry, hData chan Telemetry, dictData chan Telemetry) Dispatcher {
	arbitraryBufLen := 10

	d := Dispatcher{
		telemIn:     t,
		telemetryIn: broadcast.NewBroadcaster(arbitraryBufLen),
		hData:       hData,
		dictData:    dictData,
	}

	go d.run()

	return d
}

func (d *Dispatcher) run() {
	var telem Telemetry

	for telem = range d.telemIn {
		d.hData <- telem
		d.dictData <- telem

		d.telemetryIn.Send(telem)
	}
}
func (r *receiver) Receive() {

}

type Listener struct {
	dataIn chan interface{}
	b      broadcast.Broadcaster
}

func (l *Listener) CloseListener() {
	l.b.Unregister(l.dataIn)
}

func (l *Listener) Listen() <-chan interface{} {
	return l.dataIn
}

func (d *Dispatcher) NewListener() Listener {
	ch := make(chan interface{})

	d.telemetryIn.Register(ch)

	return Listener{ch, d.telemetryIn}
}

// ============================================================================
// Type-generic data dispatcher
// ============================================================================

// generic type for dispatcher args above
type dispatcher struct {
	dataIn chan interface{}

	alwaysSend []chan interface{}
	switcher   map[interface{}]broadcast.Broadcaster
	key        KeyGetter
}

type KeyGetter interface {
	Get(interface{}) (interface{}, error)
}

type receiver struct {
	listening map[interface{}]chan<- interface{}

	key KeyGetter
}

type Receiver interface {
	Receive() interface{}
	Subscribe(interface{})
	Unsubscribe(interface{})
}

func (d *dispatcher) run() {
	var i interface{}

	for i = range d.dataIn {
		d.Dispatch(i)
	}
}

func (d *dispatcher) Dispatch(i interface{}) error {

	key, err := d.key.Get(i)
	if err != nil {
		return errors.New("Failed to retrieve key from value")
	}

	for _, datachan := range d.alwaysSend {
		datachan <- i
	}

	arbitraryBufLen := 10

	if _, ok := d.switcher[key]; !ok {
		d.switcher[key] = broadcast.NewBroadcaster(arbitraryBufLen)
	}

	d.switcher[key].Send(i)

	return nil
}
