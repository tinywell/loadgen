package lib

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
	"tinywell/loadgen/model"
)

type ParamSet struct {
	caller     Caller
	tickets    Tickets
	lps        uint32
	timeoutNS  time.Duration
	durationNS time.Duration
	resultChan chan *model.LDResult
}

// myGenerator 载荷发生器实现
type myGenerator struct {
	caller      Caller
	tickets     Tickets
	ctx         context.Context
	cancelFunc  context.CancelFunc
	resultChan  chan *model.LDResult
	concurrency int32
	lps         uint32
	durationNS  time.Duration
	timeoutNS   time.Duration
	status      uint32
}

func NewGenerator(param ParamSet) Generator {
	gen := &myGenerator{
		caller:     param.caller,
		tickets:    param.tickets,
		resultChan: param.resultChan,
		status:     GEN_STA_ORIGIN,
		lps:        param.lps,
		timeoutNS:  param.timeoutNS,
		durationNS: param.durationNS,
	}
	return gen
}

// 初始化 gen
func (gen *myGenerator) init() error {

	gen.concurrency = (int32)((int64)(gen.lps) * (gen.timeoutNS.Nanoseconds() / 1e9))

	tickets, err := NewTickets(gen.concurrency)
	if err != nil {
		return err
	}
	gen.tickets = tickets

	return nil
}

func (gen *myGenerator) Start() bool {
	fmt.Println("starting generator")
	if ok := atomic.CompareAndSwapUint32(&gen.status, GEN_STA_ORIGIN, GEN_STA_STARTING); !ok {
		return false
	}

	err := gen.init()
	if err != nil {
		return false
	}

	gen.ctx, gen.cancelFunc =
		context.WithTimeout(context.Background(), gen.durationNS)

	var throttle <-chan time.Time
	if gen.lps > 0 {
		interval := time.Duration(1e9 / gen.lps)
		throttle = time.Tick(interval)
	}

	go func() {
		gen.genLoad(throttle)
	}()

	if ok := atomic.CompareAndSwapUint32(&gen.status, GEN_STA_STARTING, GEN_STA_STARTED); !ok {
		return false
	}
	return true
}

func (gen *myGenerator) genLoad(throttle <-chan time.Time) {
	for {
		select {
		case <-gen.ctx.Done():
			gen.Stop()
		default:
		}
		gen.asyncCall()
		select {
		case <-throttle:
		case <-gen.ctx.Done():
			gen.Stop()
		}
	}
}

// 发起一次调用
func (gen *myGenerator) asyncCall() {

}

func (gen *myGenerator) Stop() bool {
	return true
}

func (gen *myGenerator) prepareToStop() {

}

func (gen *myGenerator) Status() uint32 {
	return gen.status
}
