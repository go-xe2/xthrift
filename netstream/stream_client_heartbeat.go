/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-04 16:55
* Description:
*****************************************************************/

package netstream

import (
	"github.com/go-xe2/x/core/logger"
	"github.com/go-xe2/x/os/xlog"
	"sync/atomic"
	"time"
)

func (p *TStreamClient) heartbeatProcessLoop() {
	speed := p.options.GetHeartbeatSpeed()
	maxLoss := p.options.GetAllowMaxLoss()
	xlog.Debug("准备发送心跳,speed:", speed, ", maxLoss:", maxLoss)
	if speed == 0 || maxLoss == 0 {
		// 未设置跳心速率，不启用心跳检查
		return
	}
	if isRun := atomic.LoadInt32(&p.heartbeatIsRun); isRun != 0 {
		return
	}
	atomic.StoreInt32(&p.heartbeatIsRun, 1)
	defer atomic.StoreInt32(&p.heartbeatIsRun, 0)
	go func() {
		for {
			if p.conn.isStop() {
				// 已经关闭服务
				xlog.Debug("已关闭服务，退出心跳检测.")
				return
			}
			conn := p.conn
			xlog.Debug("丢失", conn.HeartbeatLossCount(), "次心跳")
			if conn.HeartbeatLossCount() > maxLoss {
				// 心跳包丢失超出了最大限制，断开连接尝试重连
				// 当非人工关闭时，在onDisconnect事件中会调用RetryConnect,所以该处不用调用RetryConnect
				conn.innerClose()
				//p.RetryConnect()
				p.Log(logger.LEVEL_INFO, "长时间未收到数据，已断开连接")
				return
			}
			t := conn.Heartbeat()
			diff := time.Now().Sub(t)
			if diff > speed {
				conn.UpdateHeartbeat(true)
				if e := conn.SendHeartbeat(); e != nil {
					p.Log(logger.LEVEL_WARN, "发送心跳包出错:", e)
				}
			}
			time.Sleep(speed)
		}
	}()
}

func (p *TStreamClient) RetryConnect() bool {
	maxTryCount := p.options.GetTryConnectCount()
	trySpeed := p.options.GetTryConnectSpeed()
	if trySpeed <= 0 {
		return false
	}
	for {
		n := p.tryCount + 1
		if maxTryCount > 0 {
			p.Log(logger.LEVEL_INFO, "尝试第", n, "次连接, maxCount:", maxTryCount)
		} else {
			p.Log(logger.LEVEL_INFO, "尝试第", n, "次连接")
		}
		if e := p.conn.ReOpen(); e != nil {
			p.tryCount++
			p.Log(logger.LEVEL_DEBU, "第", n, "次连接失败:", e)
			// maxTryCount为0时表达一直测试连接
			if maxTryCount > 0 && p.tryCount >= maxTryCount {
				p.Log(logger.LEVEL_ERRO, "尝试连接", p.options.GetAllowMaxLoss(), "次失败，请检查服务端是否正常运行")
				// 重连接不成功，关闭服务
				close(p.closed)
				return false
			}
		} else {
			p.tryCount = 0
			return true
		}
		time.Sleep(trySpeed)
	}
}
