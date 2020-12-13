/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-28 12:04
* Description:
*****************************************************************/

package regcenter

import "time"

type THostStoreToken struct {
	// 服务协议项目名称
	Project    string `json:"project"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Ext        int    `json:"ext"`
	isSaved    bool
	weight     int
	isFail     bool
	failExpire time.Time
	// 节点所在的文件名
	fileId int
}

func (p *THostStoreToken) SetConnectFail(expire time.Time) {
	p.isFail = true
	p.failExpire = expire
}

func (p *THostStoreToken) SetConnectSuccess() {
	p.isFail = false
	p.failExpire = time.Now()
}

func (p *THostStoreToken) IncWeight(step int) {
	p.weight += step
}

func (p *THostStoreToken) Weight() int {
	return p.weight
}

func (p *THostStoreToken) IsFail() bool {
	return p.isFail
}

func (p *THostStoreToken) CanUse(now time.Time) bool {
	if !p.isFail || (p.isFail && now.After(p.failExpire)) {
		return true
	}
	return false
}
