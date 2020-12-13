/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-16 14:29
* Description:
*****************************************************************/

package gcontext

import (
	"fmt"
	"github.com/go-xe2/x/type/xstring"
	"github.com/go-xe2/xthrift/builder"
	"github.com/go-xe2/xthrift/pdl"
	"sort"
	"strings"
)

type TClientWriter struct {
	*TWriter
}

var _ builder.ClientCodeWrite = (*TClientWriter)(nil)

func NewClientCodeWriter(cxt *TContext, fileName string) (w *TClientWriter, err error) {
	inst := &TClientWriter{}
	inst.TWriter, err = NewWriter(cxt, fileName)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func (p *TClientWriter) WriteClientBegin(ns *pdl.FileNamespace, svc *pdl.FileService) error {
	return nil
}

func (p *TClientWriter) WriteClient(ns *pdl.FileNamespace, svc *pdl.FileService) error {
	p.Write("\n")
	cliName := p.GenServiceNameCode(svc) + "Client"
	if err := p.writeClientDefine(ns, cliName, svc); err != nil {
		return err
	}

	if err := p.writeClientConstructor(ns, cliName, svc); err != nil {
		return err
	}

	if err := p.writeClientMethods(ns, cliName, svc); err != nil {
		return err
	}

	return nil
}

func (p *TClientWriter) WriteClientEnd(ns *pdl.FileNamespace, svc *pdl.FileService) error {
	return nil
}

func (p *TClientWriter) writeClientDefine(ns *pdl.FileNamespace, cliName string, svc *pdl.FileService) error {
	p.Import("github.com/go-xe2/xthrift/lib/go/xthrift", false)

	p.Write("\n")
	p.Write(fmt.Sprintf("type %s struct {\n", cliName))
	p.Write("\t*xthrift.TXClient\n")
	p.Write("}\n")
	return nil
}

func (p *TClientWriter) writeClientConstructor(ns *pdl.FileNamespace, cliName string, svc *pdl.FileService) error {
	p.Import("github.com/apache/thrift/lib/go/thrift", false)

	p.Write("\n")
	p.Write(fmt.Sprintf("func New%s(trans thrift.TTransport, in, out thrift.TProtocolFactory) *%s{\n", cliName, cliName))
	p.Write(fmt.Sprintf("\tinst := &%s{\n", cliName))
	p.Write("\t\tTXClient: xthrift.NewClient(trans, in, out),\n")
	p.Write("\t}\n")
	p.Write("\treturn inst\n")
	p.Write("}\n")
	return nil
}

func (p *TClientWriter) writeClientMethods(ns *pdl.FileNamespace, cliName string, svc *pdl.FileService) error {
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
		if err := p.writeClientMethod(ns, cliName, svc, methods[i]); err != nil {
			return err
		}
	}
	return nil
}

func (p *TClientWriter) writeClientMethod(ns *pdl.FileNamespace, cliName string, svc *pdl.FileService, method *pdl.FileServiceMethod) error {
	p.Import("golang.org/x/net/context", false)
	useNamespaces := make([]string, 0)
	if method.Result != nil {
		tmp := p.getDataTypeNamespace(ns, method.Result)
		for _, s := range tmp {
			useNamespaces = append(useNamespaces, s)
		}
	}
	if method.Exception != nil && method.Exception.Type != pdl.SPD_VOID {
		tmp := p.getDataTypeNamespace(ns, method.Exception)
		for _, s := range tmp {
			useNamespaces = append(useNamespaces, s)
		}
	}

	for _, s := range useNamespaces {
		p.Import(s, true)
	}

	p.Write("\n")
	if method.Summary != "" {
		p.Write(fmt.Sprintf("// %s\n", method.Summary))
	}
	p.Write(fmt.Sprintf("func (p *%s) %s(cxt context.Context,", cliName, p.GenServiceMethodNameCode(method)))
	args := method.Args
	sort.Slice(args, func(i, j int) bool {
		return args[i].Id-args[j].Id < 0
	})
	// 生成函数输入参数
	size := len(args)
	for i := 0; i < size; i++ {
		arg := args[i]

		namespaces := p.getDataTypeNamespace(ns, arg.FieldType)
		for _, s := range namespaces {
			p.Import(s, true)
		}

		p.Write(fmt.Sprintf("%s %s", xstring.LcFirst(arg.Name), p.GenFieldDefineTypeCode(ns, arg)))
		if i < size-1 {
			p.Write(",")
		}
	}
	p.Write(fmt.Sprintf(") (%s,", p.GenDataTypeCode(ns, method.Result)))
	if method.Exception != nil && method.Exception.Type != pdl.SPD_VOID {
		p.Write(fmt.Sprintf(" %s,", p.GenDataTypeCode(ns, method.Exception)))
	}
	p.Write(" error) {\n")
	// 初始化输入参数
	p.Write(fmt.Sprintf("\tvar args = New%s()\n", p.GenServiceNameCode(svc)+p.GenServiceMethodNameCode(method)+"Args"))
	// 设置参数值
	for i := 0; i < size; i++ {
		arg := args[i]
		p.Write(fmt.Sprintf("\targs.%s = %s\n", p.GenFieldNameCode(arg), xstring.LcFirst(arg.Name)))
	}
	p.Write(fmt.Sprintf("\tresult := New%s()\n", p.GenServiceNameCode(svc)+p.GenServiceMethodNameCode(method)+"Result"))
	p.Write(fmt.Sprintf("\terr := p.Call(cxt, \"%s\", args, result)\n", p.GenServiceMethodNameCode(method)))
	p.Write("\tif err != nil {\n")
	switch method.Result.Type {
	case pdl.SPD_STR:
		p.Write("\t\treturn \"\", err\n")
		break
	case pdl.SPD_BOOL:
		p.Write("\t\treturn false, err\n")
		break
	case pdl.SPD_I08:
		p.Write("\t\treturn 0, err\n")
		break
	case pdl.SPD_I16:
		p.Write("\t\treturn 0, err\n")
		break
	case pdl.SPD_I32:
		p.Write("\t\treturn 0, err\n")
		break
	case pdl.SPD_I64:
		p.Write("\t\treturn 0, err\n")
		break
	case pdl.SPD_DOUBLE:
		p.Write("\t\treturn 0, err\n")
		break
	case pdl.SPD_LIST, pdl.SPD_SET, pdl.SPD_STRUCT, pdl.SPD_EXCEPTION:
		p.Write("\t\treturn nil, err\n")
		break
	}
	p.Write("\t}\n")

	p.Write("\treturn result.GetSuccess(),")
	if method.Exception != nil && method.Exception.Type != pdl.SPD_VOID {
		p.Write("result.GetException(),")
	}
	p.Write("nil\n")

	// end func
	p.Write("}\n") // end func
	return nil
}
