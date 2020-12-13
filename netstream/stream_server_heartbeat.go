/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-03 17:02
* Description:
*****************************************************************/

package netstream

import (
	"github.com/go-xe2/x/core/logger"
	"sync/atomic"
	"time"
)

// 服务端心跳处理

func (p *TStreamServer) processHeartbeat() {
	atomic.StoreInt32(&p.heartbeatRun, 1)
	defer atomic.StoreInt32(&p.heartbeatRun, 0)
	maxLoss := p.options.GetAllowMaxLoss()
	speed := p.options.GetHeartbeatSpeed()
	for {
		select {
		case <-p.closed:
			return
		default:
		}
		keys := p.clients.Keys()
		if len(keys) == 0 {
			return
		}
		for _, key := range keys {
			select {
			case <-p.closed:
				return
			default:
			}
			cli := p.getClient(key)
			if cli == nil {
				continue
			}
			if cli.HeartbeatLossCount() > maxLoss {
				// 超出允许丢失心跳的的最大值,说明客户端已经断线，断开与该客户端的连接
				if e := cli.Close(); e != nil {
					p.Log(logger.LEVEL_WARN, "关闭客户端出错:", e)
				}
				continue
			}
			t := cli.Heartbeat()
			diff := time.Now().Sub(t)
			if diff > speed {
				cli.UpdateHeartbeat(true)
			}
		}
		time.Sleep(speed)
	}
}

func (p *TStreamServer) heartbeatProcessLoop() {
	if p.options.GetHeartbeatSpeed() == 0 || p.options.GetAllowMaxLoss() == 0 {
		// 不启动心跳
		return
	}
	select {
	case <-p.closed:
		return
	default:
	}
	if isRun := atomic.LoadInt32(&p.heartbeatRun); isRun != 0 {
		return
	}
	go p.processHeartbeat()
}
