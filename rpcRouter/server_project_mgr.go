package rpcRouter

import "github.com/go-xe2/xthrift/pdl"

func (p *TRouterServer) HasProject(projectName string, md5 string) bool {
	if p, ok := p.pdlProjects[projectName]; ok {
		if p.md5 == md5 {
			return true
		}
	}
	return false
}

func (p *TRouterServer) AddProject(projectName string, md5 string, proj *pdl.FileProject) {
	p.pdlProjects[projectName] = newProjectInfo(md5, proj)
}

func (p *TRouterServer) GetProject(projectName string) *pdl.FileProject {
	if p, ok := p.pdlProjects[projectName]; ok {
		return p.project
	}
	return nil
}
