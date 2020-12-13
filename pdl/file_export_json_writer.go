/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-20 17:26
* Description:
*****************************************************************/

package pdl

import (
	"bytes"
	"fmt"
	"sort"
)

type tJsonWriter struct {
	w *bytes.Buffer
}

func newJsonWriter() *tJsonWriter {
	return &tJsonWriter{
		w: bytes.NewBuffer([]byte{}),
	}
}

func (pw *tJsonWriter) write(str string) {
	if _, err := pw.w.Write([]byte(str)); err != nil {
		panic(err)
	}
}

func (pw *tJsonWriter) WriteBegin() error {
	pw.write("{\n")
	return nil
}

func (pw *tJsonWriter) WriteNamespace(namespace string) error {
	pw.write(" ")
	pw.write(`"namespace": `)
	pw.write(fmt.Sprintf(`"%s"`, namespace))
	pw.write(",\n")
	return nil
}

func (pw *tJsonWriter) WriteImports(imports []string) error {
	pw.write(" ")
	s := ""
	for _, k := range imports {
		s += fmt.Sprintf(`"%s",`, k)
	}
	if len(s) > 0 {
		pw.write(fmt.Sprintf(`"imports": [%s]`, s[:len(s)-1]))
	} else {
		pw.write(`"imports": []`)
	}
	pw.write(",\n")
	return nil
}

func (pw *tJsonWriter) WriteBasicBegin() error {
	pw.write(" ")
	pw.write(`"basic": [`)
	return nil
}

func (pw *tJsonWriter) WriteBasic(basicTypes []string) error {
	s := ""
	for _, v := range basicTypes {
		s += fmt.Sprintf(`"%s",`, v)
	}
	pw.write(s[:len(s)-1])
	return nil
}

func (pw *tJsonWriter) WriteBasicEnd() error {
	pw.write("],\n")
	return nil
}

func (pw *tJsonWriter) WriteTypeDefBegin() error {
	pw.write(" ")
	pw.write("\"typeDefs\": {")
	pw.write("\n")
	return nil
}

func (pw *tJsonWriter) WriteTypeDefs(mp map[string]*FileTypeDef) error {
	str := ""
	for k, v := range mp {
		str += "  "
		str += fmt.Sprintf(`"%s": "%s",`, k, v.OrgType.Name())
		str += "\n"
	}
	if str == "" {
		pw.write("")
	} else {
		pw.write(str[:len(str)-2])
	}
	return nil
}

func (pw *tJsonWriter) WriteTypeDefEnd() error {
	pw.write("\n },\n")
	return nil
}

func (pw *tJsonWriter) WriteTypesBegin() error {
	pw.write(" ")
	pw.write(`"types": {`)
	pw.write("\n")
	return nil
}

func (pw *tJsonWriter) WriteTypes(mp map[string]*FileStruct) error {
	size := len(mp)
	i := 0
	for k, v := range mp {
		i++
		if e := pw.WriteStructBegin("  ", k); e != nil {
			return e
		}
		if e := pw.WriteStruct("  ", v); e != nil {
			return e
		}
		if e := pw.WriteStructEnd("  "); e != nil {
			return e
		}
		if i < size {
			pw.write(",\n")
		}
	}
	return nil
}

func (pw *tJsonWriter) WriteTypesEnd() error {
	pw.write("\n },\n")
	return nil
}

func (pw *tJsonWriter) WriteInterfacesBegin() error {
	pw.write("\n ")
	pw.write(`"interfaces": {`)
	pw.write("\n")
	return nil
}

func (pw *tJsonWriter) WriteInterfaces(interfaces map[string]*FileService) error {
	var size = len(interfaces)
	i := 0
	for k, v := range interfaces {
		i++
		if e := pw.WriteServiceBegin("  ", k); e != nil {
			return e
		}
		if e := pw.WriteService("  ", v); e != nil {
			return e
		}
		if e := pw.WriteServiceEnd("  "); e != nil {
			return e
		}
		if i < size {
			pw.write("  ,\n")
		}
	}
	return nil
}

func (pw *tJsonWriter) WriteInterfacesEnd() error {
	pw.write("\n }\n")
	return nil
}

func (pw *tJsonWriter) WriteServiceBegin(indent string, svcName string) error {
	pw.write(indent)
	pw.write(fmt.Sprintf(`"%s": {`, svcName))
	pw.write("\n")
	return nil
}

func (pw *tJsonWriter) WriteService(indent string, service *FileService) error {
	methods := service.Methods
	count := len(methods)
	i := 0
	for k, v := range methods {
		i++
		if e := pw.WriteServiceMethodBegin(indent+" ", k); e != nil {
			return e
		}
		if e := pw.WriteServiceMethod(indent+" ", v); e != nil {
			return e
		}
		if e := pw.WriteServiceMethodEnd(indent + " "); e != nil {
			return e
		}
		if i < count {
			pw.write(",\n")
		}
	}
	return nil
}

