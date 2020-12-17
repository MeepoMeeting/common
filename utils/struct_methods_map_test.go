package utils

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testWebImp struct {
}

type testReq struct {
	Name string
}
type testRsp struct {
	Rsp string
}

func (webImp *testWebImp) Hello(ctx context.Context, req *testReq) (*testRsp, error) {
	rsp := &testRsp{Rsp: fmt.Sprintf("hello %s", req.Name)}
	return rsp, nil
}
func (webImp *testWebImp) HelloError(ctx context.Context, req *testReq) (*testRsp, error) {
	rsp := &testRsp{Rsp: fmt.Sprintf("hello %s", req.Name)}
	return rsp, fmt.Errorf("error")
}

func TestMethodsMap_Call(t *testing.T) {
	mm := &MethodsMap{}
	imp := &testWebImp{}
	mm.Init(imp)
	req := "{\"Name\":\"jhx\"}"
	rsp, err := mm.Call("Hello", req)
	assert.Equal(t, "{\"Rsp\":\"hello jhx\"}", rsp)
	assert.Equal(t, nil, err)
	rsp, err = mm.Call("HelloError", req)
	assert.NotEqual(t, nil, err)
	rsp, err = mm.Call("NotFound", req)
	assert.Equal(t, fmt.Errorf("%s,MethodNotFound", "NotFound"), err)
}
