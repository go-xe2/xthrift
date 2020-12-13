namespace go com.mnyun.demo
namespace java com.mnyun.demo
namespace php com.mnyun.demo
include 'com.mnyun.types.thrift'
// 文件:demo开始
typedef string str
typedef bool bl
typedef double dl
service helloService {
	// sayHello的说明
	com.mnyun.types.helloResult sayHello(1:optional str name,2:i32 age);
}
// 文件demo结束
