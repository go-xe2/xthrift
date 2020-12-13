/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-13 10:47
* Description:
*****************************************************************/

package builder

import (
	"fmt"
	"github.com/go-xe2/x/type/xstring"
	"github.com/go-xe2/xthrift/pdl"
	"sort"
	"strings"
)

type TProtoBuilder struct {
	cxt     Context
	proj    *pdl.FileProject
	buildNS map[string]*IdentMap
}

func NewProtoBuilder(cxt Context, proj *pdl.FileProject) *TProtoBuilder {
	return &TProtoBuilder{
		cxt:     cxt,
		proj:    proj,
		buildNS: make(map[string]*IdentMap),
	}
}

// 生成接口方法输入参数数据结构体
func makeMethodInputStruct(ns *pdl.FileNamespace, svc *pdl.FileService, method *pdl.FileServiceMethod) *pdl.FileStruct {
	name := fmt.Sprintf("%s%sArgs", xstring.UcFirst(svc.Name), xstring.UcFirst(method.Name))
	result := pdl.NewFileStructType(ns.Namespace, name, "工具生成:接口输入参数定义")
	for _, arg := range method.Args {
		fd := pdl.NewFileDataField(arg.Id, arg.Name, arg.FieldType, arg.Summary)
		fd.SetLimit(arg.Limit)
		fd.SetRule(arg.Rule)
		result.AddField(fd)
	}
	return result
}

// 生成接口方法输出参数数据结构体
func makeMethodResultStruct(ns *pdl.FileNamespace, svc *pdl.FileService, method *pdl.FileServiceMethod) *pdl.FileStruct {
	name := fmt.Sprintf("%s%sResult", xstring.UcFirst(svc.Name), xstring.UcFirst(method.Name))

	result := pdl.NewFileStructType(ns.Namespace, name, "工具生成:接口方法返回数据")
	success := pdl.NewFileDataField(1, "Success", method.Result, "接口成功返回数据")
	success.SetLimit(pdl.SPDLimitOptional)
	result.AddField(success)
	if method.Exception != nil && method.Exception.Type != pdl.SPD_VOID {
		fail := pdl.NewFileDataField(2, "Exception", method.Exception, "接口异常返回数据类型")
		fail.SetLimit(pdl.SPDLimitOptional)
		result.AddField(fail)
	}
	return result
}

type tBuilderMethodInfo struct {
	method *pdl.FileServiceMethod
	arg    *pdl.FileStruct
	result *pdl.FileStruct
}

type IdentDefInfo struct {
	// 所在的文件
	idType   IdentType
	fileName string
	// 所在的命名空间
	namespace string
}

func NewIdentDefInfo(it IdentType, namespace string, fileName string) *IdentDefInfo {
	return &IdentDefInfo{
		idType:    it,
		fileName:  fileName,
		namespace: namespace,
	}
}

type IdentMap struct {
	items map[string]*IdentDefInfo
}

func NewIdentMap() *IdentMap {
	return &IdentMap{
		items: make(map[string]*IdentDefInfo),
	}
}

func (p *IdentMap) Add(it IdentType, ident string, namespace, fileName string) {
	key := fmt.Sprintf("%s-%s", namespace, ident)
	if _, ok := p.items[key]; !ok {
		p.items[key] = NewIdentDefInfo(it, namespace, fileName)
	}
}

func (p *IdentMap) Exists(namespace string, ident string) *IdentDefInfo {
	if v, ok := p.items[fmt.Sprintf("%s-%s", namespace, ident)]; ok {
		return v
	}
	return nil
}

func (p *IdentMap) Assign(other *IdentMap) {
	if other == nil {
		return
	}
	for k, v := range other.items {
		if _, ok := p.items[k]; !ok {
			p.items[k] = v
		}
	}
}

func (p *IdentMap) Map() map[string]*IdentDefInfo {
	return p.items
}

func newBuildMethodInfo(method *pdl.FileServiceMethod, arg, result *pdl.FileStruct) *tBuilderMethodInfo {
	return &tBuilderMethodInfo{
		method: method,
		arg:    arg,
		result: result,
	}
}

func (p *TProtoBuilder) buildNamespace(proj *pdl.FileProject, ns *pdl.FileNamespace, idents *IdentMap) error {
	// 收集当前空间中定义的typedefs
	typedefs := make(map[string]*pdl.FileTypeDef)
	for _, f := range ns.Files {
		for defName, def := range f.Typedefs {
			if _, ok := typedefs[def.Name]; !ok {
				typedefs[defName] = def
			}
		}
	}
	// 生成typedef文件
	if len(typedefs) > 0 {
		if err := p.buildTypedefs(ns, typedefs, idents); err != nil {
			return err
		}
	}

	// 生成文件中的定义
	for _, f := range ns.Files {
		if err := p.buildFile(proj, ns, f, idents); err != nil {
			return err
		}
	}

	return nil
}

