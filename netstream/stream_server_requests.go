/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-07 11:44
* Description:
*****************************************************************/

package netstream

import "time"

type tStreamRequestStore struct {
	server *TStreamServer
	reqId  string
	// 连接ID
	connId string
	// 会话超时时间
	expireTime time.Time
}

func newStreamRequest(server *TStreamServer, seqId string, reqConn StreamConn, timeout time.Duration) *tStreamRequestStore {
	inst := &tStreamRequestStore{
		reqId:      seqId,
		server:     server,
		connId:     reqConn.Id(),
		expireTime: time.Now().Add(timeout),
	}
	return inst.checkExpire()
}

func (p *tStreamRequestStore) remove() {
	if p.server != nil {
		p.server.removeRequest(p.reqId)
	}
}

func (p *tStreamRequestStore) checkExpire() *tStreamRequestStore {
	go func() {
		time.Sleep(p.expireTime.Sub(time.Now()))
		p.remove()
	}()
	return p
}

func (p *tStreamRequestStore) ConnId() string {
	return p.connId
}

// 移出请求
func (p *TStreamServer) removeRequest(reqId string) {
	p.requests.Remove(reqId)
}

// 存储请求会话
func (p *TStreamServer) storeRequest(reqId string, reqConn StreamConn) {
	p.requests.Set(reqId, newStreamRequest(p, reqId, reqConn, p.options.GetReadTimeout()))
}

func (p *TStreamServer) restoreRequest(reqId string) StreamConn {
	req := p.requests.Get(reqId)
	if req == nil {
		return nil
	}
	connId := req.(*tStreamRequestStore).ConnId()
	clientConn := p.getClient(connId)
	return clientConn
}
