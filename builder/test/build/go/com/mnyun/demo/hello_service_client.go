package demo

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/xthrift/builder/test/build/go/com/mnyun/types"
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

// sayHello的说明
func (p *HelloServiceClient) SayHello(cxt context.Context, name *string, age int32) (*types.HelloResult, error) {
	var args = NewHelloServiceSayHelloArgs()
	args.Name = name
	args.Age = age
	result := NewHelloServiceSayHelloResult()
	err := p.Call(cxt, "SayHello", args, result)
	if err != nil {
		return nil, err
	}
	return result.GetSuccess(), nil
}
