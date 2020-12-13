/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-03 10:55
* Description:
*****************************************************************/

package netstream

import (
	"bytes"
	"errors"
	"github.com/go-xe2/x/core/exception"
	"github.com/go-xe2/x/type/xbinary"
	"math"
)

// 数据包标识
const packFlag int16 = 0x1A
const PacketHeadSize = 11

func packeData(data []byte, cmd ...StreamCmdType) []byte {
	c := StreamCmdSend
	if len(cmd) > 0 {
		c = cmd[0]
	}
	size := len(data)
	result := bytes.NewBuffer([]byte{})
	result.Write(xbinary.BeEncodeInt16(packFlag))
	result.Write(xbinary.BeEncodeInt8(int8(c)))
	result.Write(xbinary.BeEncodeInt64(int64(size)))
	result.Write(data)
	return result.Bytes()
}

func unpackeData(test []byte) (isReady bool, data []byte, cmd StreamCmdType, remainder []byte) {
	size := len(test)
	if size < PacketHeadSize {
		return false, nil, StreamCmdUnknown, test
	}
	flag := xbinary.BeDecodeToInt16(test[0:2])
	if flag != packFlag {
		return false, nil, StreamCmdUnknown, test
	}
	c := xbinary.DecodeToInt8(test[2:3])
	cmd = StreamCmdType(c)
	if cmd == StreamCmdUnknown || cmd >= StreamCmdInvalid {
		return false, nil, StreamCmdUnknown, test
	}
	dataSize := xbinary.BeDecodeToInt64(test[3:PacketHeadSize])
	if dataSize > int64(size-PacketHeadSize) {
		// 数据接收未完成
		return false, nil, StreamCmdUnknown, test
	}
	data = test[PacketHeadSize : dataSize+PacketHeadSize]
	remainder = test[dataSize+PacketHeadSize:]
	return true, data, cmd, remainder
}

func packeCallData(seqId int64, data []byte) []byte {
	result := bytes.NewBuffer([]byte{})
	result.Write(xbinary.BeEncodeInt64(seqId))
	result.Write(data)
	return result.Bytes()
}

func unpackeCallData(in []byte) (seqId int64, data []byte) {
	if len(in) < 8 {
		return 0, nil
	}
	seqId = xbinary.BeDecodeToInt64(in[0:8])
	data = in[8:]
	return
}

// 打包呼叫返回数据
func packeCallReplyData(seqId int64, data []byte, err error) []byte {
	result := bytes.NewBuffer([]byte{})
	var errLen int16 = 0
	var errText = ""
	if err != nil {
		errText = err.Error()
		if len(errText) > math.MaxInt16 {
			errText = errText[:math.MaxInt16]
		}
		errLen = int16(len(errText))
	}
	result.Write(xbinary.BeEncodeInt64(seqId))
	result.Write(xbinary.BeEncodeInt16(errLen))
	if errLen > 0 {
		result.Write(xbinary.BeEncodeString(errText))
	}
	result.Write(data)
	return result.Bytes()
}

// 解包呼叫返回数据
func unpackCallReplyData(in []byte) (seqId int64, data []byte, err error) {
	if len(in) < 10 {
		return 0, nil, exception.NewText("数据长度不足")
	}
	seqId = xbinary.BeDecodeToInt64(in[0:8])
	errLen := xbinary.BeDecodeToInt16(in[8:10])
	if errLen > 0 {
		s := xbinary.BeDecodeToString(in[10 : 10+errLen])
		err = errors.New(s)
	}
	data = in[10+errLen:]
	return
}

// 打包SendTo消息数据包
func packeSendToData(targetId string, data []byte) []byte {
	s := targetId
	if len(targetId) > math.MaxInt8 {
		s = s[:math.MaxInt8]
	}
	size := int8(len(s))
	result := bytes.NewBuffer([]byte{})
	result.Write(xbinary.BeEncodeInt8(size))
	if size > 0 {
		result.Write(xbinary.BeEncodeString(s))
	}
	result.Write(data)
	return result.Bytes()
}

