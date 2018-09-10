package lib

import (
	"tinywell/loadgen/model"
)

// Generator 载荷发生器接口
type Generator interface {
	Start() bool
	Stop() bool
	Status() uint32
}

// Caller 交易掉用模块接口
type Caller interface {
	BuildReq() *model.RawReq
	Call(req []byte) (rsp []byte, err error)
	CheckRsp(rsp []byte) *model.LDResult
}

// Tickets 并发票池接口
type Tickets interface {
	Tack()
	Return()
	Active() bool
	Total() int32
	Remain() int32
}
