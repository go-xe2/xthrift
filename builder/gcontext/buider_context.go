/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 12:44
* Description:
*****************************************************************/

package gcontext

import (
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/type/xstring"
	"github.com/go-xe2/xthrift/builder"
	"github.com/go-xe2/xthrift/pdl"
	"strings"
)

type TContext struct {
	workPath   string
	moduleName string
	nsFiles    map[string]map[string]string
}

var _ builder.Context = (*TContext)(nil)

func NewContext(workPath, moduleName string) *TContext {
	return &TContext{
		workPath:   workPath,
		moduleName: moduleName,
		nsFiles:    make(map[string]map[string]string),
	}
}

func (p *TContext) GetWorkPath() string {
	return p.workPath
}

func (p *TContext) GetModuleName() string {
	return p.moduleName
}

func (p *TContext) GetTypedefWriter(ns *pdl.FileNamespace) (builder.TypedefCodeWriter, error) {
	relativeDir := strings.Replace(ns.Namespace, ".", xfile.Separator, -1)
	fileName := xfile.Join(p.workPath, relativeDir, "typedef.go")
	return NewTypedefCodeWriter(p, fileName)
}

func (p *TContext) GetStructWriter(ns string, stru *pdl.FileStruct) (builder.StructCodeWriter, error) {
	relativeDir := strings.Replace(ns, ".", xfile.Separator, -1)
	fileName := xfile.Join(p.workPath, relativeDir, xstring.Camel2UnderScore(stru.Type.TypName, "_")+".go")
	return NewStructCodeWriter(p, fileName)
}

func (p *TContext) GetServiceWriter(ns string, svc *pdl.FileService) (builder.ServiceCodeWriter, error) {
	relativeDir := strings.Replace(ns, ".", xfile.Separator, -1)
	fileName := xfile.Join(p.workPath, relativeDir, xstring.Camel2UnderScore(svc.Name, "_")+".go")
	return NewServiceWriter(p, fileName)
}

func (p *TContext) GetProcessorFunWriter(ns string, svc *pdl.FileService) (builder.ProcessorFunCodeWriter, error) {
	relativeDir := strings.Replace(ns, ".", xfile.Separator, -1)
	fileName := xfile.Join(p.workPath, relativeDir, xstring.Camel2UnderScore(svc.Name, "_")+"_processor_func.go")
	return NewProcessorFunCodeWriter(p, fileName)
}

func (p *TContext) GetProcessorWriter(ns string, svc *pdl.FileService) (builder.ProcessorCodeWriter, error) {
	relativeDir := strings.Replace(ns, ".", xfile.Separator, -1)
	fileName := xfile.Join(p.workPath, relativeDir, xstring.Camel2UnderScore(svc.Name, "_")+"_processor.go")
	return NewProcessorCodeWriter(p, fileName)
}

func (p *TContext) GetClientWriter(ns string, svc *pdl.FileService) (builder.ClientCodeWrite, error) {
	relativeDir := strings.Replace(ns, ".", xfile.Separator, -1)
	fileName := xfile.Join(p.workPath, relativeDir, xstring.Camel2UnderScore(svc.Name, "_")+"_client.go")
	return NewClientCodeWriter(p, fileName)
}

func (p *TContext) GetServerWriter(ns string) (builder.ServerCodeWriter, error) {
	relativeDir := strings.Replace(ns, ".", xfile.Separator, -1)
	fileName := xfile.Join(p.workPath, relativeDir, "server.go")
	return NewServerCodeWriter(p, fileName)
}

func (p *TContext) GetNamespaceFiles(namespace string) map[string]string {
	if mp, ok := p.nsFiles[namespace]; ok {
		return mp
	}
	return nil
}

func (p *TContext) SetNamespaceFiles(namespace string, files map[string]string) {
	if old, ok := p.nsFiles[namespace]; ok {
		for k, v := range files {
			if _, ok1 := old[k]; !ok1 {
				old[k] = v
			}
		}
		return
	}
	p.nsFiles[namespace] = files
}
