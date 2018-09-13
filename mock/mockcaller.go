package mock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"tinywell/loadgen/model"
)

type MockCaller struct {
}

func NewMockCaller() model.Caller {
	return &MockCaller{}
}

func (caller *MockCaller) BuildReq() *model.RawReq {
	return &model.RawReq{
		ID:      time.Now().UnixNano(),
		ReqBody: []byte("Hello"),
	}
}

func (caller *MockCaller) Call(req *model.RawReq) ([]byte, error) {
	fmt.Println("== call id:", req.ID, " ===")
	rsp, err := http.Get("http://www.baidu.com")
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(rsp.Body)
}

func (caller *MockCaller) CheckRsp(req *model.RawReq, rsp []byte) *model.LDResult {
	rawRsp := &model.RawRsp{
		ID:      req.ID,
		RspBody: rsp,
		Code:    model.GEN_RTNCODE_SUCCESS,
		Err:     nil,
	}
	return &model.LDResult{
		ID:  req.ID,
		Req: req,
		Rsp: rawRsp,
		Msg: "hello",
	}
}
