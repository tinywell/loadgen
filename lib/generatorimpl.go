package lib

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
	"tinywell/loadgen/model"
)

type ParamSet struct {
	Caller     model.Caller
	Tickets    model.Tickets
	Lps        uint32
	TimeoutNS  time.Duration
	DurationNS time.Duration
	ResultChan chan *model.LDResult
}

// myGenerator 载荷发生器实现
type myGenerator struct {
	caller      model.Caller
	tickets     model.Tickets
	ctx         context.Context
	cancelFunc  context.CancelFunc
	resultChan  chan *model.LDResult
	concurrency int32
	lps         uint32
	durationNS  time.Duration
	timeoutNS   time.Duration
	status      uint32
	stopChan    chan struct{}
}

func NewGenerator(param ParamSet) model.Generator {
	gen := &myGenerator{
		caller:     param.Caller,
		tickets:    param.Tickets,
		resultChan: param.ResultChan,
		status:     model.GEN_STA_ORIGIN,
		lps:        param.Lps,
		timeoutNS:  param.TimeoutNS,
		durationNS: param.DurationNS,
		stopChan:   make(chan struct{}),
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
	fmt.Println("=== generator starting ===")
	if ok := atomic.CompareAndSwapUint32(&gen.status,
		model.GEN_STA_ORIGIN, model.GEN_STA_STARTING); !ok {
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

	if ok := atomic.CompareAndSwapUint32(&gen.status,
		model.GEN_STA_STARTING, model.GEN_STA_STARTED); !ok {
		return false
	}
	fmt.Println("=== generator started ===")
	return true
}

func (gen *myGenerator) genLoad(throttle <-chan time.Time) {
	for {
		select {
		case <-gen.ctx.Done():
			gen.prepareToStop()
			return
		default:
		}
		go gen.asyncCall()
		select {
		case <-throttle:
			// fmt.Println(" === !!!NEXT!!! ===")
		case <-gen.ctx.Done():
			gen.prepareToStop()
			return
		}
	}
}

// 发起一次调用
func (gen *myGenerator) asyncCall() {
	gen.tickets.Tack()
	defer gen.tickets.Return()
	start := time.Now()
	req := gen.caller.BuildReq()
	resp, err := gen.caller.Call(req)
	elp := time.Since(start)
	if err != nil {
		res := &model.LDResult{
			ID:  req.ID,
			Req: req,
			Rsp: &model.RawRsp{
				ID:      req.ID,
				RspBody: resp,
				Code:    model.GEN_RTNCODE_INTERERR,
				Err:     err,
			},
			Elapse: elp,
		}
		gen.resultChan <- res
		return
	}
	res := gen.caller.CheckRsp(req, resp)
	res.Elapse = elp
	// gen.resultChan <- res
	go gen.sendResult(res)
}

func (gen *myGenerator) sendResult(res *model.LDResult) {
	fmt.Println("1 send result:", res.ID)
	if gen.status == model.GEN_STA_STARTED {
		// select {
		// case <-gen.stopChan:
		// 	// close(gen.resultChan)
		// 	fmt.Println("result channel stopped,return")
		// 	// return
		// case gen.resultChan <- res:
		// 	fmt.Println("2 send result:", res.ID)
		// 	// return
		// }
		gen.resultChan <- res
		fmt.Println("2 send result:", res.ID)
	} else {
		fmt.Println("generator stopped, ignore res:", res)
	}
}

func (gen *myGenerator) Stop() bool {
	fmt.Println("=== generator stopping ===")
	if ok := atomic.CompareAndSwapUint32(&gen.status,
		model.GEN_STA_STARTED, model.GEN_STA_STOPPING); !ok {
		fmt.Println(" Swap status from GEN_STA_STARTED to GEN_STA_STOPPING error")
		return false
	}
	gen.cancelFunc() // 会触发 ctx.Done
	for {
		if gen.status == model.GEN_STA_STOPPED {
			fmt.Println("=== generator stopped,End Stop ===")
			return true
		}
		time.Sleep(time.Second)
	}
}

func (gen *myGenerator) prepareToStop() {
	fmt.Println("=== generator prepareToStop ===")
	atomic.CompareAndSwapUint32(&gen.status, model.GEN_STA_STARTED, model.GEN_STA_STOPPING)
	close(gen.resultChan)
	// close(gen.stopChan) // 通知中间信号通道
	atomic.CompareAndSwapUint32(&gen.status,
		model.GEN_STA_STOPPING, model.GEN_STA_STOPPED)
	fmt.Println("=== generator stopped:", gen.status, " ===")
}

func (gen *myGenerator) Status() uint32 {
	return gen.status
}
