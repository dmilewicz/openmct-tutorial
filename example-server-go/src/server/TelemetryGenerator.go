package server

import (
	"time"
	// "net"
)

type Telemetry struct {
	Timestamp int64           `json:"timestamp"`
	Value     interface{}     `json:"value"`
	ID        string          `json:"id"`
}


func MakeTelemetry(name string, val interface{}) Telemetry {
	telem := Telemetry {
		time.Now().UnixNano() / 100000,
		val,
		name }

	return telem
}





















