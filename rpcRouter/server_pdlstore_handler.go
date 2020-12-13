package rpcRouter

import (
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/pdl"
)

func (p *TRouterServer) Install(pdl *pdl.FileProject, md5 string) {
	if pdl == nil {
		return
	}
	xlog.Debug("rpc-router install project:", pdl.GetProjectName(), ", md5:", md5)
	p.AddProject(pdl.GetProjectName(), md5, pdl)
}

func (p *TRouterServer) UnInstall(pdl *pdl.FileProject) {
	if pdl == nil {
		delete(p.pdlProjects, pdl.GetProjectName())
	}
}
