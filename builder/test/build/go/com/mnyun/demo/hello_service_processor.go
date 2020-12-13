package demo

import (
	"github.com/go-xe2/xthrift/lib/go/xthrift"
)

type HelloServiceProcessor struct {
	*xthrift.TBaseProcessor
	handler HelloService
}

func NewHelloServiceProcessor(handler HelloService) *HelloServiceProcessor {
	inst := &HelloServiceProcessor{
		handler: handler,
	}

	inst.TBaseProcessor = xthrift.NewBaseProcessor(inst)
	return inst.registerFunctions()
}

func (p *HelloServiceProcessor) registerFunctions() *HelloServiceProcessor {
	p.RegisterFunction("SayHello", newHelloServiceSayHello(p.handler))
	return p
}