func (p *TProtoBuilder) buildFile(proj *pdl.FileProject, ns *pdl.FileNamespace, file *pdl.FileData, idents *IdentMap) error {
	result := NewIdentMap()
	// 先生成被引用的空间文件
	for _, k := range file.Imports {
		if k == ns.Namespace {
			// 不编译自己所在的命名空间
			continue
		}
		if _, ok := p.buildNS[k]; ok {
			// 已经编译的命名空间，不再编译
			continue
		}
		curNs := proj.QryNamespace(k)
		if curNs == nil {
			return fmt.Errorf("空间%s未定义", k)
		}
		err := p.buildNamespace(proj, curNs, idents)
		if err != nil {
			return err
		}
		result.Assign(idents)
	}

	// 生成类型定义文件
	for _, t := range file.Types {
		if err := p.buildStruct(ns, t, idents); err != nil {
			return err
		}
	}

	// 生成服务接口输入、输出参数数据结构
	serviceMethods := make(map[*pdl.FileService][]*tBuilderMethodInfo)
	for _, svc := range file.Services {
		methods := make([]*tBuilderMethodInfo, 0)
		for _, method := range svc.Methods {
			argType := makeMethodInputStruct(ns, svc, method)
			if e := p.buildStruct(ns, argType, idents); e != nil {
				return e
			}

			resultType := makeMethodResultStruct(ns, svc, method)
			if e := p.buildStruct(ns, resultType, idents); e != nil {
				return e
			}
			methods = append(methods, newBuildMethodInfo(method, argType, resultType))
		}
		serviceMethods[svc] = methods
	}

	// 生成接口定义
	for _, svc := range file.Services {
		if e := p.buildServiceDefine(ns, svc, idents); e != nil {
			return e
		}
	}

	// 生成处理器方法
	for svc, methods := range serviceMethods {
		if e := p.buildProcessorFuncs(ns, svc, methods, idents); e != nil {
			return e
		}
	}

	// 生成处理器
	for _, svc := range file.Services {
		if e := p.buildProcessor(ns, svc, idents); e != nil {
			return e
		}
	}

	// 生成客户端代码
	for _, svc := range file.Services {
		if err := p.buildClient(ns, svc, idents); err != nil {
			return err
		}
	}
	return nil
}

func (p *TProtoBuilder) Build() (idents *IdentMap, err error) {
	idents = NewIdentMap()
	for _, ns := range p.proj.Namespaces {
		if err := p.buildNamespace(p.proj, ns, idents); err != nil {
			return idents, err
		}
	}
	svcs := p.proj.AllServices()
	// 创建服务所在的空间
	nsItems := p.proj.AllNamespaces()
	sort.Slice(nsItems, func(i, j int) bool {
		return strings.Compare(nsItems[i], nsItems[j]) < 0
	})
	ns := pdl.NewFileNamespace(p.proj, "pdlSvr")
	if err := p.buildServer(ns, svcs, idents); err != nil {
		return idents, err
	}
	return idents, nil
}

func (p *TProtoBuilder) buildTypedefs(currentNs *pdl.FileNamespace, defs map[string]*pdl.FileTypeDef, idents *IdentMap) error {
	defWriter, err := p.cxt.GetTypedefWriter(currentNs)
	if err != nil {
		return err
	}
	defer func() {
		if e := defWriter.Flush(); e != nil {
			fmt.Println("buildTypedefs writer.Flush() error:", err)
		}
		if e := defWriter.Close(); e != nil {
			fmt.Println("BuildTypedefs writer.Close() error:", err)
		}
	}()
	// 写入命名空间
	if err = defWriter.WriteNamespace(currentNs.Namespace); err != nil {
		return err
	}
	// 生成类型别名定义文件
	for _, def := range defs {
		if ident, err := defWriter.WriteDef(currentNs, def); err != nil {
			return err
		} else {
			idents.Add(IDT_TYPE_DEF_NAME, ident, currentNs.Namespace, defWriter.FileName())
		}
	}
	return nil
}

func (p *TProtoBuilder) buildStruct(currentNs *pdl.FileNamespace, stru *pdl.FileStruct, idents *IdentMap) (err error) {
	writer, err := p.cxt.GetStructWriter(currentNs.Namespace, stru)
	if err != nil {
		return err
	}
	defer func() {
		if err := writer.Flush(); err != nil {
			fmt.Println("buildStruct writer.Flush() error:", err)
		}
		if err := writer.Close(); err != nil {
			fmt.Println("buildStruct writer.Flush() error:", err)
		}
	}()
	if err = writer.WriteNamespace(currentNs.Namespace); err != nil {
		return
	}
	// 写入数据结构
	var ident = ""
	if ident, err = writer.WriteStructBegin(currentNs, stru); err != nil {
		return err
	}
	if err = writer.WriteStruct(currentNs, stru); err != nil {
		return
	}
	if err = writer.WriteStructEnd(currentNs, stru); err != nil {
		return
	}
	idents.Add(IDT_TYPE_NAME, ident, currentNs.Namespace, writer.FileName())
	return nil
}

