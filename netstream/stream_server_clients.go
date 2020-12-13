/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-07 12:01
* Description:
*****************************************************************/

package netstream

func (p *TStreamServer) addClient(client StreamConn) string {
	p.clients.Set(client.Id(), client)
	return client.Id()
}

func (p *TStreamServer) getClient(id string) StreamConn {
	if cli := p.clients.Get(id); cli != nil {
		return cli.(StreamConn)
	}
	return nil
}

func (p *TStreamServer) removeClient(id string) {
	p.clients.Remove(id)
}

func (p *TStreamServer) AllClients() []string {
	return p.clients.Keys()
}

func (p *TStreamServer) GetClient(connId string) StreamConn {
	return p.getClient(connId)
}
