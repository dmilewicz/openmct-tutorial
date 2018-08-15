package main

import (

	// "strconv"
	"fmt"
	"reflect"
	// "github.com/syndtr/goleveldb/leveldb"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "12345"
	CONN_TYPE = "tcp"
)

type hdh struct {
	f int
	g uint


}

func main() {

	s := &hdh{}

	fmt.Println(reflect.ValueOf(s).Elem().Type())
}
