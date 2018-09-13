package mock

import "testing"

func TestBuildReq(t *testing.T) {
	caller := NewMockCaller()
	req := caller.BuildReq()
	t.Log(req)
}

func TestCall(t *testing.T) {
	caller := NewMockCaller()
	req := caller.BuildReq()
	rsp, err := caller.Call(req.ReqBody)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(rsp))
}

func TestCheck(t *testing.T) {
	caller := NewMockCaller()
	req := caller.BuildReq()
	rsp, err := caller.Call(req.ReqBody)
	if err != nil {
		t.Error(err)
	}
	res := caller.CheckRsp(req, rsp)
	t.Log(res)
}
