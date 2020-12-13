/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-20 14:42
* Description:
*****************************************************************/

package pdl

type TPDLNamespace struct {
	container *TPDLContainer
	// 命名空间名称
	Name    string                    `json:"name"`
	Imports map[string]*TPDLNamespace `json:"imports"`
	// 命名空间中定义的数据别名
	Typedefs map[string]*FileTypeDef `json:"typedefs"`
	// 空间中定义的数据结构
	Types map[string]*FileStruct `json:"types"`
	// 空间中定义的服务
	Services map[string]*FileService `json:"services"`
}

func NewPDLNamespace(cont *TPDLContainer) *TPDLNamespace {
	return &TPDLNamespace{
		container: cont,
		Imports:   make(map[string]*TPDLNamespace),
		Typedefs:  make(map[string]*FileTypeDef),
		Types:     make(map[string]*FileStruct),
		Services:  make(map[string]*FileService),
	}
}

func (p *TPDLNamespace) SetName(name string) {
	p.Name = name
}

func (p *TPDLNamespace) AddImport(ns *TPDLNamespace) {
	if _, ok := p.Imports[ns.Name]; !ok {
		p.Imports[ns.Name] = ns
	}
}

func (p *TPDLNamespace) AddTypedef(def *FileTypeDef) {
	if _, ok := p.Typedefs[def.Name]; !ok {
		p.Typedefs[def.Name] = def
	}
}

func (p *TPDLNamespace) AddType(typ *FileStruct) {
	if _, ok := p.Types[typ.Type.TypName]; !ok {
		p.Types[typ.Type.TypName] = typ
	}
}

func (p *TPDLNamespace) AddService(svc *FileService) {
	if _, ok := p.Services[svc.Name]; !ok {
		p.Services[svc.Name] = svc
	}
}

func (p *TPDLNamespace) QryService(name string) *FileService {
	if svc, ok := p.Services[name]; ok {
		return svc
	}
	return nil
}

func (p *TPDLNamespace) QryTypedef(name string) *FileTypeDef {
	if def, ok := p.Typedefs[name]; ok {
		return def
	}
	return nil
}

func (p *TPDLNamespace) QryType(typeName string) *FileStruct {
	if typ, ok := p.Types[typeName]; ok {
		return typ
	}
	return nil
}

func (p *TPDLNamespace) QryMethod(svcName string, method string) (*FileService, *FileServiceMethod) {
	if svc := p.QryService(svcName); svc != nil {
		if m := svc.QryMethod(method); m != nil {
			return svc, m
		}
		return svc, nil
	}
	return nil, nil
}
