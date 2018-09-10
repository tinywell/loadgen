package model

import (
	"time"
)

const (
	GEN_STA_ORIGIN uint32 = 0
)

type LDResult struct {
	ID     int64
	Req    *RawReq
	Rsp    *RawRsp
	Msg    string
	Elapse time.Duration
}

type RawReq struct {
	ID      int64
	ReqBody []byte
}

type RawRsp struct {
	ID      int64
	RspBody []byte
	Code    uint32
	Err     error
}
