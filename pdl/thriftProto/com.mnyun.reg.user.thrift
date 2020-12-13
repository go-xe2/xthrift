namespace go com.mnyun.reg.user
namespace java com.mnyun.reg.user
namespace php com.mnyun.reg.user
include 'com.mnyun.reg.types.thrift'
// 文件:regUserSvc开始
typedef string str
typedef bool bl
typedef double dl
service RegSvc {
	// 下级地区
	list<com.mnyun.reg.types.RegItem> GetChildListResult(1:optional i32 parId);
	// 州市列表
	list<com.mnyun.reg.types.RegItem> GetCityList(1:i32 provinceId);
	// 区县列表
	list<com.mnyun.reg.types.RegItem> GetCountList(1:i32 cityId);
	// 省份列表
	list<com.mnyun.reg.types.RegItem> GetProvincesList();
	// 地区目录树
	list<com.mnyun.reg.types.RegItem> GetRegTreeResult(1:optional i32 parId);
	// 乡镇列表
	list<com.mnyun.reg.types.RegItem> GetTownList(1:i32 countyId);
	// 详情
	com.mnyun.reg.types.RegItem RegDetailResult(1:i32 id);
}
// 文件regUserSvc结束
