package gcontext

import (
	"fmt"
	"github.com/go-xe2/xthrift/builder"
	"github.com/go-xe2/xthrift/pdl"
	"sort"
	"strings"
)

type TProcessorCodeWriter struct {
	*TWriter
}

var _ builder.ProcessorCodeWriter = (*TProcessorCodeWriter)(nil)

func NewProcessorCodeWriter(cxt *TContext, fileName string) (w *TProcessorCodeWriter, err error) {
	inst := &TProcessorCodeWriter{}
	inst.TWriter, err = NewWriter(cxt, fileName)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func (p *TProcessorCodeWriter) WriteProcessor(ns *pdl.FileNamespace, svc *pdl.FileService) (ident string, err error) {
	ident = p.GenServiceNameCode(svc) + "Processor"
	p.Write("\n")
	if err := p.writeProcessorDefine(ns, ident, svc); err != nil {
		return ident, err
	}

	if err := p.writerProcessorConstructorMethod(ns, ident, svc); err != nil {
		return ident, err
	}
	if err := p.writeProcessorRegisterMethods(ns, ident, svc); err != nil {
		return ident, err
	}
	return ident, nil
}

func (p *TProcessorCodeWriter) writeProcessorDefine(ns *pdl.FileNamespace, processorName string, svc *pdl.FileService) error {
	p.Import("github.com/go-xe2/xthrift/lib/go/xthrift", false)

	p.Write(fmt.Sprintf("type %s struct {\n", processorName))
	p.Write("\t*xthrift.TBaseProcessor\n")
	p.Write(fmt.Sprintf("\thandler %s\n", p.GenServiceNameCode(svc)))
	p.Write("}\n")
	return nil
}

func (p *TProcessorCodeWriter) writerProcessorConstructorMethod(ns *pdl.FileNamespace, processorName string, svc *pdl.FileService) error {
	p.Import("github.com/go-xe2/xthrift/lib/go/xthrift", false)

	p.Write("\n")
	p.Write(fmt.Sprintf("func New%s(handler %s) *%s {\n", processorName, p.GenServiceNameCode(svc), processorName))
	p.Write(fmt.Sprintf("\tinst := &%s{\n", processorName))
	p.Write("\t\thandler: handler,\n")
	p.Write("\t}\n")
	p.Write("\n")
	p.Write("\tinst.TBaseProcessor = xthrift.NewBaseProcessor(inst)\n")
	p.Write("\treturn inst.registerFunctions()\n")
	p.Write("}\n")
	return nil
}

func (p *TProcessorCodeWriter) writeProcessorRegisterMethods(ns *pdl.FileNamespace, processorName string, svc *pdl.FileService) error {
	p.Write("\n")
	p.Write(fmt.Sprintf("func (p *%s) registerFunctions() *%s {\n", processorName, processorName))
	size := len(svc.Methods)

	methods := make([]*pdl.FileServiceMethod, size)
	i := 0
	for _, m := range svc.Methods {
		methods[i] = m
		i++
	}

	sort.Slice(methods, func(i, j int) bool {
		return strings.Compare(methods[i].Name, methods[j].Name) < 0
	})

	for i := 0; i < size; i++ {
		method := methods[i]
		newfunName := "new" + p.GenServiceNameCode(svc) + p.GenServiceMethodNameCode(method) + "(p.handler)"
		p.Write(fmt.Sprintf("\tp.RegisterFunction(\"%s\", %s)\n", p.GenServiceMethodNameCode(method), newfunName))
	}
	p.Write("\treturn p\n")
	p.Write("}\n")
	return nil
}
