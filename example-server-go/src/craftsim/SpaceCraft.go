package craftsim

import (
	// "strings"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"math"
	"math/rand"
	"time"
)

// type Values struct {
// 	Key    string `json:"key"`
// 	Name   string `json:"name"`
// 	Units  string `json:"units,omitempty"`
// 	Format string `json:"format"`
// 	Min    int    `json:"min,omitempty"`
// 	Max    int    `json:"max,omitempty"`
// 	Hints  struct {
// 		Range int `json:"range"`
// 	} `json:"hints"`
// 	Source string `json:"source,omitempty"`
// } 

// type InstrumentData struct {
// 	Name   string `json:"name"`
// 	Key    string `json:"key"`
// 	Values []Values `json:"values"`
// }

// type SpaceCraft struct {
// 	Name           string `json:"name"`
// 	Key            string `json:"key"`
// 	InstrumentData []InstrumentData `json:"measurements"`
// }

type SpaceCraft struct {
	PropFuel      int     `json:"prop.fuel"`
	PropThrusters string  `json:"prop.thrusters"`
	CommsRecd     int     `json:"comms.recd"`
	CommsSent     int     `json:"comms.sent"`
	PwrTemp       int     `json:"pwr.temp"`
	PwrC          float64 `json:"pwr.c"`
	PwrV          int     `json:"pwr.v"`
}

func InitSpacecraft() {
	sp := SpaceCraft {
		77,    // PropFuel
		"OFF", // PropThrusters
		0,     // CommsRecd
		0,     // CommsSent
		245,   // PwrTemp
		8.15,  // PwrC
		30     // PwrV
	}

	return sp
}


func LoadCraftJSON(craftJSON string) SpaceCraft {
	dictionary, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("error:", err)
	}
	
	var sp SpaceCraft
	err = json.Unmarshal(dictionary, &sp)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(sp)
	
	return sp
}

func reduceFuel(sp SpaceCraft) {
	if (sp.PropThrusters == "OFF") {
		return
	}

	if (sp.PropFuel <= 0) {
		sp.PropFuel = sp.PropFuel - .5
	} else {
		sp.PropFuel = 0
	}
}

func updateState(sp SpaceCraft) {
	reduceFuel(&sp)
	
	sp.PwrTemp = (sp.PwrTemp * .985) + (math.Float64 * .25) + math.Sin(time.Now())

	if (sp.PropThrusters == "ON") {
		sp.PwrC = 8.15
	} else {
		sp.PwrC *= .985
	}

	sp.PwrV = 30 + (rand.Float64 ** 3)
}

