package demo

type HelloService interface {
	// 测试接口
	SayHello(name string) *HelloResult
}
