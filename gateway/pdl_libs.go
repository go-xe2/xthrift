/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-25 10:48
* Description:
*****************************************************************/

package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/x/type/xstring"
	"github.com/go-xe2/xthrift/pdl"
)

func writeStructFromMap(qry pdl.PDLQuery, stru *pdl.FileStruct, proto thrift.TProtocol, svc *pdl.FileService, mp map[string]interface{}) error {
	fields := stru.Fields
	if err := proto.WriteStructBegin(xstring.Camel2UnderScore(stru.Type.TypName, "_")); err != nil {
		return err
	}
	for _, fd := range fields {
		fdName := fd.Name
		paramName := xstring.Camel2UnderScore(fdName, "_")
		fdType, err := StripTypedef(qry, svc, fd.FieldType)
		if err != nil {
			return err
		}
		if fdType.Type == pdl.SPD_VOID || fdType.Type == pdl.SPD_UNKNOWN {
			continue
		}

		v, ok := mp[paramName]
		if !ok || v == nil {
			if fd.Limit == pdl.SPDLimitRequired {
				return fmt.Errorf("字段%s不能为空", paramName)
			}
			// 忽略该字段
			continue
		}
		if err := proto.WriteFieldBegin(paramName, fdType.Type.ThriftType(), fd.Id); err != nil {
			return err
		}
		switch fdType.Type {
		case pdl.SPD_STR:
			s := t.String(v)
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, s); err != nil {
				return err
			}
			break
		case pdl.SPD_BOOL:
			b := t.Bool(v)
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, b); err != nil {
				return err
			}
			break
		case pdl.SPD_I08:
			n8 := t.Int8(v)
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, n8); err != nil {
				return err
			}
			break
		case pdl.SPD_I16:
			n8 := t.Int16(v)
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, n8); err != nil {
				return err
			}
			break
		case pdl.SPD_I32:
			n8 := t.Int32(v)
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, n8); err != nil {
				return err
			}
			break
		case pdl.SPD_I64:
			n8 := t.Int64(v)
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, n8); err != nil {
				return err
			}
			break
		case pdl.SPD_DOUBLE:
			n8 := t.Float64(v)
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, n8); err != nil {
				return err
			}
			break
		case pdl.SPD_LIST:
			arr, ok := v.([]interface{})
			if !ok {
				if fd.Limit == pdl.SPDLimitRequired {
					return fmt.Errorf("字段%s为空或不是Array数据类型", paramName)
				}
				continue
			}
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, arr); err != nil {
				return err
			}
			break
		case pdl.SPD_SET:
			arr, ok := v.([]interface{})
			if !ok {
				if fd.Limit == pdl.SPDLimitRequired {
					return fmt.Errorf("字段%s为空或不是Array数据类型", paramName)
				}
				continue
			}
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, arr); err != nil {
				return err
			}
			break
		case pdl.SPD_MAP:
			arr, ok := v.(map[string]interface{})
			if !ok {
				if fd.Limit == pdl.SPDLimitRequired {
					return fmt.Errorf("字段%s为空或不是Object数据类型", paramName)
				}
				continue
			}
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, arr); err != nil {
				return err
			}
			break
		case pdl.SPD_STRUCT:
			arr, ok := v.(map[string]interface{})
			if !ok {
				if fd.Limit == pdl.SPDLimitRequired {
					return fmt.Errorf("字段%s为空或不是Object数据类型", paramName)
				}
				continue
			}
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, arr); err != nil {
				return err
			}
			break
		}
	}
	if err := proto.WriteFieldStop(); err != nil {
		return err
	}
	if err := proto.WriteFieldEnd(); err != nil {
		return err
	}
	if err := proto.WriteStructEnd(); err != nil {
		return err
	}
	return nil
}

func writeListFromArr(qry pdl.PDLQuery, tList *pdl.FileDataType, proto thrift.TProtocol, svc *pdl.FileService, arr []interface{}) error {
	elemType, err := StripTypedef(qry, svc, tList.ElemType)
	if err != nil {
		return err
	}
	size := len(arr)
	if err := proto.WriteListBegin(elemType.Type.ThriftType(), size); err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		if err := writeFileDataTypeToProto(qry, elemType, proto, svc, arr[i]); err != nil {
			return err
		}
	}
	if err := proto.WriteListEnd(); err != nil {
		return err
	}
	return nil
}

