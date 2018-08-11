package main

import (
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "12345"
	CONN_TYPE = "tcp"
)

func main() {

	db, err := leveldb.OpenFile("db/leveldb1", nil)

	for i := 0; i < 10; i++ {
		intstr := strconv.FormatInt(int64(i), 10)
		db.Put([]byte(intstr), []byte(intstr))
	}

	db.

}
