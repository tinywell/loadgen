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
		Lps:        200,
		TimeoutNS:  time.Second * 2,
		DurationNS: time.Second * 10,
		ResultChan: make(chan *model.LDResult, 10),
	}
	gen := lib.NewGenerator(set)
	if ok := gen.Start(); ok {
		fmt.Println("ok")
	}
	readRes(set.ResultChan)
	// time.Sleep(time.Second * 1)
	// gen.Stop()
}

func readRes(resChan chan *model.LDResult) {
	for res := range resChan {
		fmt.Println(res.ID, ":", res)
	}
}
