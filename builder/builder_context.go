/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-13 10:46
* Description:
*****************************************************************/

package builder

import (
	"github.com/go-xe2/xthrift/pdl"
)

type tNamespaceFiles struct {
	ns *pdl.FileNamespace
	// ident => file
	files map[string]string
}

func newNamespaceFiles(ns *pdl.FileNamespace, files map[string]string) *tNamespaceFiles {
	return &tNamespaceFiles{
		ns:    ns,
		files: files,
	}
}

type Context interface {
	GetWorkPath() string
	// 获取数据类型定义写入文件
	GetTypedefWriter(ns *pdl.FileNamespace) (TypedefCodeWriter, error)
	// 获取结构体定义文件
	GetStructWriter(ns string, stru *pdl.FileStruct) (StructCodeWriter, error)
	// 获取服务接口定义文件
	GetServiceWriter(ns string, svc *pdl.FileService) (ServiceCodeWriter, error)
	// 获取处理器，处理方法写入文件
	GetProcessorFunWriter(ns string, svc *pdl.FileService) (ProcessorFunCodeWriter, error)
	// 获取处理器写入文件
	GetProcessorWriter(ns string, svc *pdl.FileService) (ProcessorCodeWriter, error)
	// 获取客户端代码写入接口
	GetClientWriter(ns string, svc *pdl.FileService) (ClientCodeWrite, error)
	GetServerWriter(ns string) (ServerCodeWriter, error)

	GetNamespaceFiles(namespace string) map[string]string
	SetNamespaceFiles(namespace string, files map[string]string)
}
