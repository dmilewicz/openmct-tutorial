package server

import (
	"fmt"
	"time"
)

const MB = 1024 * 1024

type Counter struct {
	title string

	frameSeconds int
	frameCounter int
	countInFrame int
	totalCount   int64
}

// NewCounter()

func (c *Counter) run() {

	for {
		<-time.After(time.Duration(c.frameSeconds) * time.Second)
		c.totalCount += int64(c.countInFrame)
		c.report()
		c.countInFrame = 0
		c.frameCounter++
	}
}

func (c *Counter) Add(count int) {

	c.countInFrame += count

}

func (c *Counter) report() {
	fmt.Println(c.title+" frame rate: ", c.countInFrame/c.frameSeconds, " points/sec")
}