// 解包SendTo消息数据包
func unpackeSendToData(in []byte) (targetId string, data []byte) {
	if len(in) < 1 {
		return "", nil
	}
	strLen := xbinary.BeDecodeToInt8(in[:1])
	if strLen > 0 {
		targetId = xbinary.BeDecodeToString(in[1 : 1+strLen])
	}
	data = in[1+strLen:]
	return
}

// 打包呼叫客户端数据
func packeCallToData(seqId int64, targetId string, data []byte) []byte {
	s := targetId
	if len(s) > math.MaxInt8 {
		s = s[0:math.MaxInt8]
	}
	sLen := int8(len(s))
	result := bytes.NewBuffer([]byte{})
	result.Write(xbinary.BeEncodeInt64(seqId))
	result.Write(xbinary.BeEncodeInt8(sLen))
	if sLen > 0 {
		result.Write(xbinary.BeEncodeString(s))
	}
	result.Write(data)
	return result.Bytes()
}

// 解包客户端呼叫数据
func unpackeCallToData(in []byte) (seqId int64, targetId string, data []byte) {
	if len(in) < 9 {
		return 0, "", nil
	}
	seqId = xbinary.BeDecodeToInt64(in[0:8])
	sLen := xbinary.BeDecodeToInt8(in[8:9])
	if sLen > 0 {
		targetId = xbinary.BeDecodeToString(in[9 : sLen+9])
	}
	data = in[sLen+9:]
	return
}

func packeClientId(id string) []byte {
	if len(id) > math.MaxInt16 {
		id = id[:math.MaxInt16]
	}
	size := len(id)
	result := bytes.NewBuffer([]byte{})
	result.Write(xbinary.BeEncodeInt16(int16(size)))
	result.Write(xbinary.BeEncodeString(id))
	return result.Bytes()
}

func unpackeClientId(in []byte) (id string) {
	if len(in) < 2 {
		return ""
	}
	size := xbinary.BeDecodeToInt16(in[:2])
	s := xbinary.BeDecodeToString(in[2 : size+2])
	return s
}

func packeRequestData(reqId string, namespace string, body []byte) []byte {
	result := bytes.NewBuffer([]byte{})
	if len(reqId) > math.MaxInt16 {
		reqId = reqId[:math.MaxInt16]
	}
	reqIdSize := int16(len(reqId))
	if len(namespace) > math.MaxInt32 {
		namespace = namespace[:math.MaxInt32]
	}
	nameSize := int32(len(namespace))
	result.Write(xbinary.BeEncodeInt16(reqIdSize))
	result.Write(xbinary.BeEncodeInt32(nameSize))
	result.Write(xbinary.BeEncodeString(reqId))
	result.Write(xbinary.BeEncodeString(namespace))
	if body != nil {
		result.Write(body)
	}
	return result.Bytes()
}

func unpackRequestData(in []byte) (reqId string, namespace string, body []byte) {
	if len(in) < 6 {
		return "", "", nil
	}
	reqIdSize := int(xbinary.BeDecodeToInt16(in[:2]))
	nameSize := int(xbinary.BeDecodeToInt32(in[2:6]))
	if 6+reqIdSize+nameSize > len(in) {
		return "", "", nil
	}
	reqId = xbinary.BeDecodeToString(in[6 : 6+reqIdSize])
	namespace = xbinary.BeDecodeToString(in[6+reqIdSize : 6+reqIdSize+nameSize])
	body = in[6+reqIdSize+nameSize:]
	return
}

func packeResponseData(reqId string, body []byte) []byte {
	result := bytes.NewBuffer([]byte{})
	if len(reqId) > math.MaxInt16 {
		reqId = reqId[:math.MaxInt16]
	}
	reqIdSize := int16(len(reqId))
	result.Write(xbinary.BeEncodeInt16(reqIdSize))
	result.Write(xbinary.BeEncodeString(reqId))
	if body != nil {
		result.Write(body)
	}
	return result.Bytes()
}

func unpackResponseData(in []byte) (reqId string, body []byte) {
	if len(in) < 2 {
		return "", nil
	}
	reqIdSize := int(xbinary.BeDecodeToInt16(in[:2]))
	if 2+reqIdSize > len(in) {
		return "", nil
	}
	reqId = xbinary.BeDecodeToString(in[2 : 2+reqIdSize])
	body = in[2+reqIdSize:]
	return
}
