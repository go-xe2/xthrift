namespace go com.mnyun.types
namespace java com.mnyun.types
namespace php com.mnyun.types
// 文件:types开始
typedef string str
typedef bool bl
typedef double dl
typedef i8 int
typedef list<map<str,str>> rows
struct helloData {
	1:optional str name;
	2:bl sex;
}
struct helloResult {
	1:i32 status;
	2:str msg;
	3:helloData data;
}
// 文件types结束
