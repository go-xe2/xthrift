/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-03 16:03
* Description:
*****************************************************************/

package netstream

import (
	"github.com/go-xe2/x/core/logger"
	"github.com/go-xe2/x/os/xlog"
	"sync/atomic"
)

func (p *TStreamClient) OnReady(conn StreamConn) {

}

func (p *TStreamClient) OnRecv(conn StreamConn, data []byte) {
	defer func() {
		if e := recover(); e != nil {
			p.Log(logger.LEVEL_WARN, "OnRecv error:", e)
		}
	}()
	// 收到数据包处理
	if p.handler != nil {
		p.handler.OnRecv(conn, data)
	}
}

func (p *TStreamClient) OnCall(conn StreamConn, data []byte) ([]byte, error) {
	// 同步调用，返回数据
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
func (p *TStreamClient) OnSendTo(conn StreamConn, toConn string, data []byte) {
	// 转发数据
}

func (p *TStreamClient) OnCallTo(conn StreamConn, toConn string, data []byte) ([]byte, error) {
	return data, nil
}

func (p *TStreamClient) OnHeartbeat(conn StreamConn) {
	// 收到收跳时，回复心跳
	xlog.Debug("收到心跳包======>>")
	conn.UpdateHeartbeat(false)
}

func (p *TStreamClient) OnDisconnect(conn StreamConn) {
	select {
	case <-p.closed:
		// 人工关闭客户端，不再重连
		return
	default:
	}
	defer func() {
		if e := recover(); e != nil {
			p.Log(logger.LEVEL_WARN, "OnDisconnect error:", e)
		}
	}()
	xlog.Debug("断线准备重试连接中")
	if n := atomic.LoadInt32(&p.retryStatus); n != 0 {
		// 当前重试连接进程未完成，不进行重试处理
		xlog.Debug("重试连接协程已经运行.")
		return
	}
	xlog.Debug("准备重试连接")
	// 断开重连接
	//// 非人工关闭，尝试重连
	atomic.StoreInt32(&p.retryStatus, 1)
	defer atomic.StoreInt32(&p.retryStatus, 0)

	if !p.RetryConnect() {
		// 重试连接成功后，不触发OnDisconnect事件,只有重试之后仍然连接不上的情况触发
		if p.handler != nil {
			p.handler.OnDisconnect(conn)
		}
	} else {
		if p.handler != nil {
			p.handler.OnReconnect(conn)
		}
	}
}

func (p *TStreamClient) OnConnect(conn StreamConn) {
	select {
	case <-p.closed:
		p.closed = make(chan byte, 1)
	default:
	}
	// 客户端已经连接上, 通知服务端，客户端已就绪
	xlog.Debug("客户端连接上...")
	_ = p.conn.SendReady()
	defer func() {
		if e := recover(); e != nil {
			p.Log(logger.LEVEL_WARN, "OnConnect error:", e)
		}
	}()
	if p.handler != nil {
		p.handler.OnConnect(conn)
	}
	// 启动心跳检查
	p.conn.UpdateHeartbeat(false)
	p.heartbeatProcessLoop()
}

// 处理请求
func (p *TStreamClient) OnRequest(reqConn StreamConn, reqId string, namespace string, body []byte) {
	if p.handler != nil {
		p.handler.OnRequest(reqId, namespace, body)
	}
}

func (p *TStreamClient) OnResponse(resConn StreamConn, reqId string, body []byte) {
	if p.handler != nil {
		p.handler.OnResponse(reqId, body)
	}
}
