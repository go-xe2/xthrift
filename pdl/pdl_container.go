/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-20 14:40
* Description:
*****************************************************************/

package pdl

import (
	"errors"
	"fmt"
	"github.com/go-xe2/x/type/xstring"
	"sort"
	"strings"
)

// 运行时协议容器
type TPDLContainer struct {
	// key: projectName, value: (key:namespace, value: service)
	projects   map[string]map[string]*TPDLNamespace
	Namespaces map[string]*TPDLNamespace
	Services   map[string]*FileService
}

func NewPDLContainer() *TPDLContainer {
	return &TPDLContainer{
		Namespaces: make(map[string]*TPDLNamespace),
		Services:   make(map[string]*FileService),
		projects:   make(map[string]map[string]*TPDLNamespace),
	}
}

func (p *TPDLContainer) GetServiceByFullName(fullName string) *FileService {
	if svc, ok := p.Services[fullName]; ok {
		return svc
	}
	return nil
}

func (p *TPDLContainer) GetProjectNamespaces(projName string) map[string]*TPDLNamespace {
	return p.projects[projName]
}

func (p *TPDLContainer) QryService(fullName string) (*TPDLNamespace, *FileService) {
	if strings.Index(fullName, ".") > 0 {
		ns, dtn := NamespaceLastName(fullName)
		return p.QryServiceByNS(ns, dtn)
	}
	for _, ns := range p.Namespaces {
		if dt := ns.QryService(fullName); dt != nil {
			return ns, dt
		}
	}
	return nil, nil
}

func (p *TPDLContainer) QryServiceByNS(namespace string, svcName string) (*TPDLNamespace, *FileService) {
	if ns, ok := p.Namespaces[namespace]; ok {
		if dt := ns.QryService(svcName); dt != nil {
			return ns, dt
		}
		return ns, nil
	}
	return nil, nil
}

func (p *TPDLContainer) QryTypedef(fullName string) (*TPDLNamespace, *FileTypeDef) {
	if strings.Index(fullName, ".") > 0 {
		ns, dtName := NamespaceLastName(fullName)
		return p.QryTypeDefByNS(ns, dtName)
	}
	for _, ns := range p.Namespaces {
		if dt := ns.QryTypedef(fullName); dt != nil {
			return ns, dt
		}
	}
	return nil, nil
}

func (p *TPDLContainer) QryTypeDefByNS(namespace string, defName string) (*TPDLNamespace, *FileTypeDef) {
	if ns, ok := p.Namespaces[namespace]; ok {
		if v := ns.QryTypedef(defName); v != nil {
			return ns, v
		}
		return ns, nil
	}
	return nil, nil
}

func (p *TPDLContainer) QryType(fullName string) (*TPDLNamespace, *FileStruct) {
	if strings.Index(fullName, ".") > 0 {
		ns, dt := NamespaceLastName(fullName)
		return p.QryTypeByNS(ns, dt)
	}
	for _, ns := range p.Namespaces {
		if dt := ns.QryType(fullName); dt != nil {
			return ns, dt
		}
	}
	return nil, nil
}

func (p *TPDLContainer) QryTypeByNS(namespace string, typName string) (*TPDLNamespace, *FileStruct) {
	if ns, ok := p.Namespaces[namespace]; ok {
		if dt := ns.QryType(typName); dt != nil {
			return ns, dt
		}
		return ns, nil
	}
	return nil, nil
}

func (p *TPDLContainer) QryMethod(svcFullName string, methodName string) (*TPDLNamespace, *FileService, *FileServiceMethod) {
	ns, svc := p.QryService(svcFullName)
	if svc != nil {
		if dt := svc.QryMethod(methodName); dt != nil {
			return ns, svc, dt
		}
		return ns, svc, nil
	}
	return nil, nil, nil
}

func (p *TPDLContainer) QryMethodByNS(namespace string, svcName string, methodName string) (*TPDLNamespace, *FileService, *FileServiceMethod) {
	ns, svc := p.QryServiceByNS(namespace, svcName)
	if svc != nil {
		if dt := svc.QryMethod(methodName); dt != nil {
			return ns, svc, dt
		}
		return ns, svc, nil
	}
	return nil, nil, nil
}

func (p *TPDLContainer) QryServices(namespace string) (*TPDLNamespace, map[string]*FileService) {
	if ns, ok := p.Namespaces[namespace]; ok {
		return ns, ns.Services
	}
	return nil, nil
}

func (p *TPDLContainer) QryNamespace(namespace string) *TPDLNamespace {
	if ns, ok := p.Namespaces[namespace]; ok {
		return ns
	}
	return nil
}

func (p *TPDLContainer) AllNamespace() []string {
	result := make([]string, len(p.Namespaces))
	i := 0
	for k := range p.Namespaces {
		result[i] = k
		i++
	}
	sort.Slice(result, func(i, j int) bool {
		return strings.Compare(result[i], result[j]) < 0
	})
	return result
}

func (p *TPDLContainer) IsInstall(projName string) bool {
	if _, ok := p.projects[projName]; ok {
		return true
	}
	return false
}

func (p *TPDLContainer) Install(proj *FileProject) error {
	if proj == nil {
		return nil
	}
	if err := proj.Check(); err != nil {
		return err
	}
	if p.IsInstall(proj.GetProjectName()) {
		return errors.New("项目已经安装")
	}
	projectNs := make(map[string]*TPDLNamespace)
	p.projects[proj.GetProjectName()] = projectNs
	services := proj.GetServices()
	for _, svc := range services {
		p.Services[fmt.Sprintf("%s.%s", svc.Namespace, xstring.LcFirst(svc.Name))] = svc
	}
	for _, ns := range proj.Namespaces {
		target := NewPDLNamespace(p)
		if err := p.installProjectNS(proj, target, ns); err != nil {
			return err
		}
		target.SetName(ns.Namespace)
		projectNs[ns.Namespace] = target
		p.Namespaces[ns.Namespace] = target
	}
	return nil
}

func (p *TPDLContainer) installProjectNS(proj *FileProject, target *TPDLNamespace, src *FileNamespace) error {
	for _, f := range src.Files {
		if err := p.installProjectFile(proj, target, src, f); err != nil {
			return err
		}
	}
	return nil
}

func (p *TPDLContainer) installProjectFile(proj *FileProject, target *TPDLNamespace, srcNS *FileNamespace, f *FileData) error {
	// 安装引用的空间
	for _, im := range f.Imports {
		if _, ok := target.Imports[im]; ok {
			continue
		}
		importNs, ok := proj.Namespaces[im]
		if !ok {
			return fmt.Errorf("项目中缺少%s空间的定义", im)
		}
		importPdlNs := NewPDLNamespace(p)
		err := p.installProjectNS(proj, importPdlNs, importNs)
		if err != nil {
			return err
		}
		target.AddImport(importPdlNs)
	}
	// 安装类型别名定义
	for _, def := range f.Typedefs {
		target.AddTypedef(def)
	}
	// 安装类型定义
	for _, t := range f.Types {
		target.AddType(t)
	}
	// 安装服务
	for _, svc := range f.Services {
		target.AddService(svc)
	}
	return nil
}

func (p *TPDLContainer) Uninstall(projectName string) error {
	if !p.IsInstall(projectName) {
		return errors.New("项目未安装")
	}
	projNamespaces := p.projects[projectName]
	for k := range projNamespaces {
		ns, ok := p.Namespaces[k]
		if !ok {
			continue
		}
		services := ns.Services
		for k := range services {
			delete(p.Services, k)
		}
		delete(p.Namespaces, k)
	}
	delete(p.projects, projectName)
	return nil
}
