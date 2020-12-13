/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-07-24 15:50
* Description:
*****************************************************************/

package xthrift

import (
	"github.com/apache/thrift/lib/go/thrift"
)

type matchRuleItem struct {
	name string
	typ  thrift.TType
	id   int16
}

type StructFieldReader interface {
	Match() func(string, int16, thrift.TType) bool
	Reader() func(fdName string, fdId int16, fdType thrift.TType, reader *TProtocolHelper)
}

type tStructFieldReader struct {
	rule   *matchRuleItem
	reader func(fdName string, fdId int16, fdType thrift.TType, reader *TProtocolHelper)
}

var _ StructFieldReader = (*tStructFieldReader)(nil)

func (p *tStructFieldReader) Match() func(string, int16, thrift.TType) bool {
	rule := p.rule
	return func(s string, i int16, tType thrift.TType) bool {
		if (i > 0 && rule.id > 0 && i == rule.id) || (i < 0 && s != "" && s == rule.name) {
			return true
		}
		return false
	}
}

func (p *tStructFieldReader) Reader() func(fdName string, fdId int16, fdType thrift.TType, reader *TProtocolHelper) {
	return p.reader
}

type TProtocolHelper struct {
	thrift.TProtocol
}

func NewProtocolHelper(proto thrift.TProtocol) *TProtocolHelper {
	if tmp, ok := proto.(*TProtocolHelper); ok {
		return tmp
	}
	return &TProtocolHelper{
		TProtocol: proto,
	}
}

func (p *TProtocolHelper) WriteStruct(name string, fieldsWriter func(writer *TProtocolHelper)) (err error) {
	if e := p.WriteStructBegin(name); e != nil {
		return e
	}
	defer func() {
		if e := recover(); e != nil {
			_ = p.WriteFieldStop()
			err = thrift.NewTProtocolException(err)
		}
	}()
	fieldsWriter(p)
	if e := p.WriteFieldStop(); e != nil {
		return e
	}
	if e := p.WriteStructEnd(); e != nil {
		return e
	}
	return nil
}

func (p *TProtocolHelper) WriteStructField(fdName string, fdId int16, structName string, structWriter func(writer *TProtocolHelper)) {
	p.WriteField(fdName, thrift.STRUCT, fdId, func(writer *TProtocolHelper) {
		if e := writer.WriteStruct(structName, structWriter); e != nil {
			panic(e)
		}
	})
}

func (p *TProtocolHelper) WriteField(fdName string, fdType thrift.TType, fdId int16, dataWriter func(writer *TProtocolHelper)) {
	if e := p.WriteFieldBegin(fdName, fdType, fdId); e != nil {
		panic(e)
	}
	dataWriter(p)
	if e := p.WriteFieldEnd(); e != nil {
		panic(e)
	}
}

func (p *TProtocolHelper) WriteStructBegin(structName string) error {
	return p.TProtocol.WriteStructBegin(structName)
}

func (p *TProtocolHelper) WriteStructEnd() error {
	return p.TProtocol.WriteStructEnd()
}

func (p *TProtocolHelper) WriteFieldStop() error {
	return p.TProtocol.WriteFieldStop()
}

func (p *TProtocolHelper) ReadStructBegin() (structName string, err error) {
	structName, err = p.TProtocol.ReadStructBegin()
	return
}

func (p *TProtocolHelper) ReadStructEnd() error {
	return p.TProtocol.ReadStructEnd()
}

func (p *TProtocolHelper) ReadFieldBegin() (name string, fdType thrift.TType, fdId int16, err error) {
	name, fdType, fdId, err = p.TProtocol.ReadFieldBegin()
	return
}

func (p *TProtocolHelper) ReadFieldEnd() error {
	return p.TProtocol.ReadFieldEnd()
}

func (p *TProtocolHelper) WriteStrField(fdName string, fdId int16, val string) {
	p.WriteField(fdName, thrift.STRING, fdId, func(writer *TProtocolHelper) {
		if e := writer.WriteString(val); e != nil {
			panic(e)
		}
	})
}

