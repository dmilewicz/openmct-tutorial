package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type HistoryServer struct {
	Request     chan DataRequest
	HistoryData chan []Telemetry
}

type DataRequest struct {
	Value string
	Start time.Time
	End   time.Time
}

func NewHistoryServer() HistoryServer {
	return HistoryServer{
		Request:     make(chan DataRequest),
		HistoryData: make(chan []Telemetry),
	}
}

func extractRequest(u *url.URL) DataRequest {
	q := u.Query()

	// TODO: Decide to keep timestamp as int or float?
	start, _ := strconv.ParseFloat(q.Get("start"), 64)
	end, _ := strconv.ParseFloat(q.Get("end"), 64)

	dr := DataRequest{
		strings.TrimPrefix(u.Path, "/history/"),
		time.Unix(int64(start)/1000, 0),
		time.Unix(int64(end)/1000, 0)}

	return dr
}

func (hs *HistoryServer) RunServer(w http.ResponseWriter, r *http.Request) {
	// get request details
	dr := extractRequest(r.URL)

	// send data request and block for response
	hs.Request <- dr
	resp := <-hs.HistoryData

	// put response in json format
	res, _ := json.Marshal(resp)

	// set as JSON response
	w.Header().Add("Content-Type", "application/json")

	w.Write(res)
}
