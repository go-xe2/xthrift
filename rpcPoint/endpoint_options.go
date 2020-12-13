/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-18 14:59
* Description:
*****************************************************************/

package rpcPoint

import (
	"errors"
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/netstream"
	"time"
)

type TOptions struct {
	// http监听地址
	HttpAddr string `json:"httpAddr"`
	// 内部服务调用地址 `json
	ThriftAddr string `json:"thriftAddr"`
	// 协议存放路径
	PDLPath string `json:"pdlPath"`
	// 服务地址存放目录
	HostPath string `json:"hostPath"`
	// 路由器根路径
	BaseRouter    string `json:"baseRouter"`
	BaseNamespace string `json:"baseNamespace"`
	// 监听host文件变动
	WatchHostChanged bool `json:"watchHostChanged"`
	// 监听协议文件变动
	WatchPDLChanged bool `json:"watchPDLChanged"`
	// 协议文件后掇
	PDLExt string `json:"pdlExt"`
	// host文件后掇
	HostExt string `json:"hostExt"`
	// 发送数据超时时间(秒)
	WriteTimeout time.Duration `json:"writeTimeout"`
	// 接收数据超时时间(秒)
	ReadTimeout time.Duration `json:"readTimeout"`
	// 与服务器连接超时时间(秒)
	ConnectTimeout time.Duration `json:"connectTimeout"`
	// 线程池中客户端空间保持时长(秒)，0为不保持
	ClientPoolKeepAlive time.Duration `json:"clientPoolKeepAlive"`
	// 连接失败节点重试时间间隔(秒)
	ConnectFailRetry time.Duration `json:"connectFailRetry"`
	// 与服务器连接保持的心跳频率
	Heartbeat time.Duration `json:"heartbeat"`
	// 允许失丢的最大心跳数
	HeartbeatLoss int `json:"heartbeatLoss"`
	// 是否启用路由功能
	EnableRouter bool `json:"enableRouter"`
	// 路由客户端id
	RouterId string `json:"routerId"`
	// 路由服务器地址及端口号
	RouterSvr string `json:"routerSvr"`
	// 路由客户端参数
	Router *netstream.TStmClientOptions `json:"router"`
}

var defaultOptions = NewOptions(":8000")

func NewOptions(httpAddr string) *TOptions {
	//appDir := xfile.MainPkgPath()
	inst := &TOptions{
		HttpAddr:            httpAddr,
		ThriftAddr:          ":9000",
		PDLPath:             "",
		HostPath:            "",
		BaseRouter:          "",
		WatchPDLChanged:     true,
		WatchHostChanged:    true,
		PDLExt:              ".pdl",
		HostExt:             ".host",
		WriteTimeout:        1 * time.Minute,
		ReadTimeout:         1 * time.Minute,
		ConnectTimeout:      2 * time.Minute,
		ClientPoolKeepAlive: 2 * time.Minute,
		ConnectFailRetry:    2 * time.Minute,
		Router:              netstream.NewStmClientOptions(),
	}
	return inst
}

func NewOptionsFromMap(mp map[string]interface{}) (*TOptions, error) {
	opts := NewOptions(":8000")
	if s, ok := mp["httpAddr"].(string); ok {
		opts.HttpAddr = s
	}
	if s, ok := mp["thriftAddr"].(string); ok {
		opts.ThriftAddr = s
	}
	if s, ok := mp["pdlPath"].(string); ok {
		if err := opts.SetPdlPath(s); err != nil {
			return nil, err
		}
	}
	if s, ok := mp["hostPath"].(string); ok {
		if err := opts.SetHostPath(s); err != nil {
			return nil, err
		}
	}
	if s, ok := mp["baseRouter"].(string); ok {
		opts.BaseRouter = s
	}
	if s, ok := mp["baseNamespace"].(string); ok {
		opts.BaseNamespace = s
	}
	if s, ok := mp["watchHostChanged"]; ok {
		opts.WatchHostChanged = t.Bool(s)
	}
	if s, ok := mp["watchPDLChanged"]; ok {
		opts.WatchPDLChanged = t.Bool(s)
	}
	if s, ok := mp["pdlExt"].(string); ok {
		opts.PDLExt = s
	}
	if s, ok := mp["hostExt"].(string); ok {
		opts.HostExt = s
	}
	if s, ok := mp["writeTimeout"]; ok {
		opts.WriteTimeout = time.Duration(t.Int64(s)) * time.Second
	}
	if s, ok := mp["readTimeout"]; ok {
		opts.ReadTimeout = time.Duration(t.Int64(s)) * time.Second
	}
	if v, ok := mp["connectTimeout"]; ok {
		opts.ConnectTimeout = time.Duration(t.Int64(v)) * time.Second
	}
	if v, ok := mp["clientPoolKeepAlive"]; ok {
		opts.ClientPoolKeepAlive = time.Duration(t.Int64(v)) * time.Second
	}
	if v, ok := mp["connectFailRetry"]; ok {
		opts.ConnectFailRetry = time.Duration(t.Int64(v)) * time.Second
	}
	if v, ok := mp["heartbeat"]; ok {
		opts.Heartbeat = time.Duration(t.Int64(v)) * time.Second
	}
	if v, ok := mp["heartbeatLoss"]; ok {
		opts.HeartbeatLoss = t.Int(v)
	}
	if v, ok := mp["enableRouter"]; ok {
		opts.EnableRouter = t.Bool(v)
	}
	if v, ok := mp["routerId"]; ok {
		opts.RouterId = t.String(v)
	}
	if v, ok := mp["routerSvr"]; ok {
		opts.RouterSvr = t.String(v)
	}
	if opts.EnableRouter && opts.RouterId == "" {
		return nil, errors.New("未设置路由客户端id")
	}
	if opts.EnableRouter && opts.RouterSvr == "" {
		return nil, errors.New("未设置路由服务器地址")
	}
	if routerMp, ok := mp["router"].(map[string]interface{}); ok {
		opts.Router.LoadFromMap(routerMp)
	}
	return opts, nil
}

func (p *TOptions) SetHostPath(path string) error {
	if !xfile.Exists(path) {
		if err := xfile.Mkdir(path); err != nil {
			return err
		}
	}
	realPath := xfile.RealPath(path)
	p.HostPath = realPath
	return nil
}

func (p *TOptions) SetPdlPath(path string) error {
	if !xfile.Exists(path) {
		if err := xfile.Mkdir(path); err != nil {
			return err
		}
	}
	realPath := xfile.RealPath(path)
	p.PDLPath = realPath
	return nil
}

func (p *TOptions) SetThriftAddr(addr string) *TOptions {
	p.ThriftAddr = addr
	return p
}

func (p *TOptions) SetBaseRouter(baseRouter string) *TOptions {
	p.BaseRouter = baseRouter
	return p
}

func (p *TOptions) SetWatchHostChanged(enable bool) *TOptions {
	p.WatchHostChanged = enable
	return p
}

func (p *TOptions) SetWatchPDLChanged(enable bool) *TOptions {
	p.WatchPDLChanged = enable
	return p
}
