package lib

import "tinywell/loadgen/model"

// myGenerator 载荷发生器实现
type myGenerator struct {
	caller     Caller
	tickers    Tickets
	resultChan chan *model.LDResult
}

func (gen *myGenerator) Start() bool {
	return true
}

func (gen *myGenerator) Stop() bool {
	return true
}

func (gen *myGenerator) Status() uint32 {
	return model.GEN_STA_ORIGIN
}