func writeSetFromArr(qry pdl.PDLQuery, tSet *pdl.FileDataType, proto thrift.TProtocol, svc *pdl.FileService, arr []interface{}) error {
	elemType, err := StripTypedef(qry, svc, tSet.ElemType)
	if err != nil {
		return err
	}
	size := len(arr)
	if err := proto.WriteSetBegin(elemType.Type.ThriftType(), size); err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		if err := writeFileDataTypeToProto(qry, elemType, proto, svc, arr[i]); err != nil {
			return err
		}
	}
	if err := proto.WriteSetEnd(); err != nil {
		return err
	}
	return nil
}

func writeMapFromMap(qry pdl.PDLQuery, tMap *pdl.FileDataType, proto thrift.TProtocol, svc *pdl.FileService, mp map[string]interface{}) error {
	keyType, err := StripTypedef(qry, svc, tMap.KeyType)
	if err != nil {
		return err
	}
	valType, err := StripTypedef(qry, svc, tMap.ValType)
	if err != nil {
		return err
	}
	size := len(mp)
	if err := proto.WriteMapBegin(keyType.Type.ThriftType(), valType.Type.ThriftType(), size); err != nil {
		return err
	}
	for k, v := range mp {
		if err := writeFileDataTypeToProto(qry, keyType, proto, svc, k); err != nil {
			return err
		}
		if err := writeFileDataTypeToProto(qry, valType, proto, svc, v); err != nil {
			return err
		}
	}
	if err := proto.WriteMapEnd(); err != nil {
		return err
	}
	return nil
}

// map转换成协议调用消息
func MakeCallMessageFromMap(ctx context.Context, qry pdl.PDLQuery, msg string, seqId int32, svc *pdl.FileService, method *pdl.FileServiceMethod, protoFac thrift.TProtocolFactory, mp map[string]interface{}) (data []byte, err error) {
	buf := thrift.NewTMemoryBuffer()
	trans := thrift.NewTFramedTransport(buf)
	proto := protoFac.GetProtocol(trans)
	if err := proto.WriteMessageBegin(msg, thrift.CALL, seqId); err != nil {
		return nil, err
	}
	fmt.Println("svc namespace:", svc.Namespace)

	argStructName := fmt.Sprintf("%s_%s_args", xstring.Camel2UnderScore(svc.Name, "_"), xstring.Camel2UnderScore(method.Name, "_"))

	if err := proto.WriteStructBegin(argStructName); err != nil {
		return nil, err
	}
	args := method.Args
	for _, arg := range args {
		fdId := arg.Id
		limit := arg.Limit
		paramName := xstring.Camel2UnderScore(arg.Name, "_")
		fdType, err := StripTypedef(qry, svc, arg.FieldType)
		if err != nil {
			return nil, err
		}
		if fdType.Type == pdl.SPD_VOID || fdType.Type == pdl.SPD_UNKNOWN {
			continue
		}
		if err := proto.WriteFieldBegin(paramName, fdType.Type.ThriftType(), fdId); err != nil {
			return nil, err
		}
		v, ok := mp[paramName]
		if !ok || v == nil {
			if limit == pdl.SPDLimitRequired {
				return nil, fmt.Errorf("字段%s不能为空", paramName)
			}
		}
		switch fdType.Type {
		case pdl.SPD_STR:
			s := t.String(v)
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, s); err != nil {
				return nil, err
			}
			break
		case pdl.SPD_BOOL:
			b := t.Bool(v)
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, b); err != nil {
				return nil, err
			}
			break
		case pdl.SPD_I08:
			n8 := t.Int8(v)
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, n8); err != nil {
				return nil, err
			}
			break
		case pdl.SPD_I16:
			n8 := t.Int16(v)
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, n8); err != nil {
				return nil, err
			}
			break
		case pdl.SPD_I32:
			n8 := t.Int32(v)
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, n8); err != nil {
				return nil, err
			}
			break
		case pdl.SPD_I64:
			n8 := t.Int64(v)
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, n8); err != nil {
				return nil, err
			}
			break
		case pdl.SPD_DOUBLE:
			n8 := t.Float64(v)
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, n8); err != nil {
				return nil, err
			}
			break
		case pdl.SPD_LIST:
			arr, ok := v.([]interface{})
			if !ok {
				if limit == pdl.SPDLimitRequired {
					return nil, fmt.Errorf("字段%s为空或不是Array数据类型", paramName)
				}
				continue
			}
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, arr); err != nil {
				return nil, err
			}
			break
		case pdl.SPD_SET:
			arr, ok := v.([]interface{})
			if !ok {
				if limit == pdl.SPDLimitRequired {
					return nil, fmt.Errorf("字段%s为空或不是Array数据类型", paramName)
				}
				continue
			}
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, arr); err != nil {
				return nil, err
			}
			break
		case pdl.SPD_MAP:
			arr, ok := v.(map[string]interface{})
			if !ok {
				if limit == pdl.SPDLimitRequired {
					return nil, fmt.Errorf("字段%s为空或不是Object数据类型", paramName)
				}
				continue
			}
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, arr); err != nil {
				return nil, err
			}
			break
		case pdl.SPD_STRUCT:
			arr, ok := v.(map[string]interface{})
			if !ok {
				if limit == pdl.SPDLimitRequired {
					return nil, fmt.Errorf("字段%s为空或不是Object数据类型", paramName)
				}
				continue
			}
			if err := writeFileDataTypeToProto(qry, fdType, proto, svc, arr); err != nil {
				return nil, err
			}
			break
		}
	}
	if err := proto.WriteFieldStop(); err != nil {
		return nil, err
	}
	if err := proto.WriteStructEnd(); err != nil {
		return nil, err
	}
	if err := proto.WriteMessageEnd(); err != nil {
		return nil, err
	}
	if err := proto.Flush(ctx); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// typedef类型转换成实际类型
