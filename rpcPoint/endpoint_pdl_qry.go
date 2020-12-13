/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-25 17:10
* Description:
*****************************************************************/

package rpcPoint

import "github.com/go-xe2/xthrift/pdl"

func (p *TEndPointServer) GetPdlQuery() pdl.PDLQuery {
	return p.regCenter.PDLQuery()
}
