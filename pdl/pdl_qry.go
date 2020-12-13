/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-25 10:50
* Description:
*****************************************************************/

package pdl

type PDLQuery interface {
	GetServiceByFullName(fullName string) *FileService
	QryService(fullName string) (*TPDLNamespace, *FileService)
	QryServiceByNS(namespace string, svcName string) (*TPDLNamespace, *FileService)
	QryTypedef(fullName string) (*TPDLNamespace, *FileTypeDef)
	QryTypeDefByNS(namespace string, defName string) (*TPDLNamespace, *FileTypeDef)
	QryType(fullName string) (*TPDLNamespace, *FileStruct)
	QryTypeByNS(namespace string, typName string) (*TPDLNamespace, *FileStruct)
	QryMethod(svcFullName string, methodName string) (*TPDLNamespace, *FileService, *FileServiceMethod)
	QryMethodByNS(namespace string, svcName string, methodName string) (*TPDLNamespace, *FileService, *FileServiceMethod)
	QryServices(namespace string) (*TPDLNamespace, map[string]*FileService)
	AllNamespaces() []string
	QryNamespace(namespace string) *TPDLNamespace
}
