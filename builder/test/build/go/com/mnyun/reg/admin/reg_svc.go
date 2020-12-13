package admin

import (
	"github.com/go-xe2/xthrift/builder/test/build/go/com/mnyun/reg/types"
)

type RegSvc interface {
	// 修改
	UpdateResult(regId int32, parId int32, name string) bool
	// 乡镇列表
	GetTownList(countyId int32) []*types.RegItem
	// 州市列表
	GetCityList(provinceId int32) []*types.RegItem
	// 地区目录树
	GetRegTreeResult(parId int32) []*types.RegItem
	// 删除
	RemoveResult(regId int32) bool
	// 省份列表
	GetProvincesList() []*types.RegItem
	// 区县列表
	GetCountList(cityId int32) []*types.RegItem
	// 新增
	AddResult(parId int32, name string) bool
	// 下级地区
	GetChildListResult(parId int32) []*types.RegItem
	// 详情
	RegDetailResult(regId int32) *types.RegItem
}
