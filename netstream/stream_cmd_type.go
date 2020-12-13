/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-03 11:36
* Description:
*****************************************************************/

package netstream

type StreamCmdType int8

const (
	// 未知指令
	StreamCmdUnknown StreamCmdType = iota
	StreamCmdConnect
	StreamCmdDisconnect
	StreamCmdHeartbeat
	StreamCmdSend
	// 呼叫
	StreamCmdCall
	// 呼叫返回
	StreamCmdCallReply
	StreamCmdSendTo
	// 向指定客户端呼叫
	StreamCmdCallTo
	// 客户端准备就绪
	StreamCmdClientReady
	// 客户端id消息
	StreamCmdClientId
	// 请求指令
	StreamCmdRequest
	// 回复指定
	StreamCmdResponse
	// 无效指令，如果添加新的指使该条要放在最后，以便于判断指令是否有效
	StreamCmdInvalid
)

func (scmd StreamCmdType) String() string {
	switch scmd {
	case StreamCmdUnknown:
		return "StreamCmdUnknown"
	case StreamCmdConnect:
		return "StreamCmdConnect"
	case StreamCmdDisconnect:
		return "StreamCmdDisconnect"
	case StreamCmdHeartbeat:
		return "StreamCmdHeartbeat"
	case StreamCmdSend:
		return "StreamCmdSend"
	case StreamCmdCall:
		return "StreamCmdCall"
	case StreamCmdCallReply:
		return "StreamCmdCallReply"
	case StreamCmdSendTo:
		return "StreamCmdSendTo"
	case StreamCmdCallTo:
		return "StreamCmdCallTo"
	case StreamCmdClientReady:
		return "StreamCmdClientReady"
	case StreamCmdClientId:
		return "StreamCmdClientId"
	case StreamCmdRequest:
		return "StreamCmdRequest"
	case StreamCmdResponse:
		return "StreamCmdResponse"
	}
	return "StreamCmdUnknown"
}