func StripTypedef(qry pdl.PDLQuery, svc *pdl.FileService, dt *pdl.FileDataType) (*pdl.FileDataType, error) {
	if dt.Type == pdl.SPD_TYPEDEF {
		namespace := getDataTypeNamespace(svc.Namespace, dt)
		_, ot := qry.QryTypeDefByNS(namespace, dt.TypName)
		if ot == nil {
			return nil, fmt.Errorf("协议中未定义类型别名%s", dt.TypName)
		}
		return StripTypedef(qry, svc, ot.OrgType)
	}
	return dt, nil
}

func getDataTypeNamespace(ns string, dt *pdl.FileDataType) string {
	if dt.Namespace == "" {
		return ns
	}
	return dt.Namespace
}

func writeFileDataTypeToProto(qry pdl.PDLQuery, dt *pdl.FileDataType, proto thrift.TProtocol, svc *pdl.FileService, val interface{}) error {
	switch dt.Type {
	case pdl.SPD_STR:
		return proto.WriteString(t.String(val))
	case pdl.SPD_BOOL:
		return proto.WriteBool(t.Bool(val))
	case pdl.SPD_I08:
		return proto.WriteByte(t.Int8(val))
	case pdl.SPD_I16:
		return proto.WriteI16(t.Int16(val))
	case pdl.SPD_I32:
		return proto.WriteI32(t.Int32(val))
	case pdl.SPD_I64:
		return proto.WriteI64(t.Int64(val))
	case pdl.SPD_DOUBLE:
		return proto.WriteDouble(t.Float64(val))
	case pdl.SPD_LIST:
		lstMp, ok := val.([]interface{})
		if !ok || lstMp == nil {
			return fmt.Errorf("值为空或不是Array类型")
		}
		return writeListFromArr(qry, dt, proto, svc, lstMp)
	case pdl.SPD_SET:
		setMp, ok := val.([]interface{})
		if !ok || setMp == nil {
			return fmt.Errorf("值为空或不是Array类型")
		}
		return writeSetFromArr(qry, dt, proto, svc, setMp)
	case pdl.SPD_STRUCT:
		struMp, ok := val.(map[string]interface{})
		if !ok || struMp == nil {
			return fmt.Errorf("值为空或不是Object类型")
		}
		name := getDataTypeNamespace(svc.Namespace, dt)
		_, struType := qry.QryTypeByNS(name, dt.TypName)
		if struType == nil {
			return fmt.Errorf("协议中缺少数据类型%s的定义", dt.TypName)
		}
		return writeStructFromMap(qry, struType, proto, svc, struMp)
	}
	return nil
}