func (p *TProtocolHelper) WriteBoolField(fdName string, fdId int16, val bool) {
	p.WriteField(fdName, thrift.BOOL, fdId, func(writer *TProtocolHelper) {
		if e := writer.WriteBool(val); e != nil {
			panic(e)
		}
	})
}

func (p *TProtocolHelper) WriteByteField(fdName string, fdId int16, val byte) {
	p.WriteField(fdName, thrift.BYTE, fdId, func(writer *TProtocolHelper) {
		if e := writer.WriteByte(int8(val)); e != nil {
			panic(e)
		}
	})
}

func (p *TProtocolHelper) WriteInt16Field(fdName string, fdId int16, val int16) {
	p.WriteField(fdName, thrift.I16, fdId, func(writer *TProtocolHelper) {
		if e := writer.WriteI16(val); e != nil {
			panic(e)
		}
	})
}

func (p *TProtocolHelper) WriteInt32Field(fdName string, fdId int16, val int32) {
	p.WriteField(fdName, thrift.I32, fdId, func(writer *TProtocolHelper) {
		if e := writer.WriteI32(val); e != nil {
			panic(e)
		}
	})
}

func (p *TProtocolHelper) WriteInt64Field(fdName string, fdId int16, val int64) {
	p.WriteField(fdName, thrift.I64, fdId, func(writer *TProtocolHelper) {
		if e := writer.WriteI64(val); e != nil {
			panic(e)
		}
	})
}

func (p *TProtocolHelper) WriteDoubleField(fdName string, fdId int16, val float64) {
	p.WriteField(fdName, thrift.DOUBLE, fdId, func(writer *TProtocolHelper) {
		if e := writer.WriteDouble(val); e != nil {
			panic(e)
		}
	})
}

func (p *TProtocolHelper) WriteListField(fdName string, fdId int16, elemType thrift.TType, lstSize int, itemsWriter func(writer *TProtocolHelper)) {
	p.WriteField(fdName, thrift.LIST, fdId, func(writer *TProtocolHelper) {
		if e := writer.WriteListBegin(elemType, lstSize); e != nil {
			panic(e)
		}
		itemsWriter(writer)
		if e := writer.WriteListEnd(); e != nil {
			panic(e)
		}
	})
}

func (p *TProtocolHelper) WriteSetField(fdName string, fdId int16, stSize int, elemType thrift.TType, itemsWriter func(writer *TProtocolHelper)) {
	p.WriteField(fdName, thrift.SET, fdId, func(writer *TProtocolHelper) {
		if e := writer.WriteSetBegin(elemType, stSize); e != nil {
			panic(e)
		}
		itemsWriter(writer)
		if e := writer.WriteSetEnd(); e != nil {
			panic(e)
		}
	})
}

func (p *TProtocolHelper) WriteMapField(fdName string, fdId int16, kt, vt thrift.TType, size int, itemsWriter func(writer *TProtocolHelper)) {
	p.WriteField(fdName, thrift.MAP, fdId, func(writer *TProtocolHelper) {
		if e := writer.WriteMapBegin(kt, vt, size); e != nil {
			panic(e)
		}
		itemsWriter(writer)
		if e := writer.WriteMapEnd(); e != nil {
			panic(e)
		}
	})
}

func (p *TProtocolHelper) MakeFieldReader(fdName string, fdId int16, fdType thrift.TType, dataReader func(reader *TProtocolHelper)) StructFieldReader {
	result := &tStructFieldReader{
		rule: &matchRuleItem{
			name: fdName,
			typ:  fdType,
			id:   fdId,
		},
		reader: func(fdName string, fdId int16, fdType thrift.TType, reader *TProtocolHelper) {
			dataReader(reader)
		},
	}
	return result
}

