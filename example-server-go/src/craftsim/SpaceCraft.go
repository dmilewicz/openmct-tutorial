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
	PropFuel      float64 `json:"prop.fuel"`
	PropThrusters string  `json:"prop.thrusters"`
	CommsRecd     int     `json:"comms.recd"`
	CommsSent     int     `json:"comms.sent"`
	PwrTemp       float64     `json:"pwr.temp"`
	PwrC          float64 `json:"pwr.c"`
	PwrV          float64     `json:"pwr.v"`
}

func (sp *SpaceCraft) InitSpacecraft() {
	sp.PropFuel = 77
	sp.PropThrusters = "OFF"
	sp.CommsRecd = 0
	sp.CommsSent = 0
	sp.PwrTemp = 245
	sp.PwrC = 8.15
	sp.PwrV = 30
}

func NewSpacecraft() *SpaceCraft {
	return &SpaceCraft {
		77,    // PropFuel
		"OFF", // PropThrusters
		0,     // CommsRecd
		0,     // CommsSent
		245,   // PwrTemp
		8.15,  // PwrC
		30     /* PwrV */  }
}




func LoadCraftJSON(filename string) SpaceCraft {
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
/*
 * Launch that sucker
 */
func (sp *SpaceCraft) Launch() {
	sp.PropThrusters = "ON"
}

func (sp *SpaceCraft) reduceFuel() {
	if (sp.PropThrusters == "OFF") {
		return
	}

	if (sp.PropFuel <= 0) {
		sp.PropFuel = sp.PropFuel - .5
	} else {
		sp.PropFuel = 0
	}
}

func (sp *SpaceCraft) updateState() {
	sp.reduceFuel()
	
	sp.PwrTemp = (sp.PwrTemp * .985) + (rand.Float64() * .25) + math.Sin(float64(time.Now().UnixNano() / 1000000))

	if (sp.PropThrusters == "ON") {
		sp.PwrC = 8.15
	} else {
		sp.PwrC *= .985
	}

	sp.PwrV = 30 + math.Pow(rand.Float64(), 3.00)
}

func (sp *SpaceCraft) RunSim() {
	for {
		<-time.After(2 * time.Second)
		
		sp.updateState()
		// fmt.Println(sp)
	  }
}


