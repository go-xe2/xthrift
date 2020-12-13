package gcontext

import (
	"fmt"
	"github.com/go-xe2/xthrift/builder"
	"github.com/go-xe2/xthrift/pdl"
	"sort"
)

type TServiceWriter struct {
	*TWriter
}

var _ builder.ServiceCodeWriter = (*TServiceWriter)(nil)

func NewServiceWriter(cxt *TContext, fileName string) (w *TServiceWriter, err error) {
	inst := &TServiceWriter{}
	inst.TWriter, err = NewWriter(cxt, fileName)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func (p *TServiceWriter) WriteServiceBegin(ns *pdl.FileNamespace, svc *pdl.FileService) (ident string, err error) {
	ident = p.GenServiceNameCode(svc)
	p.Write("\n")
	p.Write(fmt.Sprintf("type %s interface {\n", ident))
	return ident, nil
}

func (p *TServiceWriter) WriteServiceMethod(ns *pdl.FileNamespace, svc *pdl.FileService, method *pdl.FileServiceMethod) (ident string, err error) {
	ident = p.GenServiceMethodNameCode(method)
	args := method.Args
	sort.Slice(args, func(i, j int) bool {
		return args[i].Id-args[j].Id < 0
	})
	if method.Summary != "" {
		p.Write(fmt.Sprintf("\t// %s\n", method.Summary))
	}
	p.Write(fmt.Sprintf("\t%s(", ident))
	size := len(args)
	namespaces := make([]string, 0)
	if method.Result != nil {
		tmp := p.getDataTypeNamespace(ns, method.Result)
		for _, s := range tmp {
			namespaces = append(namespaces, s)
		}
	}
	if method.Exception != nil && method.Exception.Type != pdl.SPD_VOID {
		tmp := p.getDataTypeNamespace(ns, method.Exception)
		for _, s := range tmp {
			namespaces = append(namespaces, s)
		}
	}
	for _, s := range namespaces {
		p.Import(s, true)
	}

	for i := 0; i < size; i++ {
		arg := args[i]

		namespaces := p.getDataTypeNamespace(ns, arg.FieldType)
		for _, s := range namespaces {
			p.Import(s, true)
		}

		if i > 0 {
			p.Write(" ")
		}
		p.Write(fmt.Sprintf("%s %s", arg.Name, p.GenFieldDefineTypeCode(ns, arg)))
		if i < size-1 {
			p.Write(",")
		}
	}
	p.Write(")")
	if method.Exception != nil && method.Exception.Type != pdl.SPD_VOID {
		p.Write(" (")
		p.Write(fmt.Sprintf("result %s, ", p.GenDataTypeCode(ns, method.Result)))
		p.Write(fmt.Sprintf("err %s", p.GenDataTypeCode(ns, method.Exception)))
		p.Write(")")
	} else {
		p.Write(fmt.Sprintf(" %s", p.GenDataTypeCode(ns, method.Result)))
	}
	p.Write("\n")
	return ident, nil
}

func (p *TServiceWriter) WriteServiceEnd(ns *pdl.FileNamespace, svc *pdl.FileService) error {
	p.Write("}\n")
	return nil
}
