/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-26 17:15
* Description:
*****************************************************************/

package regcenter

import "github.com/go-xe2/xthrift/pdl"

var _ pdl.PDLQuery = (*TPDLStore)(nil)

func (p *TPDLStore) QryService(fullName string) (*pdl.TPDLNamespace, *pdl.FileService) {
	return p.nsContainer.QryService(fullName)
}

func (p *TPDLStore) GetServiceByFullName(fullName string) *pdl.FileService {
	return p.nsContainer.GetServiceByFullName(fullName)
}

func (p *TPDLStore) QryServiceByNS(namespace string, svcName string) (*pdl.TPDLNamespace, *pdl.FileService) {
	return p.nsContainer.QryServiceByNS(namespace, svcName)
}

func (p *TPDLStore) QryTypedef(fullName string) (*pdl.TPDLNamespace, *pdl.FileTypeDef) {
	return p.nsContainer.QryTypedef(fullName)
}

func (p *TPDLStore) QryTypeDefByNS(namespace string, defName string) (*pdl.TPDLNamespace, *pdl.FileTypeDef) {
	return p.nsContainer.QryTypeDefByNS(namespace, defName)
}

func (p *TPDLStore) QryType(fullName string) (*pdl.TPDLNamespace, *pdl.FileStruct) {
	return p.nsContainer.QryType(fullName)
}

func (p *TPDLStore) QryTypeByNS(namespace string, typName string) (*pdl.TPDLNamespace, *pdl.FileStruct) {
	return p.nsContainer.QryTypeByNS(namespace, typName)
}

func (p *TPDLStore) QryMethod(svcFullName string, methodName string) (*pdl.TPDLNamespace, *pdl.FileService, *pdl.FileServiceMethod) {
	return p.nsContainer.QryMethod(svcFullName, methodName)
}

func (p *TPDLStore) QryMethodByNS(namespace string, svcName string, methodName string) (*pdl.TPDLNamespace, *pdl.FileService, *pdl.FileServiceMethod) {
	return p.nsContainer.QryMethodByNS(namespace, svcName, methodName)
}

func (p *TPDLStore) QryServices(namespace string) (*pdl.TPDLNamespace, map[string]*pdl.FileService) {
	return p.nsContainer.QryServices(namespace)
}

func (p *TPDLStore) AllNamespaces() []string {
	return p.nsContainer.AllNamespace()
}

func (p *TPDLStore) QryNamespace(namespace string) *pdl.TPDLNamespace {
	return p.nsContainer.QryNamespace(namespace)
}

func (p *TPDLStore) AllFiles() map[string]string {
	return p.projectFiles
}
