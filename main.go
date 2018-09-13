package main

import (
	"fmt"
	"time"
	"tinywell/loadgen/lib"
	"tinywell/loadgen/mock"
	"tinywell/loadgen/model"
)

func main() {
	caller := mock.NewMockCaller()
	tickets, _ := lib.NewTickets(2)
	set := lib.ParamSet{
		Caller:     caller,
		Tickets:    tickets,
		Lps:        10,
		TimeoutNS:  time.Second * 2,
		DurationNS: time.Second * 10,
		ResultChan: make(chan *model.LDResult),
	}
	gen := lib.NewGenerator(set)
	if ok := gen.Start(); ok {
		fmt.Println("ok")
	}
	time.Sleep(time.Second * 2)
	gen.Stop()
}