func (p *TProtocolHelper) ReadStructFields(readers ...StructFieldReader) (structName string, err error) {
	if structName, err = p.ReadStructBegin(); err != nil {
		return
	}
	ruleCount := len(readers)
	tmpReaders := make([]StructFieldReader, ruleCount)
	for i, v := range readers {
		tmpReaders[i] = v
	}
	for {
		fdName, fdType, fdId, err := p.ReadFieldBegin()
		if err != nil {
			return structName, err
		}
		if fdType == thrift.STOP {
			break
		}
		if err := p.ReadFieldEnd(); err != nil {
			return structName, err
		}
		var isMatchFn = false
		for i := 0; i < ruleCount; i++ {
			reader := tmpReaders[i]
			if reader == nil {
				continue
			}
			matchFn := reader.Match()
			isMatchFn = matchFn(fdName, fdId, fdType)
			if isMatchFn {
				reader.Reader()(fdName, fdId, fdType, p)
				tmpReaders[i] = nil
				break
			}
		}
		if !isMatchFn {
			_ = p.Skip(fdType)
		}
	}
	if err = p.ReadStructEnd(); err != nil {
		return
	}
	return
}

func (p *TProtocolHelper) ReadStrData(f *string) {
	if s, e := p.ReadString(); e != nil {
		panic(e)
	} else {
		*f = s
	}
}

func (p *TProtocolHelper) ReadI16Data(n *int16) {
	if d, e := p.ReadI16(); e != nil {
		panic(e)
	} else {
		*n = d
	}
}

func (p *TProtocolHelper) ReadI32Data(n *int32) {
	if s, e := p.ReadI32(); e != nil {
		panic(e)
	} else {
		*n = s
	}
}

func (p *TProtocolHelper) ReadI64Data(n *int64) {
	if s, e := p.ReadI64(); e != nil {
		panic(e)
	} else {
		*n = s
	}
}

func (p *TProtocolHelper) ReadByteData(n *byte) {
	if s, e := p.ReadByte(); e != nil {
		panic(e)
	} else {
		*n = byte(s)
	}
}

func (p *TProtocolHelper) ReadBoolData(n *bool) {
	if s, e := p.ReadBool(); e != nil {
		panic(e)
	} else {
		*n = s
	}
}

func (p *TProtocolHelper) ReadDoubleData(n *float64) {
	if s, e := p.ReadDouble(); e != nil {
		panic(e)
	} else {
		*n = s
	}
}

func (p *TProtocolHelper) ReadListData(listMake func(size int), itemReader func(index int, elemType thrift.TType, reader *TProtocolHelper)) {
	elemType, size, err := p.ReadListBegin()
	if err != nil {
		panic(err)
	}
	listMake(size)
	for i := 0; i < size; i++ {
		itemReader(i, elemType, p)
	}
	if err := p.ReadListEnd(); err != nil {
		panic(err)
	}
}

func (p *TProtocolHelper) ReadSetData(itemReader func(index int, elemType thrift.TType, reader *TProtocolHelper)) {
	elemType, size, err := p.ReadSetBegin()
	if err != nil {
		panic(err)
	}
	for i := 0; i < size; i++ {
		itemReader(i, elemType, p)
	}
	if err := p.ReadSetEnd(); err != nil {
		panic(err)
	}
}

func (p *TProtocolHelper) ReadMapData(itemReader func(kt, vt thrift.TType, reader *TProtocolHelper)) {
	kt, vt, size, err := p.ReadMapBegin()
	if err != nil {
		panic(err)
	}
	for i := 0; i < size; i++ {
		itemReader(kt, vt, p)
	}
	if err := p.ReadMapEnd(); err != nil {
		panic(err)
	}
}

func (p *TProtocolHelper) ReadField(dataReader func(fdName string, fdId int16, fdType thrift.TType, reader *TProtocolHelper)) {
	fdName, fdType, fdId, err := p.ReadFieldBegin()
	if err != nil {
		panic(err)
	}
	dataReader(fdName, fdId, fdType, p)
	if err := p.ReadFieldEnd(); err != nil {
		panic(err)
	}
	return
}
