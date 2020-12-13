package rpcRouter

import (
	"github.com/go-xe2/xthrift/gateway"
	"github.com/go-xe2/xthrift/netstream"
	"github.com/go-xe2/xthrift/regcenter"
)

type TRouterClient struct {
	clientId string
	pdlStore regcenter.PDLStore
	client   netstream.StreamClient
	rpcCall  gateway.RpcService
}

func NewClient(clientId string, rpcCall gateway.RpcService, pdlStore regcenter.PDLStore, svrHost string, opts *netstream.TStmClientOptions) (*TRouterClient, error) {
	clt, err := netstream.NewStreamClient(svrHost, opts)
	if err != nil {
		return nil, err
	}
	inst := &TRouterClient{
		clientId: clientId,
		pdlStore: pdlStore,
		rpcCall:  rpcCall,
		client:   clt,
	}
	return inst, nil
}

func (p *TRouterClient) Open() error {
	if err := p.client.Open(); err != nil {
		return err
	}
	p.pdlStore.SetHandler(p)
	p.client.SetHandler(p)
	return nil
}

func (p *TRouterClient) Close() error {
	return p.client.Close()
}

func (p *TRouterClient) IsOpen() bool {
	return p.client.IsOpen()
}
