package demo

import (
	"github.com/go-xe2/xthrift/builder/test/build/go/com/mnyun/types"
)

type HelloService interface {
	// sayHello的说明
	SayHello(name *string, age int32) *types.helloResult
}
