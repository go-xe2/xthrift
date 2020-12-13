/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-03 10:30
* Description:
*****************************************************************/

package netstream

type StreamHandler interface {
	OnRecv(conn StreamConn, data []byte)
	OnCall(conn StreamConn, data []byte) (result []byte, err error)
	OnConnect(conn StreamConn)
	// 尝试重连接成功
	OnReconnect(conn StreamConn)
	OnDisconnect(conn StreamConn)
}

type ClientStreamHandler interface {
	StreamHandler
	OnRequest(reqId string, namespace string, data []byte)
	OnResponse(reqId string, data []byte)
}

type ServerStreamHandler interface {
	StreamHandler
	OnSendTo(conn StreamConn, target StreamConn, data []byte)
	OnCallTo(conn StreamConn, target StreamConn, data []byte) (result []byte, err error)
	OnRequest(reqConn StreamConn, reqId string, namespace string, data []byte)
}
