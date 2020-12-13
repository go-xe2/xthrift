/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-06 10:52
* Description:
*****************************************************************/

package netstream

import (
	"github.com/go-xe2/x/type/t"
	"time"
)

type TStmClientOptions struct {
	*TConnOptions
	// 心跳检查频率
	heartbeatSpeed time.Duration
	// 允许心跳丢失的最大次数
	allowMaxLoss int
	// 是否启用调试输出
	enableDebug bool
	// 断线重试连接次数
	tryConnectCount int
	// 断线重试连接频率
	tryConnectSpeed time.Duration
}

var DefaultStmClientOptions = NewStmClientOptions()

func NewStmClientOptions() *TStmClientOptions {
	return &TStmClientOptions{
		TConnOptions:    NewConnOptions(),
		heartbeatSpeed:  2 * time.Minute,
		allowMaxLoss:    5,
		enableDebug:     false,
		tryConnectCount: 5,
		tryConnectSpeed: 2 * time.Minute,
	}
}

func (p *TStmClientOptions) LoadFromMap(mp map[string]interface{}) {
	if mp == nil {
		return
	}
	if p.TConnOptions == nil {
		p.TConnOptions = NewConnOptions()
	}
	p.TConnOptions.LoadFromMap(mp)
	if v, ok := mp["allowMaxLoss"]; ok {
		p.allowMaxLoss = t.Int(v)
	}
	if v, ok := mp["heartbeat"]; ok {
		p.heartbeatSpeed = t.Duration(v) * time.Second
	}
	if v, ok := mp["enableDebug"]; ok {
		p.enableDebug = t.Bool(v)
	}
	if v, ok := mp["tryConnectCount"]; ok {
		p.tryConnectCount = t.Int(v)
	}
	if v, ok := mp["tryConnectSpeed"]; ok {
		p.tryConnectSpeed = t.Duration(v) * time.Second
	}
}

func (p *TStmClientOptions) SetHeartbeatSpeed(speed time.Duration) *TStmClientOptions {
	p.heartbeatSpeed = speed
	return p
}

func (p *TStmClientOptions) SetAllowMaxLoss(maxLoss int) *TStmClientOptions {
	p.allowMaxLoss = maxLoss
	return p
}

func (p *TStmClientOptions) SetEnableDebug(debug bool) *TStmClientOptions {
	p.enableDebug = debug
	return p
}

func (p *TStmClientOptions) SetTryConnectCount(count int) *TStmClientOptions {
	p.tryConnectCount = count
	return p
}

func (p *TStmClientOptions) SetTryConnectSpeed(speed time.Duration) *TStmClientOptions {
	p.tryConnectSpeed = speed
	return p
}

func (p *TStmClientOptions) GetHeartbeatSpeed() time.Duration {
	return p.heartbeatSpeed
}

func (p *TStmClientOptions) GetAllowMaxLoss() int {
	return p.allowMaxLoss
}

func (p *TStmClientOptions) GetEnableDebug() bool {
	return p.enableDebug
}

func (p *TStmClientOptions) GetTryConnectCount() int {
	return p.tryConnectCount
}

func (p *TStmClientOptions) GetTryConnectSpeed() time.Duration {
	return p.tryConnectSpeed
}
