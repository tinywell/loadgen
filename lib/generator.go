package lib

import (
	"tinywell/loadgen/model"
)

const (
	GEN_STA_ORIGIN   uint32 = 0
	GEN_STA_STARTING uint32 = 1
	GEN_STA_STARTED  uint32 = 2
	GEN_STA_STOPPING uint32 = 3
	GEN_STA_STOPPED  uint32 = 4
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
