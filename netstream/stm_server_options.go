/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-06 10:10
* Description:
*****************************************************************/

package netstream

import (
	"time"
)

type TStmServerOptions struct {
	*TConnOptions
	// 心跳检查频率
	heartbeatSpeed time.Duration
	// 允许心跳丢失的最大次数
	allowMaxLoss int
	// 是否启用调试输出
	enableDebug bool
	// 客户连接池大小
	clientPoolSize int
}

// 默认流服务端参数
var DefaultStmServerOptions = NewStmServerOptions()

func NewStmServerOptions() *TStmServerOptions {
	return &TStmServerOptions{
		TConnOptions:   NewConnOptions(),
		heartbeatSpeed: 2 * time.Minute,
		allowMaxLoss:   5,
		enableDebug:    false,
	}
}

func (p *TStmServerOptions) SetHeartbeatSpeed(speed time.Duration) *TStmServerOptions {
	p.heartbeatSpeed = speed
	return p
}

func (p *TStmServerOptions) SetAllowMaxLoss(maxLoss int) *TStmServerOptions {
	p.allowMaxLoss = maxLoss
	return p
}

func (p *TStmServerOptions) SetEnableDebug(debug bool) *TStmServerOptions {
	p.enableDebug = debug
	return p
}

func (p *TStmServerOptions) GetHeartbeatSpeed() time.Duration {
	return p.heartbeatSpeed
}

func (p *TStmServerOptions) GetAllowMaxLoss() int {
	return p.allowMaxLoss
}

func (p *TStmServerOptions) GetEnableDebug() bool {
	return p.enableDebug
}
