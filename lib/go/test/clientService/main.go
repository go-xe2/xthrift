/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-27 15:27
* Description:
*****************************************************************/

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/xthrift/lib/go/test/comm"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
	"runtime/debug"
)

var processorFactory thrift.TProcessorFactory
var transportFactory = thrift.NewTTransportFactory()
var outputTransportFactory = thrift.NewTTransportFactory()

var inputProtocolFactory = xthrift.NewBinaryProtocolExFactory()
var outputProtocolFactory = xthrift.NewBinaryProtocolExFactory()
var forwardHeaders []string
var closed = make(chan byte, 1)
var pipo thrift.TTransport

func processRequests(in, out thrift.TTransport) error {
	processor := processorFactory.GetProcessor(out)
	inputTransport, err := transportFactory.GetTransport(in)
	if err != nil {
		return err
	}
	inputProtocol := inputProtocolFactory.GetProtocol(inputTransport)
	var outputTransport thrift.TTransport
	var outputProtocol thrift.TProtocol

	// for THeaderProtocol, we must use the same protocol instance for
	// input and output so that the response is in the same dialect that
	// the server detected the request was in.
	headerProtocol, ok := inputProtocol.(*thrift.THeaderProtocol)
	if ok {
		outputProtocol = inputProtocol
	} else {
		oTrans, err := outputTransportFactory.GetTransport(out)
		if err != nil {
			return err
		}
		outputTransport = oTrans
		outputProtocol = outputProtocolFactory.GetProtocol(outputTransport)
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
		if !in.IsOpen() {
			return errors.New("客户端被关闭")
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
			ctx = thrift.AddReadTHeaderToContext(ctx, headerProtocol.GetReadHeaders())
			ctx = thrift.SetWriteHeaderList(ctx, forwardHeaders)
		}
		fmt.Println("before processor process")
		ok, err := processor.Process(ctx, inputProtocol, outputProtocol)
		fmt.Println("after processor process ok:", ok, ", err:", err)
		if err, ok := err.(thrift.TTransportException); ok && err.TypeId() == thrift.END_OF_FILE {
			return err
		} else if err != nil {
			return err
		}
		if err, ok := err.(thrift.TApplicationException); ok && err.TypeId() == thrift.UNKNOWN_METHOD {
			continue
		}
		if !ok {
			break
		}
	}
	return nil
}

func innerAccept() (int32, error) {
	//closed := atomic.LoadInt32(&closed)
	//if closed != 0 {
	//	return closed, nil
	//}
	if !pipo.IsOpen() {
		xlog.Info("未连接服务器")
		return 0, nil
	}
	//var buf = make([]byte, thrift.DEFAULT_MAX_LENGTH)
	//var readBuf = make([]byte, thrift.DEFAULT_MAX_LENGTH)
	//stmTrans := thrift.NewStreamTransportR(bytes.NewBuffer(readBuf))
	//inTrans := thrift.NewTFramedTransport(stmTrans)

	//frameClient := thrift.NewTFramedTransport(pipo)
	//for {
	if err := processRequests(pipo, pipo); err != nil {
		fmt.Println("processRequests error:", err)
		close(closed)
	}
	//}

	return 0, nil
}

func AcceptLoop() error {
	go func() {
		innerAccept()
	}()
	select {
	case <-closed:
	}
	return nil
}

func main() {
	socket, err := thrift.NewTSocket("127.0.0.1:8001")
	if err != nil {
		panic(err)
	}
	namespaces := xthrift.NamespaceProcessor()
	handler := comm.NewSayHelloService()
	processor1 := comm.NewSayHelloServiceProcessor(handler)

	_ = namespaces.RegisterNamespace("demo.hello", processor1)

	processorFactory = xthrift.NewNamespaceProcessorFactory(namespaces)

	if err = socket.Open(); err != nil {
		panic(err)
	}
	pipo = socket
	fmt.Println("客户端服务连接成功，等待访问")
	if e := AcceptLoop(); e != nil {
		xlog.Error(err)
	}
	fmt.Println("关闭客户端.")
}
