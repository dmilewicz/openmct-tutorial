package server

import (
	"encoding/binary"
	"fmt"
	"time"

	"labix.org/v2/mgo/bson"
)

type point struct {
}

// type packages struct {
// 	name   string   `json:"name"`
// 	points []string `json:"points"`
// }

// type Telemetry struct {
// 	Values []value `json:"values"`
// }

type OpenMCTTime time.Time

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

// type value struct {
// 	Key     string `json:"key" bson:"key"`
// 	Name    string `json:"name" bson:"name"`
// 	Units   string `json:"units,omitempty" bson:"units,omitempty"`
// 	Format  string `json:"format,omitempty" bson:"format,omitempty"`
// 	Min     int    `json:"min,omitempty" bson:"min,omitempty"`
// 	Max     int    `json:"max,omitempty" bson:"max,omitempty"`
// 	RawType int32  `json:"raw_type,omitempty" bson:"raw_type,omitempty"`
// 	Hints   struct {
// 		Range  int `json:"range" bson:"range"`
// 		Domain int `json:"domain,omitempty" bson:"domain,omitempty"`
// 	} `json:"hints"`
// 	Source string `json:"source,omitempty"`
// }

// type Dictionary struct {
// 	nameseen map[string]bool
// 	// chan dataIn

// }

// func LoadTelemetry(v interface{}) Telemetry {
// 	val := reflect.ValueOf(v)

// 	t := Telemetry{v, make(map[string]int), time.Now().UnixNano() / int64(time.Millisecond)}

// 	for i := 0; i < val.Type().NumField(); i++ {
// 		tag := val.Type().Field(i).Tag.Get("json")

// 		t.idx[tag] = i
// 	}

// 	return t
// }

// func (t Telemetry) GetIdx(key string) int {
// 	return t.idx[key]
// }

// func (t Telemetry) Get(key string) interface{} {
// 	tval := reflect.ValueOf(t.value)
// 	return tval.Field(t.idx[key]).Interface()
// }

// func (t Telemetry) Datum(name string) Datum {
// 	return Datum{
// 		t.Timestamp,
// 		t.Get(name),
// 		name}
// }

// func (t Telemetry) Len() int {
// 	return reflect.ValueOf(t.value).Type().NumField()
// }
