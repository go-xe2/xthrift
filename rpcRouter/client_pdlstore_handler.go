package rpcRouter

import (
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/pdl"
)

func (p *TRouterClient) Install(pdl *pdl.FileProject, md5 string) {
	if pdl == nil {
		return
	}
	if err := p.sendRegProject(pdl, md5); err != nil {
		xlog.Error(err)
	}
	//pktId := time.Now().UnixNano()
	//
	//data := makeRegData(context.Background(), p.clientId, pktId, pdl, md5)
	//if err := p.send(data); err != nil {
	//	xlog.Error(err)
	//}
}

func (p *TRouterClient) UnInstall(pdl *pdl.FileProject) {

}
