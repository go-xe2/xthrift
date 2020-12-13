/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-03 12:02
* Description:
*****************************************************************/

package netstream

import (
	"fmt"
	"github.com/go-xe2/x/core/logger"
)

func (p *TStreamServer) OnReady(conn StreamConn) {
}

func (p *TStreamServer) OnRecv(conn StreamConn, data []byte) {
	defer func() {
		if e := recover(); e != nil {
			p.Log(logger.LEVEL_WARN, "OnRecv error:", e)
		}
	}()
	if p.handler != nil {
		p.handler.OnRecv(conn, data)
	}
}

func (p *TStreamServer) OnCall(conn StreamConn, data []byte) ([]byte, error) {
	defer func() {
		if e := recover(); e != nil {
			p.Log(logger.LEVEL_WARN, "OnCall error:", e)
		}
	}()
	if p.handler != nil {
		return p.handler.OnCall(conn, data)
	}
	return data, nil
}

// 收到消息并回复
func (p *TStreamServer) OnSendTo(conn StreamConn, toConn string, data []byte) {
	// 转发数据
	targetConn := p.getClient(toConn)
	if targetConn == nil {
		return
	}
	if p.handler != nil {
		(func() {
			defer func() {
				if e := recover(); e != nil {
					p.Log(logger.LEVEL_WARN, "OnSendTo error:", e)
				}
			}()
			if h1, ok := p.handler.(ServerStreamHandler); ok {
				h1.OnSendTo(conn, targetConn, data)
				return
			}
		})()
	}
	if _, err := targetConn.Send(data); err != nil {
		p.Log(logger.LEVEL_WARN, "OnSendTo error:", err)
	}
}

func (p *TStreamServer) OnCallTo(conn StreamConn, toConn string, data []byte) ([]byte, error) {
	targetConn := p.getClient(toConn)
	if targetConn == nil {
		return nil, fmt.Errorf("客户端%s不在线", toConn)
	}
	if p.handler != nil {
		defer func() {
			if e := recover(); e != nil {
				p.Log(logger.LEVEL_WARN, "OnCallTo error:", e)
			}
		}()
		if h1, ok := p.handler.(ServerStreamHandler); ok {
			return h1.OnCallTo(conn, targetConn, data)
		}
	}
	// 默认处理方式
	return targetConn.Call(data)
}

func (p *TStreamServer) OnHeartbeat(conn StreamConn) {
	// 更新客户端心跳时间并回复
	conn.UpdateHeartbeat(false)
	if e := conn.SendHeartbeat(); e != nil {
		p.Log(logger.LEVEL_WARN, "回复客户端[", conn.Id(), "]心跳包出错")
	}
}

func (p *TStreamServer) OnDisconnect(conn StreamConn) {
	id := conn.Id()
	p.removeClient(id)
	if p.handler != nil {
		defer func() {
			if e := recover(); e != nil {
				p.Log(logger.LEVEL_WARN, "OnDisconnect error:", e)
			}
		}()
		p.handler.OnDisconnect(conn)
	}
}

func (p *TStreamServer) OnConnect(conn StreamConn) {
	p.addClient(conn)
	if p.handler != nil {
		defer func() {
			if e := recover(); e != nil {
				p.Log(logger.LEVEL_WARN, "OnConnect error:", e)
			}
		}()
		p.handler.OnConnect(conn)
	}
}

func (p *TStreamServer) OnRequest(reqConn StreamConn, reqId string, namespace string, body []byte) {
	p.storeRequest(reqId, reqConn)
	if p.handler != nil {
		p.handler.OnRequest(reqConn, reqId, namespace, body)
	}
}

func (p *TStreamServer) OnResponse(resConn StreamConn, reqId string, body []byte) {
	reqConn := p.restoreRequest(reqId)
	if reqConn == nil {
		// 请求超时或客户端已断线
		return
	}
	if err := reqConn.Response(reqId, body); err != nil {
		p.Log(logger.LEVEL_WARN, "转发回复数据出错:", err)
	}
}
