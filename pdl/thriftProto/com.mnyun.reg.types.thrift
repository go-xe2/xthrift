namespace go com.mnyun.reg.types
namespace java com.mnyun.reg.types
namespace php com.mnyun.reg.types
// 文件:types开始
typedef string str
typedef bool bl
typedef double dl
// 地区资料参数
struct RegItem {
	// 分类id
	1:i32 Id;
	// 上级地区id
	2:i32 ParentId;
	// 名称
	3:str Name;
	// 层级数
	4:i8 Level;
	// 子分类数
	5:i32 ChildCount;
	// 父级id列表 以逗号分隔开
	6:str ParIds;
	// 完整地区名
	7:str Path;
	// 创建日期
	8:i64 Time;
	// 测试
	9:list<RegItem> Data;
}
// 文件types结束
