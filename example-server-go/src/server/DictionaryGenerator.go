package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"labix.org/v2/mgo/bson"
)

// ============================================================================
// Data Structures
// ============================================================================

type DictionaryGenerator struct {
	dataIn <-chan TelemetryBuffer

	packages   map[string]Telemetry
	dictString string
}

type OpenMCTTime time.Time

type TelemetryBuffer struct {
	Name      string      `json:"name" bson:"name"`
	Key       string      `json:"-" bson:"key"`
	Flags     int64       `json:"flags,omitempty" bson:"flags,omitempty"`
	Timestamp OpenMCTTime `json:"timestamp" bson:"timestamp"`
	Raw_Type  int64       `json:"-" bson:"raw_type"`
	Raw_Value interface{} `json:"raw_value" bson:"raw_value"`
	Eng_Type  int64       `json:"eng_type,omitempty" bson:"eng_type,omitempty"`
	Eng_Val   interface{} `json:"eng_val,omitempty" bson:"eng_val,omitempty"`
}

type Val struct {
	Key      string `json:"key" bson:"key"`
	Name     string `json:"name" bson:"name"`
	Units    string `json:"units,omitempty" bson:"units,omitempty"`
	Format   string `json:"format,omitempty" bson:"format,omitempty"`
	Min      int    `json:"min,omitempty" bson:"min,omitempty"`
	Max      int    `json:"max,omitempty" bson:"max,omitempty"`
	Raw_Type int64  `json:"raw_type,omitempty" bson:"raw_type,omitempty"`
	Eng_Type int64  `json:"eng_type,omitempty" bson:"eng_type,omitempty"`
	Hints    hint   `json:"hints" bson:"hints"`
	Source   string `json:"source,omitempty"`
}

type hint struct {
	Range  int `json:"range,omitempty" bson:"range,omitempty"`
	Domain int `json:"domain,omitempty" bson:"domain,omitempty"`
}

type Telemetry struct {
	Name   string `json:"name"`
	Key    string `json:"key"`
	ID     string `json:"id"`
	Values []Val  `json:"values"`
}

type Package struct {
}

// ============================================================================
// Member Functions
// ============================================================================

func NewDictionaryGenerator(dataIn <-chan TelemetryBuffer) DictionaryGenerator {
	d := DictionaryGenerator{
		dataIn:   dataIn,
		packages: make(map[string]Telemetry),
	}

	// read in the points file?

	go d.runGenerator()

	return d
}

func (dg *DictionaryGenerator) runGenerator() {
	for point := range dg.dataIn {
		dg.writePoint(point)
	}
}

func (dg *DictionaryGenerator) writePoint(point TelemetryBuffer) {
	if _, ok := dg.packages[point.Name]; !ok {
		pointMetaData := UnmarshalTelemetry(point)

		dg.packages[point.Name] = pointMetaData

		dg.Save()

	}

}

func (dg *DictionaryGenerator) Save() {
	s, err := json.Marshal(dg.packages)
	// s, err := json.MarshalIndent(dg.packages, "", "\t")

	if err != nil {
		panic("ERROR")
	}

	ioutil.WriteFile("points.json", s, 0644) // TODO: what is the black magic?
}

func (t *Telemetry) String() string {
	// s, err := json.MarshalIndent(*t, "", "\t")
	s, err := json.Marshal(*t)

	if err != nil {
		panic("ERROR")
	}
	return string(s)
}

// ============================================================================
// Unmarshalling Functions
// ============================================================================

func (t OpenMCTTime) MarshalJSON() ([]byte, error) {
	// fmt.Println(time.Time(t))
	// return time.Time(t).MarshalJSON()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(time.Time(t).UnixNano()/int64(time.Millisecond)))
	// return b, nil
	return []byte(fmt.Sprintf("%d", time.Time(t).UnixNano()/int64(time.Millisecond))), nil

}

func (t *OpenMCTTime) SetBSON(b bson.Raw) error {
	bt := new(time.Time)
	err := b.Unmarshal(&bt)
	*t = OpenMCTTime(*bt)
	return err
}

func UnmarshalTelemetry(tBuf TelemetryBuffer) Telemetry {
	var vals []Val

	vals = append(vals, Val{
		Key:    "utc",
		Source: "timestamp",
		Name:   "Timestamp",
		Format: "utc",
		Hints: hint{
			Domain: 1,
		},
	},
		Val{
			Key:      "raw_value",
			Name:     "Raw Value",
			Raw_Type: tBuf.Raw_Type,
			Hints: hint{
				Range: 2,
			},
		})

	if tBuf.Eng_Val != nil {
		vals = append(vals, Val{
			Key:      "eng_val",
			Name:     "Engineering Value",
			Eng_Type: tBuf.Eng_Type,
			Hints: hint{
				Range: 1,
			},
		})
	}

	return Telemetry{
		Name:   tBuf.Name,
		Key:    tBuf.Name,
		ID:     tBuf.Name,
		Values: vals,
	}
}
