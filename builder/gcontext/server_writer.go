/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-16 14:30
* Description:
*****************************************************************/

package gcontext

import (
	"fmt"
	"github.com/go-xe2/x/type/xstring"
	"github.com/go-xe2/xthrift/builder"
	"github.com/go-xe2/xthrift/pdl"
)

type TServerWriter struct {
	*TWriter
}

var _ builder.ServerCodeWriter = (*TServerWriter)(nil)

func NewServerCodeWriter(cxt *TContext, fileName string) (w *TServerWriter, err error) {
	inst := &TServerWriter{}
	inst.TWriter, err = NewWriter(cxt, fileName)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func (p *TServerWriter) WriteServicesBegin(ns *pdl.FileNamespace) error {
	return nil
}

func (p *TServerWriter) WriteServices(ns *pdl.FileNamespace, svcs map[string]*pdl.FileService) error {
	_, s := pdl.NamespaceLastName(ns.Namespace)
	svrName := xstring.UcFirst(s) + "Server"

	if err := p.writerServerDefine(ns, svrName); err != nil {
		return err
	}
	for _, svc := range svcs {
		if svc.Namespace != ns.Namespace {
			p.Import(svc.Namespace, true)
		}
		if err := p.writeServerRegisterServiceHandler(ns, svc, svrName); err != nil {
			return err
		}
	}
	return nil
}

func (p *TServerWriter) WriteServicesEnd(ns *pdl.FileNamespace) error {
	return nil
}

func (p *TServerWriter) writerServerDefine(ns *pdl.FileNamespace, svrName string) error {
	p.Import("github.com/apache/thrift/lib/go/thrift", false)
	p.Import("github.com/go-xe2/xthrift/lib/go/xthrift", false)

	p.Write("\n")
	p.Write(fmt.Sprintf("type %s struct {\n", svrName))
	p.Write("\tserver *xthrift.TXServer\n")
	p.Write("\tlistenAddr string\n")
	p.Write("\tprocessors []thrift.TProcessor\n")
	p.Write("\tnamespaces *xthrift.TNamespaceProcessor\n")
	p.Write("}\n")
	p.Write("\n")

	// constructor
	p.Write("\n")
	p.Write(fmt.Sprintf("func New%s(addr string) (*%s, error) {\n", svrName, svrName))
	p.Write("\tsvr, err := xthrift.NewServer(addr)\n")
	p.Write("\tif err != nil {\n")
	p.Write("\t\treturn nil, err\n")
	p.Write("\t}\n")

	p.Write(fmt.Sprintf("\tinst := &%s{\n", svrName))
	p.Write("\t\tlistenAddr: addr,\n")
	p.Write("\t\tnamespaces: xthrift.NamespaceProcessor(),\n")
	p.Write("\t\tserver: svr,\n")
	p.Write("\t\tprocessors: make([]thrift.TProcessor, 0),\n")
	p.Write("\t}\n")
	p.Write("\treturn inst, nil\n")
	p.Write("}\n")

	// RegisterProcessor
	p.Write("\n")
	p.Write(fmt.Sprintf("func (p *%s) RegisterProcessor(namespace string, processor thrift.TProcessor) {\n", svrName))
	p.Write("\t_ = p.namespaces.RegisterNamespace(namespace, processor)\n")
	p.Write("\tp.processors = append(p.processors, processor)\n")
	p.Write("}\n")

	// Serve
	p.Write("\n")
	p.Write(fmt.Sprintf("func (p *%s) Serve() error {\n", svrName))
	p.Write("\treturn p.server.Serve()\n")
	p.Write("}\n")
	// Stop
	p.Write("\n")
	p.Write(fmt.Sprintf("func (p *%s) Stop() error {\n", svrName))
	p.Write("\treturn p.server.Stop()\n")
	p.Write("}\n")

	return nil
}

func (p *TServerWriter) writeServerRegisterServiceHandler(ns *pdl.FileNamespace, svc *pdl.FileService, svrName string) error {
	prefix := ""
	svcPrefix := ""
	if svc.Namespace != ns.Namespace {
		_, s := pdl.NamespaceLastName(svc.Namespace)
		prefix = s + "."
		svcPrefix = xstring.UcFirst(s)
	}

	p.Write("\n")
	p.Write(fmt.Sprintf(fmt.Sprintf("func (p *%s) Register%s%s(handler %s) {\n", svrName, svcPrefix, p.GenServiceNameCode(svc), prefix+p.GenServiceNameCode(svc))))
	p.Write(fmt.Sprintf("\tprocessor := %sNew%s(handler)\n", prefix, p.GenServiceNameCode(svc)+"Processor"))
	p.Write(fmt.Sprintf("\tp.RegisterProcessor(\"%s.%s\", processor)\n", svc.Namespace, xstring.LcFirst(svc.Name)))
	p.Write("}\n")
	return nil
}
