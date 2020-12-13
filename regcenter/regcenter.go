/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-28 10:37
* Description:
*****************************************************************/

package regcenter

import (
	"github.com/go-xe2/xthrift/pdl"
)

type TRegCenter struct {
	hostStore *THostStore
	pdlStore  *TPDLStore
}

func NewRegCenter(hostPath, hostExt, pdlPath, pdlExt string, hostWatch, pdlWatch bool) *TRegCenter {
	inst := &TRegCenter{}
	inst.hostStore = NewHostStore(hostPath, hostExt, hostWatch)
	inst.pdlStore = NewPDLStore(pdlPath, pdlExt, pdlWatch)
	return inst
}

func (p *TRegCenter) Close() {
	if p.hostStore.IsFileWatch() {
		p.hostStore.DisableFileWatch()
	}
	if p.pdlStore.IsFileWatch() {
		p.pdlStore.DisableFileWatch()
	}
}

func (p *TRegCenter) WatchHostChanged() error {
	return p.hostStore.EnableFileWatch()
}

func (p *TRegCenter) WatchPDLChanged() error {
	return p.pdlStore.EnableFileWatch()
}

func (p *TRegCenter) Load() error {
	if err := p.hostStore.Load(); err != nil {
		return err
	}
	if err := p.pdlStore.Load(); err != nil {
		return err
	}
	return nil
}

func (p *TRegCenter) PdlStore() PDLStore {
	return p.pdlStore
}

func (p *TRegCenter) HostStore() HostStore {
	return p.hostStore
}

func (p *TRegCenter) PDLQuery() pdl.PDLQuery {
	return p.pdlStore
}
