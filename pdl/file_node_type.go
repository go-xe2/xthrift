/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-20 11:45
* Description:
*****************************************************************/

package pdl

type TFileNodeType int8

const (
	// 数据类型节点
	FNT_DATATYPE_NODE TFileNodeType = iota
	// 数据别名节点
	FNT_TYPEDEF_NODE
	// 字段数据节点
	FNT_FIELD_NODE
	// struct结构体节点
	FNT_STRUCT_NODE
	// 接口数据节点
	FNT_METHOD_NODE
	// 服务数据节点
	FNT_SERVICE_NODE
	// 命名空间数据节点
	FNT_NAMESPACE_NODE
	// 文件数据节点
	FNT_FILE_NODE
	// 协议项目数据节点
	FNT_PROJECT_NODE
	// 未知类型数据节点
	FNT_UNKNOWN_NODE
)

func (ft TFileNodeType) String() string {
	switch ft {
	case FNT_DATATYPE_NODE:
		return "datatype_node"
	case FNT_TYPEDEF_NODE:
		return "typedef_node"
	case FNT_STRUCT_NODE:
		return "struct_node"
	case FNT_FIELD_NODE:
		return "field_node"
	case FNT_METHOD_NODE:
		return "method_node"
	case FNT_SERVICE_NODE:
		return "service_node"
	case FNT_NAMESPACE_NODE:
		return "namespace_node"
	case FNT_FILE_NODE:
		return "file_node"
	case FNT_PROJECT_NODE:
		return "project_node"
	default:
		return "unknown_node"
	}
}