func (p *TProtoBuilder) buildServiceDefine(currentNs *pdl.FileNamespace, svc *pdl.FileService, idents *IdentMap) (err error) {
	writer, err := p.cxt.GetServiceWriter(currentNs.Namespace, svc)
	if err != nil {
		return err
	}
	defer func() {
		if err := writer.Flush(); err != nil {
			fmt.Println("buildServiceDefine writer.Flush() error:", err)
		}
		if err := writer.Close(); err != nil {
			fmt.Println("buildServiceDefine writer.Close() error:", err)
		}
	}()
	if err = writer.WriteNamespace(currentNs.Namespace); err != nil {
		return
	}
	// 写入数据结构
	if ident, err := writer.WriteServiceBegin(currentNs, svc); err != nil {
		return err
	} else {
		idents.Add(IDT_SERVICE_NAME, ident, currentNs.Namespace, writer.FileName())
	}
	size := len(svc.Methods)
	methods := make([]*pdl.FileServiceMethod, size)
	i := 0
	for _, m := range svc.Methods {
		methods[i] = m
		i++
	}
	// 对服务接口排序输出
	sort.Slice(methods, func(i, j int) bool {
		return strings.Compare(methods[i].Name, methods[i].Name) < 0
	})

	for i := 0; i < size; i++ {
		if ident, err := writer.WriteServiceMethod(currentNs, svc, methods[i]); err != nil {
			return err
		} else {
			idents.Add(IDT_METHOD_NAME, ident, currentNs.Namespace, writer.FileName())
		}
	}
	if err = writer.WriteServiceEnd(currentNs, svc); err != nil {
		return
	}
	return nil
}

func (p *TProtoBuilder) buildProcessorFuncs(currentNs *pdl.FileNamespace, svc *pdl.FileService, methods []*tBuilderMethodInfo, idents *IdentMap) (err error) {
	writer, err := p.cxt.GetProcessorFunWriter(currentNs.Namespace, svc)
	if err != nil {
		return err
	}
	defer func() {
		if err := writer.Flush(); err != nil {
			fmt.Println("buildProcessorFuncs writer.Flush() error:", err)
		}
		if err := writer.Close(); err != nil {
			fmt.Println("buildProcessorFuncs writer.Close() error:", err)
		}
	}()

	if err = writer.WriteNamespace(currentNs.Namespace); err != nil {
		return
	}
	// 写入处理器方法
	size := len(methods)
	for i := 0; i < size; i++ {
		if ident, err := writer.WriteFunction(currentNs, svc, methods[i].method, methods[i].arg, methods[i].result); err != nil {
			return err
		} else {
			idents.Add(IDT_PROCESSOR_FUN_NAME, ident, currentNs.Namespace, writer.FileName())
		}
	}
	return err
}

func (p *TProtoBuilder) buildProcessor(currentNs *pdl.FileNamespace, svc *pdl.FileService, idents *IdentMap) (err error) {
	writer, err := p.cxt.GetProcessorWriter(currentNs.Namespace, svc)
	if err != nil {
		return err
	}
	defer func() {
		if err := writer.Flush(); err != nil {
			fmt.Println("buildProcessor writer.Flush() error:", err)
		}
		if err := writer.Close(); err != nil {
			fmt.Println("buildProcessor writer.Close() error:", err)
		}
	}()

	if err = writer.WriteNamespace(currentNs.Namespace); err != nil {
		return
	}
	if ident, err := writer.WriteProcessor(currentNs, svc); err != nil {
		return err
	} else {
		idents.Add(IDT_PROCESSOR_NAME, ident, currentNs.Namespace, writer.FileName())
	}
	return nil
}

func (p *TProtoBuilder) buildClient(currentNs *pdl.FileNamespace, svc *pdl.FileService, idents *IdentMap) (err error) {
	writer, err := p.cxt.GetClientWriter(currentNs.Namespace, svc)
	if err != nil {
		return err
	}
	defer func() {
		if err := writer.Flush(); err != nil {
			fmt.Println("buildClient writer.Flush() error:", err)
		}
		if err := writer.Close(); err != nil {
			fmt.Println("buildClient writer.Close() error:", err)
		}
	}()

	if err = writer.WriteNamespace(currentNs.Namespace); err != nil {
		return
	}

	if err = writer.WriteClientBegin(currentNs, svc); err != nil {
		return
	}
	if err = writer.WriteClient(currentNs, svc); err != nil {
		return
	}
	if err = writer.WriteClientEnd(currentNs, svc); err != nil {
		return
	}
	return
}

func (p *TProtoBuilder) buildServer(currentNs *pdl.FileNamespace, svcs map[string]*pdl.FileService, idents *IdentMap) (err error) {
	writer, err := p.cxt.GetServerWriter(currentNs.Namespace)
	if err != nil {
		return err
	}
	defer func() {
		if err := writer.Flush(); err != nil {
			fmt.Println("buildServer writer.Flush() error:", err)
		}
		if err := writer.Close(); err != nil {
			fmt.Println("buildServer writer.Close() error:", err)
		}
	}()
	if err = writer.WriteNamespace(currentNs.Namespace); err != nil {
		return
	}
	if err = writer.WriteServicesBegin(currentNs); err != nil {
		return
	}
	if err = writer.WriteServices(currentNs, svcs); err != nil {
		return
	}
	if err = writer.WriteServicesEnd(currentNs); err != nil {
		return
	}
	return nil
}
