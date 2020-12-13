/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-25 09:49
* Description:
*****************************************************************/

package rpcPoint

import (
	"errors"
	"github.com/go-xe2/x/container/xarray"
	"github.com/go-xe2/x/container/xpool"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/x/utils/xrand"
	"github.com/go-xe2/xthrift/regcenter"
	"time"
)

// 访问服务客户端线程池
type tHostClientPool struct {
	server           *TEndPointServer
	pool             *xpool.TPool
	hosts            []*regcenter.THostStoreToken
	readTimeout      time.Duration
	writeTimeout     time.Duration
	connectTimeout   time.Duration
	ConnectFailRetry time.Duration
	keepAlive        time.Duration
	// 心跳频率
	heartbeat time.Duration
	// 允许丢失心跳的最大次数
	heartbeatLoss int
}

func newHostClientPool(server *TEndPointServer, hosts []*regcenter.THostStoreToken, readTimeout, writeTimeout, connectTimeout, failReTry, keeyAlive, heartbeat time.Duration, hearbeatLoss int) *tHostClientPool {
	inst := &tHostClientPool{
		hosts:            hosts,
		server:           server,
		readTimeout:      readTimeout,
		writeTimeout:     writeTimeout,
		connectTimeout:   connectTimeout,
		ConnectFailRetry: failReTry,
		keepAlive:        keeyAlive,
		heartbeat:        heartbeat,
		heartbeatLoss:    hearbeatLoss,
	}
	inst.pool = xpool.New(heartbeat, inst.poolNewFun, inst.poolExpireFun)
	return inst
}

func (p *tHostClientPool) poolNewFun() (interface{}, error) {
	tmpArr := xarray.New(true)
	tryItems := make([]*regcenter.THostStoreToken, 0)
	now := time.Now()
	for _, host := range p.hosts {
		if host.CanUse(now) {
			tmpArr.PushLeft(host)
		} else {
			tryItems = append(tryItems, host)
		}
	}
	idx := 0
	if tmpArr.Len() == 0 {
		// 无可用资源时，尝试使用连接失败的服务
		for _, h := range tryItems {
			tmpArr.PushLeft(h)
		}
	}
	if tmpArr.Len() > 1 {
		idx = xrand.N(0, tmpArr.Len()-1)
	}
	for {
		if tmpArr.Len() <= 0 {
			return nil, errors.New("没有可用服务")
		}
		cur := tmpArr.Get(idx).(*regcenter.THostStoreToken)
		c, err := p.newHostClient(cur)
		if err == nil && c.IsOpen() {
			cur.SetConnectSuccess()
			c.expire = time.Now().Add(p.keepAlive)
			return c, nil
		}
		// 失败的节点，2分钟后再重试使用链接
		cur.SetConnectFail(time.Now().Add(p.ConnectFailRetry))
		tmpArr.Remove(idx)
		idx = 0
		if tmpArr.Len() > 1 {
			idx = xrand.N(0, tmpArr.Len())
		}
	}
}

func (p *tHostClientPool) Get() *tInnerClient {
	v, err := p.pool.Get()
	if err != nil || v == nil {
		return nil
	}
	return v.(*tInnerClient)
}

func (p *tHostClientPool) Put(item *tInnerClient) {
	p.pool.Put(item)
}

func (p *tHostClientPool) newHostClient(host *regcenter.THostStoreToken) (*tInnerClient, error) {
	return newInnerClient(p, p.server.inProtoFac, host.Host, host.Port, p.writeTimeout, p.readTimeout, p.connectTimeout, p.heartbeatLoss)
}

func (p *tHostClientPool) poolExpireFun(item interface{}) {
	c := item.(*tInnerClient)
	if !c.IsOpen() {
		// 连接已经断开
		return
	}
	now := time.Now()
	if now.After(c.expire) {
		// 连接空闲超时，关闭连接
		xlog.Debug("连接空闲超时，断开链接")
		if err := c.Close(); err != nil {
			xlog.Error(err)
		}
	} else {
		// 放回线程池
		// 心跳检查
		if p.heartbeat > 0 {
			now := time.Now()
			if p.heartbeatLoss == 0 {
				p.heartbeatLoss = 3
			}
			if now.After(c.lastAlive.Add(p.heartbeat)) {
				if p.heartbeatLoss > 0 && c.heartFailCount > p.heartbeatLoss {
					// 心跳丢失次数过多, 关闭当前链接
					if err := c.Close(); err != nil {
						xlog.Error(err)
					}
					return
				}
				// 发送心跳
				c.sendHeartbeat()
			}
		}
		// 放回线程池
		p.Put(c)
	}
}
