package user

import (
	"github.com/go-xe2/xthrift/builder/test/build/go/com/mnyun/reg/types"
)

type RegSvc interface {
	// 下级地区
	GetChildListResult(parId *int32) []*types.RegItem
	// 州市列表
	GetCityList(provinceId int32) []*types.RegItem
	// 区县列表
	GetCountList(cityId int32) []*types.RegItem
	// 省份列表
	GetProvincesList() []*types.RegItem
	// 地区目录树
	GetRegTreeResult(parId *int32) []*types.RegItem
	// 乡镇列表
	GetTownList(countyId int32) []*types.RegItem
	// 详情
	RegDetailResult(id int32) *types.RegItem
}
