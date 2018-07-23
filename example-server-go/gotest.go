package main

import "fmt"
import "reflect"
import "time"

type SpaceCraft struct {
	PropFuel      int     `json:"prop.fuel"`
	PropThrusters string  `json:"prop.thrusters"`
	CommsRecd     int     `json:"comms.recd"`
	CommsSent     int     `json:"comms.sent"`
	PwrTemp       int     `json:"pwr.temp"`
	PwrC          float64 `json:"pwr.c"`
	PwrV          int     `json:"pwr.v"`
}

func main() {
	messages := make(chan string, 2)
	// go func() { messages <- "pong" ; messages <- pong}()

	// messages <- "buffered"
	// messages <- "channel"
	// go func() { messages <- "ping" }()

	// msg := <-messages
	// fmt.Println(msg)
	// fmt.Println(<-messages)
	// messages <- "third"

	// msg = <-messages
	// fmt.Println(msg)

	go lister("1", messages)
	go lister("2", messages)

	messages <- "heloo"
	// fmt.Println(reflect.TypeOf(msg))

	x := SpaceCraft{
		77,    // PropFuel
		"OFF", // PropThrusters
		0,     // CommsRecd
		0,     // CommsSent
		245,   // PwrTemp
		8.15,  // PwrC
		30 /* PwrV */}

	v := reflect.ValueOf(x)

	values := make([]interface{}, v.NumField())

	fmt.Println(reflect.TypeOf(v))

	for i := 0; i < v.NumField(); i++ {
		fmt.Println(v.Type().Field(i).Tag.Get("json"))
	}
	fmt.Println(values)

	fmt.Println(time.Now())
}

func lister(name string, msg chan string) {
	fmt.Println(name + " " + <-msg)
}
