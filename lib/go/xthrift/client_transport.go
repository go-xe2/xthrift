/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-29 10:56
* Description:
*****************************************************************/

package xthrift

import (
	"github.com/apache/thrift/lib/go/thrift"
	"gopkg.in/mgo.v2/bson"
)

type ClientTransport interface {
	thrift.TTransport
	// 客户端id
	ID() string
	RemoteAddr() string
}

type tClientTransport struct {
	id string
	thrift.TTransport
}

func NewClientTransport(trans thrift.TTransport) ClientTransport {
	return &tClientTransport{
		TTransport: trans,
		id:         bson.NewObjectId().Hex(),
	}
}

func (p *tClientTransport) ID() string {
	return p.id
}

func (p *tClientTransport) RemoteAddr() string {
	if sck, ok := p.TTransport.(*thrift.TSocket); ok {
		return sck.Conn().RemoteAddr().String()
	}
	return ""
}
