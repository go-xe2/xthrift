/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-14 16:06
* Description:
*****************************************************************/

package gcontext

import (
	"bytes"
	"fmt"
)

func (p *TWriter) GenStrWriteCode(protoVar, val, err string) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("if %s := %s.WriteString(%s); %s != nil {\n", err, protoVar, val, err))
	result.WriteString(fmt.Sprintf("\treturn %s\n", err))
	result.WriteString("}\n")
	return result.String()
}

func (p *TWriter) GenBoolWriteCode(protoVar, val, err string) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("if  %s := %s.WriteBool(%s); %s != nil {\n", err, protoVar, val, err))
	result.WriteString(fmt.Sprintf("\treturn %s\n", err))
	result.WriteString("}\n")
	return result.String()
}

func (p *TWriter) GenI08WriteCode(protoVar, val, err string) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("if %s := %s.WriteByte(%s); %s != nil {\n", err, protoVar, val, err))
	result.WriteString(fmt.Sprintf("\treturn %s\n", err))
	result.WriteString("}\n")
	return result.String()
}

func (p *TWriter) GenI16WriteCode(protoVar, val, err string) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("if %s := %s.WriteI16(%s); %s != nil {\n", err, protoVar, val, err))
	result.WriteString(fmt.Sprintf("\treturn %s\n", err))
	result.WriteString("}\n")
	return result.String()
}

func (p *TWriter) GenI32WriteCode(protoVar, val, err string) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("if %s := %s.WriteI32(%s); %s != nil {\n", err, protoVar, val, err))
	result.WriteString(fmt.Sprintf("\treturn %s\n", err))
	result.WriteString("}\n")
	return result.String()
}

func (p *TWriter) GenI64WriteCode(protoVar, val, err string) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("if %s := %s.WriteI64(%s); %s != nil {\n", err, protoVar, val, err))
	result.WriteString(fmt.Sprintf("\treturn %s\n", err))
	result.WriteString("}\n")
	return result.String()
}

func (p *TWriter) GenDoubleWriteCode(protoVar, val, err string) string {
	result := bytes.NewBufferString("")
	result.WriteString(fmt.Sprintf("if %s := %s.WriteDouble(%s); %s != nil {\n", err, protoVar, val, err))
	result.WriteString(fmt.Sprintf("\treturn %s\n", err))
	result.WriteString("}\n")
	return result.String()
}
