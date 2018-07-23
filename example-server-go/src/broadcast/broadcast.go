package broadcast

// import (
// 	"fmt"
// )

// type broadcast struct {
// 	c chan broadcast
// 	v interface{}
// }

// type Broadcaster struct {
// 	// private fields:
// 	Listenc chan chan (chan broadcast)
// 	Sendc   chan<- interface{}
// }

// type Receiver struct {
// 	// private fields:
// 	C chan broadcast
// }

// // create a new broadcaster object.
// func NewBroadcaster() Broadcaster {
// 	listenc := make(chan (chan (chan broadcast)))
// 	sendc := make(chan interface{})
// 	go func() {
// 		currc := make(chan broadcast)
// 		for {
// 			select {
// 			case v := <-sendc:
// 				if v == nil {
// 					currc <- broadcast{}
// 					return
// 				}
// 				// fmt.Print("hello")
// 				// fmt.Println(v)

// 				c := make(chan broadcast)
// 				b := broadcast{c: c, v: v}
// 				select {
// 				case currc <- b:
// 					fmt.Println("yep ", v)
// 					currc = c
// 				default:
// 					fmt.Println("no message")
// 				}
// 			case r := <-listenc:
// 				fmt.Println("listenc")
// 				r <- currc
// 			}
// 		}
// 	}()
// 	return Broadcaster{
// 		Listenc: listenc,
// 		Sendc:   sendc,
// 	}
// }

// // start listening to the broadcasts.
// func (b Broadcaster) Listen() Receiver {
// 	c := make(chan chan broadcast, 0)
// 	b.Listenc <- c
// 	return Receiver{<-c}
// }

// // broadcast a value to all listeners.
// func (b Broadcaster) Write(v interface{}) {

// 	b.Sendc <- v
// 	fmt.Println("done writing")
// }

// var count = 1

// // read a value that has been broadcast,
// // waiting until one is available if necessary.
// func (r *Receiver) Read() interface{} {
// 	b := <-r.C
// 	v := b.v
// 	fmt.Println("count: ", count)
// 	count++
// 	r.C <- b
// 	r.C = b.c
// 	return v
// }

// func (r *Receiver) Destroy() {
// 	b := r.C
// }
