namespace go com.mnyun.reg.admin
namespace java com.mnyun.reg.admin
namespace php com.mnyun.reg.admin
include 'com.mnyun.reg.types.thrift'
// 文件:regAdminSvc开始
typedef string str
typedef bool bl
typedef double dl
service RegSvc {
	// 下级地区
	list<com.mnyun.reg.types.RegItem> GetChildListResult(1:i32 parId);
	// 乡镇列表
	list<com.mnyun.reg.types.RegItem> GetTownList(1:i32 countyId);
	// 州市列表
	list<com.mnyun.reg.types.RegItem> GetCityList(1:i32 provinceId);
	// 修改
	bl UpdateResult(1:i32 regId,2:i32 parId,3:str name);
	// 删除
	bl RemoveResult(1:i32 regId);
	// 详情
	com.mnyun.reg.types.RegItem RegDetailResult(1:i32 regId);
	// 新增
	bl AddResult(1:i32 parId,2:str name);
	// 区县列表
	list<com.mnyun.reg.types.RegItem> GetCountList(1:i32 cityId);
	// 省份列表
	list<com.mnyun.reg.types.RegItem> GetProvincesList();
	// 地区目录树
	list<com.mnyun.reg.types.RegItem> GetRegTreeResult(1:i32 parId);
}
// 文件regAdminSvc结束
