/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-13 10:30
* Description:
*****************************************************************/

package builder

import (
	"github.com/go-xe2/xthrift/pdl"
)

type IdentType int8

const (
	// typedef数据类型标识
	IDT_TYPE_DEF_NAME IdentType = iota
	// type类型标识
	IDT_TYPE_NAME
	// 字段标识
	IDT_FIELD_NAME
	// 服务名标识
	IDT_SERVICE_NAME
	// 服务接口名标识
	IDT_METHOD_NAME
	// 处理器名标识
	IDT_PROCESSOR_NAME
	// 处理器处理方法名标识
	IDT_PROCESSOR_FUN_NAME
)

func (it IdentType) String() string {
	switch it {
	case IDT_TYPE_DEF_NAME:
		return "typedef ident"
	case IDT_TYPE_NAME:
		return "type ident"
	case IDT_FIELD_NAME:
		return "field ident"
	case IDT_SERVICE_NAME:
		return "service ident"
	case IDT_METHOD_NAME:
		return "method ident"
	case IDT_PROCESSOR_NAME:
		return "processor ident"
	case IDT_PROCESSOR_FUN_NAME:
		return "processorFun ident"
	}
	return "unknown ident"
}

type CodeWriter interface {
	FileName() string
	// 写入引用文件
	WriteNamespace(namespace string) error
	Import(namespace string, inner bool)
	// 写入到缓存
	Write(str string)
	Context() Context
	Flush() error
	Close() error
}

type TypedefCodeWriter interface {
	CodeWriter
	WriteDef(ns *pdl.FileNamespace, def *pdl.FileTypeDef) (ident string, err error)
}

type StructCodeWriter interface {
	CodeWriter
	WriteStructBegin(ns *pdl.FileNamespace, stru *pdl.FileStruct) (ident string, err error)
	WriteStruct(ns *pdl.FileNamespace, stru *pdl.FileStruct) error
	WriteStructEnd(ns *pdl.FileNamespace, stru *pdl.FileStruct) error
}

type ServiceCodeWriter interface {
	CodeWriter
	WriteServiceBegin(ns *pdl.FileNamespace, svc *pdl.FileService) (ident string, err error)
	WriteServiceMethod(ns *pdl.FileNamespace, svc *pdl.FileService, method *pdl.FileServiceMethod) (ident string, err error)
	WriteServiceEnd(ns *pdl.FileNamespace, svc *pdl.FileService) error
}

type ProcessorFunCodeWriter interface {
	CodeWriter
	WriteFunction(ns *pdl.FileNamespace, svc *pdl.FileService, method *pdl.FileServiceMethod, input, result *pdl.FileStruct) (ident string, err error)
}

type ProcessorCodeWriter interface {
	CodeWriter
	WriteProcessor(ns *pdl.FileNamespace, svc *pdl.FileService) (ident string, err error)
}

type ClientCodeWrite interface {
	CodeWriter
	WriteClientBegin(ns *pdl.FileNamespace, svc *pdl.FileService) error
	WriteClient(ns *pdl.FileNamespace, svc *pdl.FileService) error
	WriteClientEnd(ns *pdl.FileNamespace, svc *pdl.FileService) error
}

type ServerCodeWriter interface {
	CodeWriter
	WriteServicesBegin(ns *pdl.FileNamespace) error
	WriteServices(ns *pdl.FileNamespace, svcs map[string]*pdl.FileService) error
	WriteServicesEnd(ns *pdl.FileNamespace) error
}
