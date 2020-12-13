/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-24 11:19
* Description:
*****************************************************************/

package rpcPoint

import (
	"errors"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
)

func (p *TEndPointServer) InitInnerServer() error {
	var err error
	p.thriftSvc, err = xthrift.NewServerFromListener(p.thriftLst)
	if err != nil {
		return err
	}
	p.transFac = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	p.inProtoFac = xthrift.NewBinaryProtocolExFactory()
	p.outProtoFac = xthrift.NewBinaryProtocolExFactory()
	p.thriftSvc.SetInputTransportFac(p.transFac)
	//outTransFac := thrift.NewTTransportFactory()
	// 此处不用frameTransport，数据已经在process中使用frame打包处理
	p.thriftSvc.SetOutputTransportFac(thrift.NewTTransportFactory())
	p.thriftSvc.SetInputProtocolFac(p.inProtoFac)
	p.thriftSvc.SetOutputProtocolFac(p.outProtoFac)
	p.thriftSvc.SetProcessorFac(newInnerProcessorFactory(p))
	return nil
}

func (p *TEndPointServer) InnerServerStart() error {
	if p.thriftSvc == nil {
		return errors.New("未初始化协议内部服务")
	}
	return p.thriftSvc.Serve()
}

func (p *TEndPointServer) InnerServerStop() error {
	if p.thriftSvc != nil {
		return p.thriftSvc.Stop()
	}
	return nil
}