func MapFromCallReply(qry pdl.PDLQuery, msg string, svc *pdl.FileService, seqId int32, method *pdl.FileServiceMethod, proto thrift.TProtocol) (val interface{}, err error) {
	_, msgType, _, err := proto.ReadMessageBegin()
	if err != nil {
		xlog.Error(err)
		return nil, err
	}
	//if nSeqId != seqId {
	//	return nil, fmt.Errorf("期望序列%d实际返回%d", nSeqId, seqId)
	//}
	if msgType == thrift.EXCEPTION {
		// 服务端返回异常类型
		appErr := thrift.NewTApplicationException(0, "")
		if err := appErr.Read(proto); err != nil {
			xlog.Error(err)
			return nil, err
		}
		if err := proto.ReadMessageEnd(); err != nil {
			xlog.Error(err)
			return nil, err
		}
		xlog.Error(appErr)
		return nil, appErr
	}

	if msgType != thrift.REPLY {
		return nil, errors.New("数据协议错误")
	}
	var successData, exceptionData interface{} = nil, nil

	if _, err := proto.ReadStructBegin(); err != nil {
		xlog.Error(err)
		return nil, err
	}
	if method.Result == nil && method.Result.Type == pdl.SPD_VOID {
		return true, nil
	}
	resultType, err := StripTypedef(qry, svc, method.Result)
	if err != nil {
		xlog.Error(err)
		return nil, err
	}
	exceptionType := pdl.NewBaseFileDataType(pdl.SPD_VOID)
	if method.Exception != nil {
		exceptionType, err = StripTypedef(qry, svc, exceptionType)
		if err != nil {
			xlog.Error(err)
			return nil, err
		}
	}

	for {
		fdName, fdType, fdId, err := proto.ReadFieldBegin()
		if err != nil {
			xlog.Error(err)
			return nil, err
		}
		if fdType == thrift.STOP {
			break
		}
		if fdType == thrift.VOID {
			if err := proto.ReadFieldEnd(); err != nil {
				xlog.Error(err)
				return nil, err
			}
			continue
		}

		if fdId == 1 || fdName == "success" {
			// 成功
			successData, err = makeDataFromProto(qry, proto, svc, fdType, resultType)
			if err != nil {
				xlog.Error(err)
				return nil, err
			}
		} else if fdId == 2 || fdName == "exception" {
			// 失败
			exceptionData, err = makeDataFromProto(qry, proto, svc, fdType, exceptionType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
		} else {
			if err := proto.Skip(fdType); err != nil {
				xlog.Debug(err)
				return nil, err
			}
		}
		if err := proto.ReadFieldEnd(); err != nil {
			xlog.Debug(err)
			return nil, err
		}
	}
	if err := proto.ReadStructEnd(); err != nil {
		xlog.Debug(err)
		return nil, err
	}
	if err := proto.ReadMessageEnd(); err != nil {
		xlog.Debug(err)
		return nil, err
	}
	if successData != nil {
		return successData, nil
	}
	return exceptionData, err
}

func makeDataFromProto(qry pdl.PDLQuery, proto thrift.TProtocol, svc *pdl.FileService, t thrift.TType, dt *pdl.FileDataType) (result interface{}, err error) {
	switch dt.Type {
	case pdl.SPD_STR:
		return readBasicDataFromProto(qry, proto, thrift.STRING, dt)
	case pdl.SPD_BOOL:
		return readBasicDataFromProto(qry, proto, thrift.BOOL, dt)
	case pdl.SPD_I08:
		return readBasicDataFromProto(qry, proto, thrift.BYTE, dt)
	case pdl.SPD_I16:
		return readBasicDataFromProto(qry, proto, thrift.I16, dt)
	case pdl.SPD_I32:
		return readBasicDataFromProto(qry, proto, thrift.I32, dt)
	case pdl.SPD_I64:
		return readBasicDataFromProto(qry, proto, thrift.I64, dt)
	case pdl.SPD_DOUBLE:
		return readBasicDataFromProto(qry, proto, thrift.DOUBLE, dt)
	case pdl.SPD_LIST:
		return readListDataFromProto(qry, proto, svc, thrift.LIST, dt)
	case pdl.SPD_SET:
		return readSetDataFromProto(qry, proto, svc, thrift.SET, dt)
	case pdl.SPD_MAP:
		return readMapDataFromProto(qry, proto, svc, thrift.MAP, dt)
	case pdl.SPD_STRUCT, pdl.SPD_EXCEPTION:
		return readStructDataFromProto(qry, proto, svc, thrift.STRUCT, dt)
	}
	return
}

func readBasicDataFromProto(qry pdl.PDLQuery, proto thrift.TProtocol, t thrift.TType, dt *pdl.FileDataType) (result interface{}, err error) {
	switch t {
	case thrift.STRING:
		return proto.ReadString()
	case thrift.BOOL:
		return proto.ReadBool()
	case thrift.BYTE:
		return proto.ReadByte()
	case thrift.I16:
		return proto.ReadI16()
	case thrift.I32:
		return proto.ReadI32()
	case thrift.I64:
		return proto.ReadI64()
	case thrift.DOUBLE:
		return proto.ReadDouble()
	}
	return nil, nil
}

func readListDataFromProto(qry pdl.PDLQuery, proto thrift.TProtocol, svc *pdl.FileService, t thrift.TType, dt *pdl.FileDataType) (result []interface{}, err error) {
	elemType, size, err := proto.ReadListBegin()
	if err != nil {
		return nil, err
	}
	result = make([]interface{}, size)
	for i := 0; i < size; i++ {
		var v interface{} = nil
		switch elemType {
		case thrift.STRING:
			v, err = proto.ReadString()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.BOOL:
			v, err = proto.ReadBool()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.BYTE:
			v, err = proto.ReadByte()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.I16:
			v, err = proto.ReadI16()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.I32:
			v, err = proto.ReadI32()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.I64:
			v, err = proto.ReadI64()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.DOUBLE:
			v, err = proto.ReadDouble()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.LIST:
			v, err = readListDataFromProto(qry, proto, svc, thrift.LIST, dt.ElemType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.SET:
			v, err = readSetDataFromProto(qry, proto, svc, thrift.SET, dt.ElemType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.MAP:
			v, err = readMapDataFromProto(qry, proto, svc, thrift.MAP, dt.ElemType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.STRUCT:
			v, err = readStructDataFromProto(qry, proto, svc, thrift.STRUCT, dt.ElemType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		}
		result[i] = v
	}
	if err := proto.ReadListEnd(); err != nil {
		xlog.Debug(err)
		return nil, err
	}
	return result, nil
}

func readSetDataFromProto(qry pdl.PDLQuery, proto thrift.TProtocol, svc *pdl.FileService, t thrift.TType, dt *pdl.FileDataType) (result []interface{}, err error) {
	elemType, size, err := proto.ReadSetBegin()
	if err != nil {
		return nil, err
	}
	result = make([]interface{}, size)
	for i := 0; i < size; i++ {
		var v interface{} = nil
		switch elemType {
		case thrift.STRING:
			v, err = proto.ReadString()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.BOOL:
			v, err = proto.ReadBool()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.BYTE:
			v, err = proto.ReadByte()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.I16:
			v, err = proto.ReadI16()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.I32:
			v, err = proto.ReadI32()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.I64:
			v, err = proto.ReadI64()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.DOUBLE:
			v, err = proto.ReadDouble()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.LIST:
			v, err = readListDataFromProto(qry, proto, svc, thrift.LIST, dt.ElemType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.SET:
			v, err = readSetDataFromProto(qry, proto, svc, thrift.SET, dt.ElemType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.MAP:
			v, err = readMapDataFromProto(qry, proto, svc, thrift.MAP, dt.ElemType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.STRUCT:
			v, err = readStructDataFromProto(qry, proto, svc, thrift.STRUCT, dt.ElemType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		}
		result[i] = v
	}
	if err := proto.ReadSetEnd(); err != nil {
		return nil, err
	}
	return result, nil
}

func readMapDataFromProto(qry pdl.PDLQuery, proto thrift.TProtocol, svc *pdl.FileService, tt thrift.TType, dt *pdl.FileDataType) (result map[string]interface{}, err error) {
	keyType, valType, size, err := proto.ReadMapBegin()
	if err != nil {
		return nil, err
	}
	result = make(map[string]interface{})
	for i := 0; i < size; i++ {
		var k, v interface{} = nil, nil
		switch keyType {
		case thrift.STRING:
			k, err = proto.ReadString()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.BOOL:
			k, err = proto.ReadBool()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.BYTE:
			k, err = proto.ReadByte()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.I16:
			k, err = proto.ReadI16()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.I32:
			k, err = proto.ReadI32()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.I64:
			k, err = proto.ReadI64()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.DOUBLE:
			k, err = proto.ReadDouble()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		}

		switch valType {
		case thrift.STRING:
			v, err = proto.ReadString()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.BOOL:
			v, err = proto.ReadBool()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.BYTE:
			v, err = proto.ReadByte()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.I16:
			v, err = proto.ReadI16()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.I32:
			v, err = proto.ReadI32()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.I64:
			v, err = proto.ReadI64()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.DOUBLE:
			v, err = proto.ReadDouble()
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.LIST:
			v, err = readListDataFromProto(qry, proto, svc, thrift.LIST, dt.ElemType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.SET:
			v, err = readSetDataFromProto(qry, proto, svc, thrift.SET, dt.ElemType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.MAP:
			v, err = readMapDataFromProto(qry, proto, svc, thrift.MAP, dt.ElemType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.STRUCT:
			v, err = readStructDataFromProto(qry, proto, svc, thrift.STRUCT, dt.ElemType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		}
		result[t.String(k)] = v
	}
	return result, nil
}

func readStructDataFromProto(qry pdl.PDLQuery, proto thrift.TProtocol, svc *pdl.FileService, t thrift.TType, dt *pdl.FileDataType) (result map[string]interface{}, err error) {
	namespace := getDataTypeNamespace(svc.Namespace, dt)
	_, stru := qry.QryTypeByNS(namespace, dt.TypName)
	if stru == nil {
		return nil, fmt.Errorf("协议中未定义类型%s", dt.TypName)
	}
	fields := make(map[int16]*pdl.FileDataField)
	for _, f := range stru.Fields {
		fields[f.Id] = f
	}
	_, err = proto.ReadStructBegin()
	if err != nil {
		return nil, err
	}
	result = make(map[string]interface{})

	for {
		var v interface{} = nil
		fdName, fdType, fdId, err := proto.ReadFieldBegin()
		if err != nil {
			xlog.Debug(err)
			return nil, err
		}
		if fdType == thrift.STOP {
			break
		}
		if fdType == thrift.VOID {
			if err = proto.ReadFieldEnd(); err != nil {
				xlog.Debug(err)
				return nil, err
			}
			continue
		}
		fdInfo, fdOk := fields[fdId]
		if !fdOk {
			// 字段id不存在
			if err := proto.Skip(fdType); err != nil {
				xlog.Debug(err)
				return nil, err
			}
			if err := proto.ReadFieldEnd(); err != nil {
				xlog.Debug(err)
				return nil, err
			}
			continue
		}
		switch fdType {
		case thrift.STRING, thrift.BOOL, thrift.BYTE, thrift.I16, thrift.I32, thrift.I64, thrift.DOUBLE:
			v, err = readBasicDataFromProto(qry, proto, fdType, fdInfo.FieldType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.LIST:
			v, err = readListDataFromProto(qry, proto, svc, fdType, fdInfo.FieldType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.SET:
			v, err = readSetDataFromProto(qry, proto, svc, fdType, fdInfo.FieldType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.MAP:
			v, err = readMapDataFromProto(qry, proto, svc, fdType, fdInfo.FieldType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		case thrift.STRUCT:
			v, err = readStructDataFromProto(qry, proto, svc, fdType, fdInfo.FieldType)
			if err != nil {
				xlog.Debug(err)
				return nil, err
			}
			break
		}
		if err := proto.ReadFieldEnd(); err != nil {
			xlog.Debug(err)
			return nil, err
		}
		if v != nil {
			fdResultName := xstring.Camel2UnderScore(fdName, "_")
			result[fdResultName] = v
		}
	}
	if err := proto.ReadStructEnd(); err != nil {
		xlog.Debug(err)
		return nil, err
	}
	return result, nil
}
