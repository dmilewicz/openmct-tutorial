package server

import (
	"broadcast"
	"errors"
	"fmt"
	"sync"
)

// ============================================================================
// Telemetry-specific Dispatcher
// ============================================================================

type TelemetryDispatcher struct {
	telemIn chan Telemetry

	telemetryIn broadcast.Broadcaster

	hData    chan<- Telemetry
	dictData chan<- Telemetry

	switcher map[string]*broadcast.Broadcaster
}

func NewDispatch(t chan Telemetry, hData chan Telemetry, dictData chan Telemetry) TelemetryDispatcher {
	arbitraryBufLen := 10

	d := TelemetryDispatcher{
		telemIn:     t,
		telemetryIn: broadcast.NewBroadcaster(arbitraryBufLen),
		hData:       hData,
		dictData:    dictData,
	}

	go d.run()

	return d
}

func (d *TelemetryDispatcher) run() {
	var telem Telemetry

	for telem = range d.telemIn {
		d.hData <- telem
		d.dictData <- telem

		d.telemetryIn.Send(telem)
	}
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

func (d *TelemetryDispatcher) NewListener() Listener {
	ch := make(chan interface{})

	d.telemetryIn.Register(ch)

	return Listener{ch, d.telemetryIn}
}

type TelemID bool

func (t TelemID) Get(i interface{}) (interface{}, error) {
	return i.(Telemetry).Name, nil
}

// ============================================================================
// Type-generic data dispatcher
// ============================================================================

// generic type for dispatcher args above
type dispatcher struct {
	dataIn chan interface{}

	alwaysSend []chan Telemetry
	msgMux     *sync.Map
	key        KeyGetter
}

func NewTelemetryDispatcher(dataIn chan interface{}, kg KeyGetter, senders ...chan Telemetry) dispatcher {
	d := dispatcher{
		dataIn:     dataIn,
		alwaysSend: senders,
		msgMux:     new(sync.Map),
		key:        kg,
	}

	go d.run()
	return d
}

func (d *dispatcher) AddAlwaysSend(c chan Telemetry) {
	d.alwaysSend = append(d.alwaysSend, c)
}

func (d *dispatcher) run() {
	var i interface{}

	for {
		select {
		case i = <-d.dataIn:
			d.Dispatch(i)
		}
	}
}

func (d *dispatcher) NewReceiver() Receiver {
	r := receiver{
		dataOut:  make(chan interface{}),
		switcher: d.msgMux,
		key:      d.key,
	}

	return &r
}

func (d *dispatcher) Dispatch(i interface{}) error {

	key, err := d.key.Get(i)
	if err != nil {
		return errors.New("Failed to retrieve key from value")
	}

	for _, datachan := range d.alwaysSend {
		datachan <- i.(Telemetry)
	}

	arbitraryBufLen := 10

	b, ok := d.msgMux.Load(key)
	for !ok {
		d.msgMux.Store(key, broadcast.NewBroadcaster(arbitraryBufLen))
		b, ok = d.msgMux.Load(key)
	}

	b.(broadcast.Broadcaster).Send(i)

	return nil
}

type KeyGetter interface {
	Get(interface{}) (interface{}, error)
}

type receiver struct {
	dataOut     chan interface{}
	unsubscribe sync.Map
	switcher    *sync.Map

	key KeyGetter
}

type Receiver interface {
	Receive() <-chan interface{}
	Subscribe(interface{})
	Unsubscribe(interface{})
}

func (r *receiver) Receive() <-chan interface{} {
	return r.dataOut
}

func (r *receiver) Subscribe(key interface{}) {
	if _, ok := r.unsubscribe.Load(key); ok {
		// already have channel listening
		return
	}

	fmt.Println(r.switcher)
	b, ok := r.switcher.Load(key)
	for !ok {
		r.switcher.Store(key, broadcast.NewBroadcaster(10))
		b, ok = r.switcher.Load(key)
	}

	brdcstr := b.(broadcast.Broadcaster)

	dataChan := make(chan interface{}, 25)
	closerChan := make(chan bool)

	brdcstr.Register(dataChan)
	r.unsubscribe.Store(key, closerChan)

	go func(key interface{}) {
		var i interface{}

		for {
			select {
			case i = <-dataChan:
				r.dataOut <- i
			case <-closerChan:
				brdcstr.Unregister(dataChan)
				r.unsubscribe.Delete(key)
			}
		}
	}(key)
}

func (r *receiver) Unsubscribe(key interface{}) {
	if cChan, ok := r.unsubscribe.Load(key); ok {
		cChan.(chan bool) <- true
	}
}

// func CopyMsgMux(m msgMux) msgMux {
// 	var newM msgMux
// 	newM = make(map[interface{}]broadcast.Broadcaster)

// 	for k, v := range m {
// 		newM[k] = v
// 	}

// 	return newM
// }