func (pw *tJsonWriter) WriteServiceEnd(indent string) error {
	pw.write("\n")
	pw.write(indent)
	pw.write("}")
	return nil
}

func (pw *tJsonWriter) WriteServiceMethodBegin(indent string, methodName string) error {
	pw.write(indent)
	pw.write(fmt.Sprintf(`"%s": {`, methodName))
	pw.write("\n")
	return nil
}

func (pw *tJsonWriter) WriteServiceMethod(indent string, method *FileServiceMethod) error {
	pw.write(indent + " ")
	pw.write(fmt.Sprintf(`"summary": "%s",`, method.Summary))
	args := method.Args
	pw.write("\n")
	pw.write(indent + " ")
	pw.write(`"args": {`)
	pw.write("\n")
	size := len(args)
	i := 0
	sort.Slice(args, func(i, j int) bool {
		return args[i].Id-args[j].Id < 0
	})
	for _, f := range args {
		i++
		if e := pw.WriteFieldBegin(indent+"  ", f.Name); e != nil {
			return e
		}
		if e := pw.WriteField(indent+"  ", f); e != nil {
			return e
		}
		if e := pw.WriteFieldEnd(indent + "  "); e != nil {
			return e
		}
		if i < size {
			pw.write(",\n")
		}
	}
	pw.write("\n")
	pw.write(indent + " },\n")
	pw.write(indent + "\"results\": \"" + method.Result.Name() + "\"")
	if method.Exception == nil || method.Exception.Type == SPD_VOID {
		pw.write("\n")
	} else {
		pw.write(indent + ",\n\"throw\": \"" + method.Exception.Name() + "\"")
	}
	return nil
}

func (pw *tJsonWriter) WriteServiceMethodEnd(indent string) error {
	pw.write("\n")
	pw.write(indent + "}")
	return nil
}

func (pw *tJsonWriter) WriteFieldBegin(indent string, fieldName string) error {
	pw.write(indent)
	pw.write(fmt.Sprintf("\"%s\": {", fieldName))
	return nil
}

func (pw *tJsonWriter) WriteField(indent string, field *FileDataField) error {
	pw.write(fmt.Sprintf(`"id": %d,`, field.Id))
	pw.write(fmt.Sprintf(`"type": "%s"`, field.FieldType.Name()))
	if field.Summary != "" {
		pw.write(fmt.Sprintf(`,"summary": "%s"`, field.Summary))
	}
	if field.Rule != "" {
		pw.write(fmt.Sprintf(`,"valid": "%s"`, field.Rule))
	}
	if field.Limit != SPDLimitRequired {
		pw.write(fmt.Sprintf(`,"limit": "%s"`, field.Limit))
	}
	return nil
}

func (pw *tJsonWriter) WriteFieldEnd(indent string) error {
	pw.write("}")
	return nil
}

func (pw *tJsonWriter) WriteStructBegin(indent string, struName string) error {
	pw.write(indent)
	pw.write(fmt.Sprintf("\"%s\": {\n", struName))
	return nil
}

func (pw *tJsonWriter) WriteStruct(indent string, stru *FileStruct) error {
	pw.write(indent + " " + `"type": `)
	if stru.Type.Type == SPD_STRUCT {
		pw.write("\"struct\",")
	} else {
		pw.write("\"exception\",")
	}
	pw.write("\n")

	pw.write(indent + " ")
	pw.write(fmt.Sprintf(`"summary": "%s"`, stru.Summary))
	pw.write(",")
	pw.write("\n")

	pw.write(indent + " ")
	pw.write(`"fields": {`)
	pw.write("\n")
	fields := make([]*FileDataField, len(stru.Fields))
	i := 0
	for _, f := range stru.Fields {
		fields[i] = f
		i++
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Id-fields[j].Id < 0
	})

	size := len(fields)
	i = 0

	for _, v := range fields {
		i++
		if e := pw.WriteFieldBegin(indent+"  ", v.Name); e != nil {
			return e
		}
		if e := pw.WriteField(indent+"  ", v); e != nil {
			return e
		}
		if e := pw.WriteFieldEnd(indent + "  "); e != nil {
			return e
		}
		if i < size {
			pw.write(",\n")
		}
	}
	pw.write("\n")
	pw.write(indent + " }")
	return nil
}

func (pw *tJsonWriter) WriteStructEnd(indent string) error {
	pw.write("\n")
	pw.write(indent + "}")
	return nil
}

func (pw *tJsonWriter) WriteEnd() error {
	pw.write("\n}")
	return nil
}

func (p *tJsonWriter) Data() []byte {
	return p.w.Bytes()
}
