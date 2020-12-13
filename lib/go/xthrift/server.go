/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-21 12:10
* Description:
*****************************************************************/

package xthrift

import (
	"context"
	. "github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/core/exception"
	"github.com/go-xe2/x/os/xlog"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

type TXServer struct {
	closed int32
	wg     sync.WaitGroup
	mu     sync.Mutex

	processorFactory       TProcessorFactory
	serverTransport        TServerTransport
	inputTransportFactory  TTransportFactory
	outputTransportFactory TTransportFactory
	inputProtocolFactory   TProtocolFactory
	outputProtocolFactory  TProtocolFactory

	forwardHeaders []string
}

var _ TServer = (*TXServer)(nil)

func NewServer(addr string) (*TXServer, error) {
	trans, err := NewStandardServerSocket(addr)
	if err != nil {
		return nil, err
	}
	return &TXServer{
		processorFactory:       NewNamespaceProcessorFactoryDefault(),
		serverTransport:        trans,
		inputTransportFactory:  NewTFramedTransportFactory(NewTTransportFactory()),
		outputTransportFactory: NewTFramedTransportFactory(NewTTransportFactory()),
		inputProtocolFactory:   NewBinaryProtocolExFactory(),
		outputProtocolFactory:  NewBinaryProtocolExFactory(),
	}, nil
}

func NewServerFromListener(listener net.Listener) (*TXServer, error) {
	trans, err := NewStandardSocketFromListener(listener)
	if err != nil {
		return nil, err
	}
	return &TXServer{
		processorFactory:       NewNamespaceProcessorFactoryDefault(),
		serverTransport:        trans,
		inputTransportFactory:  NewTFramedTransportFactory(NewTTransportFactory()),
		outputTransportFactory: NewTFramedTransportFactory(NewTTransportFactory()),
		inputProtocolFactory:   NewBinaryProtocolExFactory(),
		outputProtocolFactory:  NewBinaryProtocolExFactory(),
	}, nil
}

func NewServerByOptions(trans TServerTransport, processorFac TProcessorFactory, inTransFac, outTransFac TTransportFactory, inFac, outFac TProtocolFactory) *TXServer {
	return &TXServer{
		processorFactory:       processorFac,
		serverTransport:        trans,
		inputTransportFactory:  inTransFac,
		outputTransportFactory: outTransFac,
		inputProtocolFactory:   inFac,
		outputProtocolFactory:  outFac,
	}
}

func (p *TXServer) SetServer(server TServerTransport) *TXServer {
	p.serverTransport = server
	return p
}

func (p *TXServer) SetProcessorFac(processorFac TProcessorFactory) *TXServer {
	p.processorFactory = processorFac
	return p
}

func (p *TXServer) SetInputTransportFac(transFac TTransportFactory) *TXServer {
	p.inputTransportFactory = transFac
	return p
}

func (p *TXServer) SetOutputTransportFac(transFac TTransportFactory) *TXServer {
	p.outputTransportFactory = transFac
	return p
}

func (p *TXServer) SetInputProtocolFac(protoFac TProtocolFactory) *TXServer {
	p.inputProtocolFactory = protoFac
	return p
}

func (p *TXServer) SetOutputProtocolFac(protoFac TProtocolFactory) *TXServer {
	p.outputProtocolFactory = protoFac
	return p
}

func (p *TXServer) ProcessorFactory() TProcessorFactory {
	return p.processorFactory
}

func (p *TXServer) ServerTransport() TServerTransport {
	return p.serverTransport
}

func (p *TXServer) InputTransportFactory() TTransportFactory {
	return p.inputTransportFactory
}

func (p *TXServer) OutputTransportFactory() TTransportFactory {
	return p.outputTransportFactory
}

func (p *TXServer) InputProtocolFactory() TProtocolFactory {
	return p.inputProtocolFactory
}

func (p *TXServer) OutputProtocolFactory() TProtocolFactory {
	return p.outputProtocolFactory
}

func (p *TXServer) Listen() error {
	if p.serverTransport == nil {
		panic(exception.NewText("未绑定基础服务"))
	}
	return p.serverTransport.Listen()
}

func (p *TXServer) Addr() net.Addr {
	if p.serverTransport == nil {
		return nil
	}
	if svc, ok := p.serverTransport.(*TStandardServerSocket); ok {
		return svc.Addr()
	} else if svc, ok := p.serverTransport.(*TServerSocket); ok {
		return svc.Addr()
	}
	return nil
}

func (p *TXServer) SetForwardHeaders(headers []string) {
	size := len(headers)
	if size == 0 {
		p.forwardHeaders = nil
		return
	}

	keys := make([]string, size)
	copy(keys, headers)
	p.forwardHeaders = keys
}

func (p *TXServer) innerAccept() (int32, error) {
	client, err := p.serverTransport.Accept()
	p.mu.Lock()
	defer p.mu.Unlock()
	closed := atomic.LoadInt32(&p.closed)
	if closed != 0 {
		return closed, nil
	}
	if err != nil {
		return 0, err
	}
	if client != nil {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			cTrans := NewClientTransport(client)
			if err := p.processRequests(cTrans); err != nil {
				xlog.Error("error processing request:", err)
			}
		}()
	}
	return 0, nil
}

