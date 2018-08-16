package main

import (
	// "strconv"
	"fmt"
	// "github.com/syndtr/goleveldb/leveldb"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "12345"
	CONN_TYPE = "tcp"
)

type ha map[string]string

func nums(n ...int, f ...string) {
	for _, v := range f {
		fmt.Println(v)
	}
}

func main() {
	nums(1, 2, 3)
}
