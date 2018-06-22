package server

import (
	"reflect"
	"time"
	// "net"
)

type Datum struct {
	Timestamp int64       `json:"timestamp"`
	Value     interface{} `json:"value"`
	ID        string      `json:"id"`
}

type Telemetry struct {
	value     interface{}
	idx       map[string]int
	Timestamp int64
}

func LoadTelemetry(v interface{}) Telemetry {
	val := reflect.ValueOf(v)

	t := Telemetry{v, make(map[string]int), time.Now().UnixNano() / int64(time.Millisecond)}

	for i := 0; i < val.Type().NumField(); i++ {
		tag := val.Type().Field(i).Tag.Get("json")

		t.idx[tag] = i
	}

	return t
}

func (t Telemetry) Get(key string) interface{} {
	tval := reflect.ValueOf(t)
	return tval.Field(t.idx[key]).Interface()
}

func (t Telemetry) Datum(name string) Datum {
	return Datum{
		t.Timestamp,
		t.Get(name),
		name}
}
