package server

// "net"

// type Telemetry struct {
// 	value     interface{}
// 	idx       map[string]int
// 	Timestamp int64
// }

// type Telemetry interface

/*
 * Loads a telemetry object for
 */

// func UnmarshalTelemetry(buf []byte) Telemetry {
// 	var tBuf TelemetryBuffer

// 	bson.Unmarshal(buf, tBuf)

// 	var vals []Val

// 	append(vals, Val{
// 		Key:    "utc",
// 		Source: "timestamp",
// 		Name:   "Timestamp",
// 		Format: "utc",
// 		Hints: {
// 			domain: 1,
// 		},
// 	},
// 		Val{
// 			Key:      "raw_value",
// 			Name:     "Raw Value",
// 			Raw_Type: tBuf.Raw_Type,
// 			Hints: {
// 				Range: 2,
// 			},
// 		})

// 	if tBuf.Eng_Val != nil {
// 		append(vals, Val{
// 			Key:      "eng_val",
// 			Name:     "Engineering Value",
// 			Eng_Type: tBuf.Eng_Type,
// 			Hints: {
// 				Range: 1,
// 			},
// 		})
// 	}

// 	t := Telemetry{
// 		Name:   tBuf.Name,
// 		Key:    tBuf.Key,
// 		Values: vals,
// 	}

// }

// // func BuildValues() {

// }

// type packages struct {
// 	name   string   `json:"name"`
// 	points []string `json:"points"`
// }

// type TelemetryBuffer struct {
// 	Name      string      // `json:"name" bson:"name"`
// 	Key       string      // `json:"key" bson:"key"`
// 	Flags     int64       // `json:"flags,omitempty" bson:"flags,omitempty"`
// 	Timestamp float64     // `json:"timestamp" bson:"timestamp"`
// 	Raw_Type  int64       // `json:"raw_type" bson:"raw_type"`
// 	Raw_Value interface{} // `json:"raw_value" bson:"raw_value"`
// 	Eng_Type  int64       // `json:"eng_type,omitempty" bson:"eng_type,omitempty"`
// 	Eng_Val   interface{} // `json:"eng_val,omitempty" bson:"eng_val,omitempty"`
// }

// type Telemetry struct {
// 	Name   string `json:"name"`
// 	Key    string `json:"key"`
// 	Values []Val  `json:"values"`
// 	// []struct {
// 	// 	Key     string `json:"key" bson:"key"`
// 	// 	Name    string `json:"name" bson:"name"`
// 	// 	Units   string `json:"units,omitempty" bson:"units,omitempty"`
// 	// 	Format  string `json:"format,omitempty" bson:"format,omitempty"`
// 	// 	Min     int    `json:"min,omitempty" bson:"min,omitempty"`
// 	// 	Max     int    `json:"max,omitempty" bson:"max,omitempty"`
// 	// 	RawType int32  `json:"raw_type,omitempty" bson:"raw_type,omitempty"`
// 	// 	Hints   hint   `json:"hints" bson:"hints"`
// 	// 	Source  string `json:"source,omitempty"`
// 	// } `json:"values"`
// }

// type Val struct {
// 	Key    string `json:"key" bson:"key"`
// 	Name   string `json:"name" bson:"name"`
// 	Units  string `json:"units,omitempty" bson:"units,omitempty"`
// 	Format string `json:"format,omitempty" bson:"format,omitempty"`
// 	// Min     int    `json:"min,omitempty" bson:"min,omitempty"`
// 	// Max     int    `json:"max,omitempty" bson:"max,omitempty"`
// 	EngType int32 `json:"raw_type,omitempty" bson:"raw_type,omitempty"`

// 	RawType int32  `json:"raw_type,omitempty" bson:"raw_type,omitempty"`
// 	Hints   hint   `json:"hints" bson:"hints"`
// 	Source  string `json:"source,omitempty"`
// }

// type hint struct {
// 	Range  int `json:"range,omitempty" bson:"range,omitempty"`
// 	Domain int `json:"domain,omitempty" bson:"domain,omitempty"`
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
