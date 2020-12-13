/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-16 10:31
* Description:
*****************************************************************/

package pdl

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/xf/ef/xqi"
)

type TStructFieldInfo struct {
	Id     int16
	FdType thrift.TType
	Setter func(obj DynamicStruct, val interface{}) bool
}

type DynamicStruct interface {
	NewInstance() DynamicStruct
	AllFields() map[string]*TStructFieldInfo
	FieldNameMaps() map[string]string
	SetFieldValue(fdName string, val interface{}) bool
	AssignFromMap(mp map[string]interface{}) bool
	SliceFromMaps(mps []map[string]interface{}) []thrift.TStruct
	AssignFromDataSet(ds xqi.Dataset) bool
	SliceFromDataSet(ds xqi.Dataset) []thrift.TStruct
}

func NewStructFieldInfo(id int16, fdType thrift.TType, setter func(obj DynamicStruct, val interface{}) bool) *TStructFieldInfo {
	return &TStructFieldInfo{
		Id:     id,
		FdType: fdType,
		Setter: setter,
	}
}
