package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"strconv"

)


type HistoryServer struct {
	Telem  chan interface{}
	PubSub chan string
}

type DataRequest struct {
	Value string
	Start int64
	End   int64
}

func extractRequest(u *url.URL) DataRequest {
	q := u.Query()

	start, _ := strconv.ParseInt(q["start"][0], 10, 64)
	end, _   := strconv.ParseInt(q["end"][0], 10, 64)

	dr := DataRequest {
		strings.TrimPrefix(u.Path, "/history/"),
		start,
		end }
	
	return dr
}


func (hs *HistoryServer) RunServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("History server starting...")

	dr := extractRequest(r.URL)


	fmt.Println(dr)


}

