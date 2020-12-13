/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-06 10:30
* Description:
*****************************************************************/

package netstream

import (
	"github.com/go-xe2/x/core/logger"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/x/xf/log/xdefaultLogger"
	"time"
)

type TConnOptions struct {
	// 日志输出接口
	logger logger.ILogger
	// 数据发送缓存队列大小
	sendBufSize int
	// 数据接收缓存队列大小
	recvBufSize int
	// 连接超时时间
	connectTimeout time.Duration
	// 数据读取超时时间
	readTimeout time.Duration
	// 数据发送超时时间
	writeTimeout time.Duration
}

// 默认流连接参数
var DefaultConnOptions = NewConnOptions()

func NewConnOptions() *TConnOptions {
	inst := &TConnOptions{
		logger:         xdefaultLogger.New(),
		sendBufSize:    1024,
		recvBufSize:    1024,
		connectTimeout: 3 * time.Minute,
		readTimeout:    3 * time.Minute,
		writeTimeout:   3 * time.Minute,
	}
	return inst
}

func (p *TConnOptions) LoadFromMap(mp map[string]interface{}) {
	if mp == nil {
		return
	}
	if v, ok := mp["sendBufSize"]; ok {
		p.recvBufSize = t.Int(v)
	}
	if v, ok := mp["sendBufSize"]; ok {
		p.sendBufSize = t.Int(v)
	}
	if v, ok := mp["connectTimeout"]; ok {
		p.connectTimeout = t.Duration(v) * time.Second
	}
	if v, ok := mp["readTimeout"]; ok {
		p.readTimeout = t.Duration(v) * time.Second
	}
	if v, ok := mp["writeTimeout"]; ok {
		p.writeTimeout = t.Duration(v) * time.Second
	}
}

func (p *TConnOptions) SetLogger(logger logger.ILogger) *TConnOptions {
	p.logger = logger
	return p
}

func (p *TConnOptions) SetSendBufSize(size int) *TConnOptions {
	p.sendBufSize = size
	return p
}

func (p *TConnOptions) SetRecvBufSize(size int) *TConnOptions {
	p.recvBufSize = size
	return p
}

func (p *TConnOptions) SetConnectTimeout(timeout time.Duration) *TConnOptions {
	p.connectTimeout = timeout
	return p
}

func (p *TConnOptions) SetWriteTimeout(timeout time.Duration) *TConnOptions {
	p.writeTimeout = timeout
	return p
}

func (p *TConnOptions) SetReadTimeout(timeout time.Duration) *TConnOptions {
	p.readTimeout = timeout
	return p
}

func (p *TConnOptions) GetLogger() logger.ILogger {
	return p.logger
}

func (p *TConnOptions) GetSendBufSize() int {
	return p.sendBufSize
}

func (p *TConnOptions) GetRecvBufSize() int {
	return p.recvBufSize
}

func (p *TConnOptions) GetConnectTimeout() time.Duration {
	return p.connectTimeout
}

func (p *TConnOptions) GetWriteTimeout() time.Duration {
	return p.writeTimeout
}

func (p *TConnOptions) GetReadTimeout() time.Duration {
	return p.readTimeout
}
