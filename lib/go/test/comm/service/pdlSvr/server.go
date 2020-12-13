package pdlSvr

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/xthrift/lib/go/test/comm/service/com/mnyun/demo"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
)

type PdlSvrServer struct {
	server     *xthrift.TXServer
	listenAddr string
	processors []thrift.TProcessor
	namespaces *xthrift.TNamespaceProcessor
}

func NewPdlSvrServer(addr string) (*PdlSvrServer, error) {
	svr, err := xthrift.NewServer(addr)
	if err != nil {
		return nil, err
	}
	inst := &PdlSvrServer{
		listenAddr: addr,
		namespaces: xthrift.NamespaceProcessor(),
		server:     svr,
		processors: make([]thrift.TProcessor, 0),
	}
	return inst, nil
}

func (p *PdlSvrServer) RegisterProcessor(namespace string, processor thrift.TProcessor) {
	_ = p.namespaces.RegisterNamespace(namespace, processor)
	p.processors = append(p.processors, processor)
}

func (p *PdlSvrServer) Serve() error {
	return p.server.Serve()
}

func (p *PdlSvrServer) Stop() error {
	return p.server.Stop()
}

func (p *PdlSvrServer) RegisterDemoHelloService(handler demo.HelloService) {
	processor := demo.NewHelloServiceProcessor(handler)
	p.RegisterProcessor("com.mnyun.demo.helloService", processor)
}
