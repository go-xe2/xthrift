/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-09-01 16:37
* Description:
*****************************************************************/

package rpcRouter

import (
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/type/t"
	"time"
)

type TSvrOptions struct {
	// http监听地址
	HttpAddr string `json:"httpAddr"`
	// 内部服务调用地址 `json
	RouterAddr string `json:"routerAddr"`
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
	// 允许丢失的最大心跳包数
	AllowMaxLoss int `json:"allowMaxLoss"`
	// 数据接收队列大小
	RecvBufSize int `json:"recvBufSize"`
	// 数据发送队列大小
	SendBufSize int `json:"sendBufSize"`
}

var defaultSvrOptions = NewSvcOptions()

func NewSvcOptions() *TSvrOptions {
	//appDir := xfile.MainPkgPath()
	inst := &TSvrOptions{
		HttpAddr:            ":3003",
		RouterAddr:          ":3004",
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
		Heartbeat:           5 * time.Minute,
		AllowMaxLoss:        3,
		RecvBufSize:         512,
		SendBufSize:         512,
	}
	return inst
}

func NewSvrOptionsFromMap(mp map[string]interface{}) (*TSvrOptions, error) {
	opts := NewSvcOptions()
	if s, ok := mp["httpAddr"].(string); ok {
		opts.HttpAddr = s
	}
	if s, ok := mp["routerAddr"].(string); ok {
		opts.RouterAddr = s
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
	if v, ok := mp["allowMaxLoss"]; ok {
		opts.AllowMaxLoss = t.Int(v)
	}
	if v, ok := mp["recvBufSize"]; ok {
		opts.RecvBufSize = t.Int(v)
	}
	if v, ok := mp["sendBufSize"]; ok {
		opts.SendBufSize = t.Int(v)
	}
	return opts, nil
}

func (p *TSvrOptions) SetHostPath(path string) error {
	if !xfile.Exists(path) {
		if err := xfile.Mkdir(path); err != nil {
			return err
		}
	}
	realPath := xfile.RealPath(path)
	p.HostPath = realPath
	return nil
}

func (p *TSvrOptions) SetPdlPath(path string) error {
	if !xfile.Exists(path) {
		if err := xfile.Mkdir(path); err != nil {
			return err
		}
	}
	realPath := xfile.RealPath(path)
	p.PDLPath = realPath
	return nil
}

func (p *TSvrOptions) SetRouterAddr(addr string) *TSvrOptions {
	p.RouterAddr = addr
	return p
}

func (p *TSvrOptions) SetBaseRouter(baseRouter string) *TSvrOptions {
	p.BaseRouter = baseRouter
	return p
}

func (p *TSvrOptions) SetWatchHostChanged(enable bool) *TSvrOptions {
	p.WatchHostChanged = enable
	return p
}

func (p *TSvrOptions) SetWatchPDLChanged(enable bool) *TSvrOptions {
	p.WatchPDLChanged = enable
	return p
}
