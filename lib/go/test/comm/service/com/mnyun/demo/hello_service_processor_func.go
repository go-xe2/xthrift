package demo

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
)

type HelloServiceSayHello struct {
	*xthrift.TBaseProcessorFunction
	handler HelloService
}

func newHelloServiceSayHello(handler HelloService) *HelloServiceSayHello {
	inst := &HelloServiceSayHello{handler: handler}
	inst.TBaseProcessorFunction = xthrift.NewBaseProcessorFunction(inst)
	return inst
}

func (p *HelloServiceSayHello) Invoke(args thrift.TStruct) (thrift.TStruct, error) {
	if input, ok := args.(*HelloServiceSayHelloArgs); ok {
		result := NewHelloServiceSayHelloResult()

		success := p.handler.SayHello(input.Name)
		result.Success = success

		return result, nil
	}
	return nil, thrift.NewTApplicationException(thrift.INVALID_DATA, "输入参数错误")
}

func (p *HelloServiceSayHello) GetInputArgsInstance() thrift.TStruct {
	return NewHelloServiceSayHelloArgs()

}