func (p *TXServer) AcceptLoop() error {
	for {
		closed, err := p.innerAccept()
		if err != nil {
			return err
		}
		if closed != 0 {
			return nil
		}
	}
}

func (p *TXServer) Serve() error {
	err := p.Listen()
	if err != nil {
		return err
	}
	//addr := p.Addr()
	//addrName := ""
	//if addr != nil {
	//	addrName = addr.String()
	//}
	//xlog.Info("服务已运行，监听端口:", addrName)
	_ = p.AcceptLoop()
	return nil
}

func (p *TXServer) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if atomic.LoadInt32(&p.closed) != 0 {
		return nil
	}
	atomic.StoreInt32(&p.closed, 1)
	_ = p.serverTransport.Interrupt()
	p.wg.Wait()
	return nil
}

func (p *TXServer) processRequests(client TTransport) error {
	processor := p.processorFactory.GetProcessor(client)
	inputTransport, err := p.inputTransportFactory.GetTransport(client)
	if err != nil {
		return err
	}
	inputProtocol := p.inputProtocolFactory.GetProtocol(inputTransport)
	var outputTransport TTransport
	var outputProtocol TProtocol

	// for THeaderProtocol, we must use the same protocol instance for
	// input and output so that the response is in the same dialect that
	// the server detected the request was in.
	headerProtocol, ok := inputProtocol.(*THeaderProtocol)
	if ok {
		outputProtocol = inputProtocol
	} else {
		oTrans, err := p.outputTransportFactory.GetTransport(client)
		if err != nil {
			return err
		}
		outputTransport = oTrans
		outputProtocol = p.outputProtocolFactory.GetProtocol(outputTransport)
	}

	defer func() {
		if e := recover(); e != nil {
			xlog.Error("panic in processor: %s: %s", e, debug.Stack())
		}
	}()

	if inputTransport != nil {
		defer inputTransport.Close()
	}
	if outputTransport != nil {
		defer outputTransport.Close()
	}
	for {
		if atomic.LoadInt32(&p.closed) != 0 {
			return nil
		}
		ctx := context.Background()
		if headerProtocol != nil {
			// We need to call ReadFrame here, otherwise we won't
			// get any headers on the AddReadTHeaderToContext call.
			//
			// ReadFrame is safe to be called multiple times so it
			// won't break when it's called again later when we
			// actually start to read the message.
			if err := headerProtocol.ReadFrame(); err != nil {
				return err
			}
			ctx = AddReadTHeaderToContext(ctx, headerProtocol.GetReadHeaders())
			ctx = SetWriteHeaderList(ctx, p.forwardHeaders)
		}
		ok, err := processor.Process(ctx, inputProtocol, outputProtocol)
		if err, ok := err.(TTransportException); ok && err.TypeId() == END_OF_FILE {
			return nil
		} else if err != nil {
			return err
		}
		if err, ok := err.(TApplicationException); ok && err.TypeId() == UNKNOWN_METHOD {
			continue
		}
		if !ok {
			break
		}
	}
	return nil
}
