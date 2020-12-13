/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-03 15:50
* Description:
*****************************************************************/

package netstream

import "sync"

type TWaitGroup struct {
	delta int
	wg    sync.WaitGroup
}

func (p *TWaitGroup) Add(delta int) {
	p.delta += delta
	p.wg.Add(delta)
}

func (p *TWaitGroup) Done() {
	p.delta -= 1
	p.wg.Done()
}

func (p *TWaitGroup) DoneAll() {
	for {
		if p.delta <= 0 {
			return
		}
		p.Done()
	}
}

func (p *TWaitGroup) Wait() {
	p.wg.Wait()
}
