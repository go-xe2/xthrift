package demo

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
	"golang.org/x/net/context"
)

type HelloServiceClient struct {
	*xthrift.TXClient
}

func NewHelloServiceClient(trans thrift.TTransport, in, out thrift.TProtocolFactory) *HelloServiceClient {
	inst := &HelloServiceClient{
		TXClient: xthrift.NewClient(trans, in, out),
	}
	return inst
}

// 测试接口
func (p *HelloServiceClient) SayHello(cxt context.Context, name string) (*HelloResult, error) {
	var args = NewHelloServiceSayHelloArgs()
	args.Name = name
	result := NewHelloServiceSayHelloResult()
	err := p.Call(cxt, "SayHello", args, result)
	if err != nil {
		return nil, err
	}
	return result.GetSuccess(), nil
}
