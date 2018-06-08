package server

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
)

type Values struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Units  string `json:"units,omitempty"`
	Format string `json:"format"`
	Min    int    `json:"min,omitempty"`
	Max    int    `json:"max,omitempty"`
	Hints  struct {
		Range int `json:"range"`
	} `json:"hints"`
	Source string `json:"source,omitempty"`
}


type InstrumentData struct {
	Name   string `json:"name"`
	Key    string `json:"key"`
	Values []Values `json:"values"`
}


type SpaceCraft struct {
	Name           string `json:"name"`
	Key            string `json:"key"`
	InstrumentData []InstrumentData `json:"measurements"`
}


func main() {
	var sp SpaceCraft
	dictionary, err := ioutil.ReadFile("../dictionary.json")
	if err != nil {
		fmt.Println("error:", err)
	}

	// fmt.Printf("File contents: %s", d)

	err = json.Unmarshal(dictionary, &sp)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", sp)
}
