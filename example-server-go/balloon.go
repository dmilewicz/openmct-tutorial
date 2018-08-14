package main

import (
	"fmt"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "12345"
	CONN_TYPE = "tcp"
)

func main() {

	db, _ := leveldb.OpenFile("db/leveldb1", nil)

	// for i := 0; i < 10; i++ {
	// 	intstr := strconv.FormatInt(int64(i), 10)
	// 	db.Put([]byte(intstr), []byte(intstr), nil)
	// }

	// for i := 0; i < 10; i++ {
	// 	intstr := strconv.FormatInt(int64(i), 10)
	// 	data, err := db.Get([]byte(intstr), nil)

	// 	fmt.Println(string(data))
	// 	fmt.Println(err)

	// }

	r := util.Range{
		Start: []byte(strconv.FormatInt(int64(4), 10)),
		Limit: []byte(strconv.FormatInt(int64(7), 10)),
	}

	it := db.NewIterator(&r, nil)

	for it.Next() {
		fmt.Println(string(it.Value()))
	}

	// var f, g interface{}
	// f = "hello"
	// g = "hells"

	// h := make(map[interface{}]interface{})

	// h[f] = 0

	// if val, ok := h[g]; ok {
	// 	fmt.Println("val good: ", val)
	// } else {
	// 	fmt.Println("nope")
	// }

}
