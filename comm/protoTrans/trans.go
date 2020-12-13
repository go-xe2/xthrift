package protoTrans

import "github.com/apache/thrift/lib/go/thrift"

func ProtocolTransform(in thrift.TProtocol, out thrift.TProtocol) (msgName string, msgType thrift.TMessageType, seqId int32, err error) {
	msgName, msgType, seqId, err = in.ReadMessageBegin()
	if err != nil {
		return
	}
	if err = out.WriteMessageBegin(msgName, msgType, seqId); err != nil {
		return
	}
	if msgType == thrift.ONEWAY+1 {
		if err = in.ReadMessageEnd(); err != nil {
			return
		}
		if err = out.WriteMessageEnd(); err != nil {
			return
		}
		return
	}

	if err = protocolStructTransform(in, out); err != nil {
		return
	}
	if err = in.ReadMessageEnd(); err != nil {
		return
	}
	if err = out.WriteMessageEnd(); err != nil {
		return
	}
	return
}

func protocolListTransform(in thrift.TProtocol, out thrift.TProtocol) error {
	elem, size, err := in.ReadListBegin()
	if err != nil {
		return err
	}
	if err := out.WriteListBegin(elem, size); err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		switch elem {
		case thrift.STRING:
			s, err := in.ReadString()
			if err != nil {
				return err
			}
			if err := out.WriteString(s); err != nil {
				return err
			}
			break
		case thrift.BOOL:
			b, err := in.ReadBool()
			if err != nil {
				return err
			}
			if err := out.WriteBool(b); err != nil {
				return err
			}
			break
		case thrift.BYTE:
			n8, err := in.ReadByte()
			if err != nil {
				return err
			}
			if err := out.WriteByte(int8(n8)); err != nil {
				return err
			}
			break
		case thrift.I16:
			n16, err := in.ReadI16()
			if err != nil {
				return err
			}
			if err := out.WriteI16(n16); err != nil {
				return err
			}
			break
		case thrift.I32:
			n32, err := in.ReadI32()
			if err != nil {
				return err
			}
			if err := out.WriteI32(n32); err != nil {
				return err
			}
			break
		case thrift.I64:
			n64, err := in.ReadI64()
			if err != nil {
				return err
			}
			if err := out.WriteI64(n64); err != nil {
				return err
			}
			break
		case thrift.DOUBLE:
			dl, err := in.ReadDouble()
			if err != nil {
				return err
			}
			if err := out.WriteDouble(dl); err != nil {
				return err
			}
			break
		case thrift.LIST:
			if err := protocolListTransform(in, out); err != nil {
				return err
			}
			break
		case thrift.SET:
			if err := protocolSetTransform(in, out); err != nil {
				return err
			}
			break
		case thrift.STRUCT:
			if err := protocolStructTransform(in, out); err != nil {
				return err
			}
			break
		}
	}
	if err := in.ReadListEnd(); err != nil {
		return err
	}
	if err := out.WriteListEnd(); err != nil {
		return err
	}
	return nil
}

func protocolSetTransform(in thrift.TProtocol, out thrift.TProtocol) error {
	elem, size, err := in.ReadSetBegin()
	if err != nil {
		return err
	}
	if err := out.WriteSetBegin(elem, size); err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		switch elem {
		case thrift.STRING:
			s, err := in.ReadString()
			if err != nil {
				return err
			}
			if err := out.WriteString(s); err != nil {
				return err
			}
			break
		case thrift.BOOL:
			b, err := in.ReadBool()
			if err != nil {
				return err
			}
			if err := out.WriteBool(b); err != nil {
				return err
			}
			break
		case thrift.BYTE:
			n8, err := in.ReadByte()
			if err != nil {
				return err
			}
			if err := out.WriteByte(int8(n8)); err != nil {
				return err
			}
			break
		case thrift.I16:
			n16, err := in.ReadI16()
			if err != nil {
				return err
			}
			if err := out.WriteI16(n16); err != nil {
				return err
			}
			break
		case thrift.I32:
			n32, err := in.ReadI32()
			if err != nil {
				return err
			}
			if err := out.WriteI32(n32); err != nil {
				return err
			}
			break
		case thrift.I64:
			n64, err := in.ReadI64()
			if err != nil {
				return err
			}
			if err := out.WriteI64(n64); err != nil {
				return err
			}
			break
		case thrift.DOUBLE:
			dl, err := in.ReadDouble()
			if err != nil {
				return err
			}
			if err := out.WriteDouble(dl); err != nil {
				return err
			}
			break
		case thrift.LIST:
			if err := protocolListTransform(in, out); err != nil {
				return err
			}
			break
		case thrift.SET:
			if err := protocolSetTransform(in, out); err != nil {
				return err
			}
			break
		case thrift.STRUCT:
			if err := protocolStructTransform(in, out); err != nil {
				return err
			}
			break
		}
	}
	if err := in.ReadSetEnd(); err != nil {
		return err
	}
	if err := out.WriteSetEnd(); err != nil {
		return err
	}
	return nil
}

func protocolStructTransform(in thrift.TProtocol, out thrift.TProtocol) error {
	if s, err := in.ReadStructBegin(); err != nil {
		return err
	} else {
		if err := out.WriteStructBegin(s); err != nil {
			return err
		}
	}
	for {
		fdName, fdType, fdId, err := in.ReadFieldBegin()
		if err != nil {
			return err
		}
		if fdType == thrift.STOP {
			if err := out.WriteFieldStop(); err != nil {
				return err
			}
			break
		}
		if fdType == thrift.VOID {
			if err := in.ReadFieldEnd(); err != nil {
				return err
			}
			continue
		}
		if err := out.WriteFieldBegin(fdName, fdType, fdId); err != nil {
			return err
		}
		switch fdType {
		case thrift.STRING:
			s, err := in.ReadString()
			if err != nil {
				return err
			}
			if err := out.WriteString(s); err != nil {
				return err
			}
			break
		case thrift.BOOL:
			b, err := in.ReadBool()
			if err != nil {
				return err
			}
			if err := out.WriteBool(b); err != nil {
				return err
			}
			break
		case thrift.BYTE:
			n8, err := in.ReadByte()
			if err != nil {
				return err
			}
			if err := out.WriteByte(int8(n8)); err != nil {
				return err
			}
			break
		case thrift.I16:
			n16, err := in.ReadI16()
			if err != nil {
				return err
			}
			if err := out.WriteI16(n16); err != nil {
				return err
			}
			break
		case thrift.I32:
			n32, err := in.ReadI32()
			if err != nil {
				return err
			}
			if err := out.WriteI32(n32); err != nil {
				return err
			}
			break
		case thrift.I64:
			n64, err := in.ReadI64()
			if err != nil {
				return err
			}
			if err := out.WriteI64(n64); err != nil {
				return err
			}
			break
		case thrift.DOUBLE:
			dl, err := in.ReadDouble()
			if err != nil {
				return err
			}
			if err := out.WriteDouble(dl); err != nil {
				return err
			}
			break
		case thrift.LIST:
			if err := protocolListTransform(in, out); err != nil {
				return err
			}
			break
		case thrift.SET:
			if err := protocolSetTransform(in, out); err != nil {
				return err
			}
			break
		case thrift.STRUCT:
			if err := protocolStructTransform(in, out); err != nil {
				return err
			}
			break
		}
		if err := in.ReadFieldEnd(); err != nil {
			return err
		}
		if err := out.WriteFieldEnd(); err != nil {
			return err
		}
	}
	if err := in.ReadStructEnd(); err != nil {
		return err
	}
	if err := out.WriteStructEnd(); err != nil {
		return err
	}
	return nil
}
