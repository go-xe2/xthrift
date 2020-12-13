package gcontext

import (
	"fmt"
	"github.com/go-xe2/x/type/xstring"
	"github.com/go-xe2/xthrift/builder"
	"github.com/go-xe2/xthrift/builder/comm"
	"github.com/go-xe2/xthrift/pdl"
	"sort"
)

type TProcessorFunCodeWriter struct {
	*TWriter
}

var _ builder.ProcessorFunCodeWriter = (*TProcessorFunCodeWriter)(nil)

func NewProcessorFunCodeWriter(cxt *TContext, fileName string) (w *TProcessorFunCodeWriter, err error) {
	inst := &TProcessorFunCodeWriter{}
	inst.TWriter, err = NewWriter(cxt, fileName)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

//
//func (p *TProcessorFunCodeWriter) WriteInclude(namespaces []string, files map[string]string) error {
//	if err := p.TWriter.WriteInclude(namespaces, files); err != nil {
//		return err
//	}
//	p.Write("import \"github.com/apache/thrift/lib/go/thrift\"\n")
//	p.Write("import \"github.com/go-xe2/xthrift/lib/go/xthrift\"\n")
//	return nil
//}

func (p *TProcessorFunCodeWriter) WriteFunction(ns *pdl.FileNamespace, svc *pdl.FileService, method *pdl.FileServiceMethod, input, result *pdl.FileStruct) (ident string, err error) {
	ident = p.GenServiceNameCode(svc) + p.GenServiceMethodNameCode(method)
	p.Write("\n")

	if err := p.writeFunStructDefine(ident, svc); err != nil {
		return ident, err
	}
	// 写构造方法
	p.Write("\n")
	if err := p.writeFunStructConstructMethod(ident, svc); err != nil {
		return ident, err
	}
	// 创建调用方法
	p.Write("\n")
	if err := p.writeFunInvokeMethod(ns, ident, method, input, result); err != nil {
		return ident, err
	}
	// 创建新建输入参数实现方法
	if err := p.writeStructNewInputArgInstanceMethod(ns, ident, method, input); err != nil {
		return ident, err
	}
	return ident, nil
}

func (p *TProcessorFunCodeWriter) writeFunStructDefine(fnStruName string, svc *pdl.FileService) error {
	p.Import("github.com/go-xe2/xthrift/lib/go/xthrift", false)

	p.Write(fmt.Sprintf("type %s struct {\n", fnStruName))
	p.Write("\t*xthrift.TBaseProcessorFunction\n")
	p.Write(fmt.Sprintf("\thandler %s\n", p.GenServiceNameCode(svc)))
	p.Write("}\n")
	return nil
}

func (p *TProcessorFunCodeWriter) writeFunStructConstructMethod(fnStruName string, svc *pdl.FileService) error {
	p.Import("github.com/go-xe2/xthrift/lib/go/xthrift", false)

	p.Write("\n")
	p.Write(fmt.Sprintf("func new%s(handler %s) *%s {\n", fnStruName, p.GenServiceNameCode(svc), fnStruName))
	p.Write(fmt.Sprintf("\tinst := &%s{handler: handler}\n", fnStruName))
	p.Write("\tinst.TBaseProcessorFunction = xthrift.NewBaseProcessorFunction(inst)\n")
	p.Write("\treturn inst\n")
	p.Write("}\n")
	return nil
}

func (p *TProcessorFunCodeWriter) writeFunInvokeMethod(ns *pdl.FileNamespace, fnStruName string, method *pdl.FileServiceMethod, input *pdl.FileStruct, result *pdl.FileStruct) error {
	p.Import("github.com/apache/thrift/lib/go/thrift", false)

	p.Write("\n")
	p.Write(fmt.Sprintf("func (p *%s) Invoke(args thrift.TStruct) (thrift.TStruct, error) {\n", fnStruName))
	if len(method.Args) > 0 {
		p.Write(fmt.Sprintf("\tif input, ok := args.(%s); ok {\n", p.GenDataTypeCode(ns, input.Type)))
	} else {
		p.Write(fmt.Sprintf("\tif _, ok := args.(%s); ok {\n", p.GenDataTypeCode(ns, input.Type)))
	}
	// create result instance
	p.Write("\t")
	p.Write(comm.AppendIndent(1, p.GenCreateStructTypeCode(ns, "result", ":=", result.Type)))
	p.Write("\n")

	// 写接收参数
	resultSize := len(result.Fields)
	resultFields := make([]*pdl.FileDataField, resultSize)
	i := 0
	for _, f := range result.Fields {
		resultFields[i] = f
		i++
	}

	sort.Slice(resultFields, func(i, j int) bool {
		return resultFields[i].Id-resultFields[j].Id < 0
	})

	p.Write("\t\t")
	for i := 0; i < resultSize; i++ {
		fd := resultFields[i]
		p.Write(xstring.LcFirst(fd.Name))
		if i < resultSize-1 {
			p.Write(",")
		}
	}
	p.Write(fmt.Sprintf(" := p.handler.%s(", p.GenServiceMethodNameCode(method)))

	// 写输入参数
	size := len(input.Fields)
	args := make([]*pdl.FileDataField, size)
	i = 0
	for _, f := range input.Fields {
		args[i] = f
		i++
	}
	sort.Slice(args, func(i, j int) bool {
		return args[i].Id-args[j].Id < 0
	})
	for i := 0; i < size; i++ {
		arg := args[i]
		p.Write(fmt.Sprintf("input.%s", p.GenFieldNameCode(arg)))
		if i < size-1 {
			p.Write(",")
		}
	}
	// call handler method end
	p.Write(")\n")

	// 获取返回参数
	for i := 0; i < resultSize; i++ {
		fd := resultFields[i]
		s := p.GenFieldAssignValueCode(ns, "result", fd, xstring.LcFirst(fd.Name))
		p.Write(comm.AppendIndent(2, s))
		p.Write("\n")
	}

	p.Write("\t\treturn result, nil\n")
	p.Write("\t}\n")
	p.Write("\treturn nil, thrift.NewTApplicationException(thrift.INVALID_DATA, \"输入参数错误\")")
	p.Write("}\n")
	return nil
}

func (p *TProcessorFunCodeWriter) writeStructNewInputArgInstanceMethod(ns *pdl.FileNamespace, fnStruName string, method *pdl.FileServiceMethod, input *pdl.FileStruct) error {
	p.Import("github.com/apache/thrift/lib/go/thrift", false)

	p.Write("\n")
	p.Write(fmt.Sprintf("func (p *%s) GetInputArgsInstance() thrift.TStruct {\n", fnStruName))
	p.Write(comm.AppendIndent(1, p.GenCreateStructTypeCode(ns, "return ", "", input.Type)))
	p.Write("\n")
	p.Write("}\n")
	return nil
}
