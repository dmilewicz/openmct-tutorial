package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type HistoryServer struct {
	Request     chan DataRequest
	HistoryData chan []Telemetry
}

type DataRequest struct {
	Value string
	Start int64
	End   int64
}

func extractRequest(u *url.URL) DataRequest {
	q := u.Query()

	// TODO: Decide to keep timestamp as int or float?
	start, _ := strconv.ParseFloat(q.Get("start"), 64)
	end, _ := strconv.ParseFloat(q.Get("end"), 64)

	dr := DataRequest{
		strings.TrimPrefix(u.Path, "/history/"),
		int64(start),
		int64(end)}

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
