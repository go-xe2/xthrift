package rpcRouter

import (
	"context"
	"errors"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/x/utils/xrand"
	"github.com/go-xe2/xthrift/netstream"
	"github.com/go-xe2/xthrift/pdl"
	"github.com/go-xe2/xthrift/regcenter"
)

func (p *TRouterServer) processRecv(conn netstream.StreamConn, pktId int64, proto RouterProto, pktData []byte, pkt TPacketType) error {
	switch pkt {
	case REG_PACKET:
		// 注册服务
		return p.regProcess(conn, proto, pktId)
	case CALL_PACKET:
		// 调用服务
		return p.callProcess(conn, proto, pktData, pktId)
	case REPLY_PACKET:
		// 调用回复
		return p.replyProcess(conn, pktId, pktData)
	}
	return nil
}

func (p *TRouterServer) regProcess(conn netstream.StreamConn, proto RouterProto, pktId int64) error {
	clientId, project, md5, err := proto.ReadRegBegin()
	if err != nil {
		return err
	}
	p.saveClient(clientId, conn.Id())
	if p.HasProject(project, md5) {
		// 协议项目已经存在
		proj := p.GetProject(project)
		xlog.Info("协议项目:", project, "已经存在，md5:", md5)
		p.center.HostStore().AddHostWithProject(proj, clientId, 0, 1)
		if err := p.center.HostStore().Save(); err != nil {
			return err
		}
		resData := makeRegResData(context.Background(), pktId, project, md5)
		if _, err := conn.Send(resData); err != nil {
			return err
		}
		return nil
	}
	pdlData, err := proto.ReadData()
	if err != nil {
		return err
	}
	if err := proto.ReadRegEnd(); err != nil {
		return err
	}
	var proj *pdl.FileProject
	projInfo := p.center.PdlStore().GetProjectByName(project)
	if projInfo == nil {
		proj, err = p.center.PdlStore().AddProjectFromContent(pdlData)
		if err != nil {
			xlog.Debug("AddProjectFromContent err:", err)
			return err
		}
	} else {
		// 已经存在，卸载之前协议
		if err := p.center.PdlStore().RemoveProject(project); err != nil {
			xlog.Debug("removeProject error:", err)
		}
		proj, err = p.center.PdlStore().AddProjectFromContent(pdlData)
		if err != nil {
			xlog.Debug("AddProjectFromContent err:", err)
			return err
		}
	}
	p.center.HostStore().AddHostWithProject(proj, clientId, 0, 1)
	if err := p.center.HostStore().Save(); err != nil {
		return err
	}
	p.AddProject(project, md5, proj)
	resData := makeRegResData(context.Background(), pktId, project, md5)
	if _, err := conn.Send(resData); err != nil {
		return err
	}
	return nil
}

func (p *TRouterServer) callProcess(conn netstream.StreamConn, proto RouterProto, pktData []byte, pktId int64) error {
	namespace, _, _, err := proto.ReadCallBegin()
	if err != nil {
		return err
	}
	defer proto.ReadCallEnd()

	hosts := p.center.HostStore().GetSvcHosts(namespace)
	onLines := make([]*regcenter.THostStoreToken, 0)
	for _, h := range hosts {
		if h.Ext > 0 {
			onLines = append(onLines, h)
		}
	}
	if len(onLines) == 0 {
		return errors.New("没有可用的服务资源")
	}
	if err := proto.ReadCallEnd(); err != nil {
		return err
	}
	return p.transportCall(conn, onLines, pktId, pktData)
}

func (p *TRouterServer) replyProcess(fromConn netstream.StreamConn, pktId int64, pktData []byte) error {
	// 数据调用返回
	fromId := p.getClientIdByConnId(fromConn.Id())
	return p.callReply(fromId, pktId, pktData)
	// 数据调用返回
}

// 服务内部数据转换使用
func (p *TRouterServer) transportCall(fromConn netstream.StreamConn, hosts []*regcenter.THostStoreToken, pktId int64, pktData []byte) error {
	// 随机获取一条链接
	idx := 0
	onLineCount := len(hosts)
	if onLineCount > 0 {
		idx = xrand.N(0, onLineCount-1)
	}
	host := hosts[idx]
	clientId := host.Host
	fromClientId := p.getClientIdByConnId(fromConn.Id())
	targetConn := p.getClientConn(clientId)
	if targetConn == nil {
		return errors.New("服务已离线或没有可用的服务资源")
	}
	sender := newServerTransSender(p, fromClientId)
	callTargetId := p.getClientIdByConnId(targetConn.Id())
	p.callTimeout(callTargetId, pktId, sender)
	if _, err := targetConn.Send(pktData); err != nil {
		return err
	}
	return nil
}
