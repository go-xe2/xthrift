/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 14:17
* Description:
*****************************************************************/

package gcontext

import (
	"bytes"
	"fmt"
)

func (p *TWriter) GenStrReadCode(protoVar, val, err string) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("%s, %s := %s.ReadString()\n", val, err, protoVar))
	result.WriteString(fmt.Sprintf("if %s != nil {\n", err))
	result.WriteString(fmt.Sprintf("\treturn %s\n", err))
	result.WriteString("}\n")
	return result.String()
}

func (p *TWriter) GenBoolReadCode(protoVar, val, err string) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("%s, %s := %s.ReadBool()\n", val, err, protoVar))
	result.WriteString(fmt.Sprintf("if %s != nil {\n", err))
	result.WriteString(fmt.Sprintf("\treturn %s\n", err))
	result.WriteString("}\n")
	return result.String()
}

func (p *TWriter) GenI08ReadCode(protoVar, val, err string) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("%s, %s := %s.ReadByte()\n", val, err, protoVar))
	result.WriteString(fmt.Sprintf("if %s != nil {\n", err))
	result.WriteString(fmt.Sprintf("\treturn %s\n", err))
	result.WriteString("}\n")
	return result.String()
}

func (p *TWriter) GenI16ReadCode(protoVar, val, err string) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("%s, %s := %s.ReadI16()\n", val, err, protoVar))
	result.WriteString(fmt.Sprintf("if %s != nil {\n", err))
	result.WriteString(fmt.Sprintf("\treturn %s\n", err))
	result.WriteString("}\n")
	return result.String()
}

func (p *TWriter) GenI32ReadCode(protoVar, val, err string) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("%s, %s := %s.ReadI32()\n", val, err, protoVar))
	result.WriteString(fmt.Sprintf("if %s != nil {\n", err))
	result.WriteString(fmt.Sprintf("\treturn %s\n", err))
	result.WriteString("}\n")
	return result.String()
}

func (p *TWriter) GenI64ReadCode(protoVar, val, err string) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("%s, %s := %s.ReadI64()\n", val, err, protoVar))
	result.WriteString(fmt.Sprintf("if %s != nil {\n", err))
	result.WriteString(fmt.Sprintf("\treturn %s\n", err))
	result.WriteString("}\n")
	return result.String()
}

func (p *TWriter) GenDoubleReadCode(protoVar, val, err string) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("%s, %s := %s.ReadDouble()\n", val, err, protoVar))
	result.WriteString(fmt.Sprintf("if %s != nil {\n", err))
	result.WriteString(fmt.Sprintf("\treturn %s\n", err))
	result.WriteString("}\n")
	return result.String()
}
