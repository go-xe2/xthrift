package pdl

import (
	"fmt"
	"github.com/go-xe2/x/encoding/xparser"
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/type/t"
	"sort"
)

func (p *FileData) loadImports(arr []interface{}) error {
	for _, v := range arr {
		s := t.String(v)
		if s != "" {
			p.AddImport(s)
		}
	}
	return nil
}

func (p *FileData) loadTypedefs(mp map[string]interface{}) error {
	for k, v := range mp {
		s := t.String(v)
		if s == "" {
			return fmt.Errorf("类型别名定义错误, 类型别名%s的数据类型未定义", k)
		}
		if k == s {
			return fmt.Errorf("类型别名定义错误，定义类型名%s与类型%s相同", k, s)
		}
		p.AddTypedef(NewFileTypeDef(k, NewFileDataTypeFromStr(s)))
	}
	return nil
}

func (p *FileData) loadTypes(mp map[string]interface{}) error {
	for k, v := range mp {
		m, ok := v.(map[string]interface{})
		if !ok {
			fmt.Errorf("数据类型%s定义错误，未定义type,fields等待相关字段", k)
		}
		szType := t.String(m["type"])
		summary := t.String(m["summary"])
		fields := t.Map(m["fields"])
		if szType == "" {
			fmt.Errorf("数据类型%s定义错误,未定义数据类型字段type", k)
		}
		if fields == nil {
			fmt.Errorf("数据类型%s定义错误，未定义字段列表fields", k)
		}
		if szType != "struct" && szType != "exception" {
			fmt.Errorf("数据类型%s定义类型%s错误，只能定义为struct和exception类型", k, szType)
		}
		var inst *FileStruct
		if szType == "struct" {
			inst = NewFileStructType(p.Namespace, k, summary)
		} else {
			inst = NewFileExceptType(p.Namespace, k, summary)
		}

		fdIds := make(map[int16]int, 0)
		for fdName, fv := range fields {
			fm, ok := fv.(map[string]interface{})
			if !ok {
				return fmt.Errorf("数据类型%s的字段%s定义错误，未定义id,type等相关字段", k, fdName)
			}
			fdId := t.Int16(fm["id"])
			szFdType := t.String(fm["type"])
			fdLimit := t.String(fm["limit"])
			fdRule := t.String(fm["rule"])
			fdSummary := t.String(fm["summary"])
			if fdId < 1 {
				return fmt.Errorf("数据类型%s的字段%sID未定义或定义或定义小于1", k, fdName)
			}
			if _, ok := fdIds[fdId]; ok {
				return fmt.Errorf("数据类型%s的字段%sID定义重复%d", k, fdName, fdId)
			}
			nLimit := SPDLimitRequired
			if fdLimit == "optional" {
				nLimit = SPDLimitOptional
			}
			fdIds[fdId] = 1
			fdType := NewFileDataTypeFromStr(szFdType)
			field := NewFileDataField(fdId, fdName, fdType, fdSummary)
			field.SetLimit(nLimit)
			field.SetRule(fdRule)
			inst.AddField(field)
		}
		if inst.Fields != nil {
			sort.Slice(inst.Fields, func(i, j int) bool {
				return inst.Fields[i].Id-inst.Fields[j].Id < 0
			})
		}
		p.AddType(inst)
	}
	return nil
}

func (p *FileData) loadServices(mp map[string]interface{}) error {
	for svcName, v := range mp {
		mmp, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("服务%s定义错误，未定义服务接口", svcName)
		}
		serviceSummary := ""
		service := NewFileService(p.Namespace, svcName, "")
		for methodName, v1 := range mmp {
			if methodName == svcName+"Summary" {
				serviceSummary = t.String(v1)
				continue
			}
			methodMp, ok := v1.(map[string]interface{})
			if !ok {
				return fmt.Errorf("服务接口%s.%s定义错误，未定义输入参数及返回参数类型等字段", svcName, methodName)
			}
			methodSummary := t.String(methodMp["summary"])
			argsMp := t.Map(methodMp["args"])
			result := t.String(methodMp["results"])
			except := t.String(methodMp["throw"])
			if result == "" {
				return fmt.Errorf("服务接口%s.%s未定义返回类型", svcName, methodName)
			}
			if except == "" {
				except = "void"
			}
			method := NewFileServiceMethod(methodName, methodSummary)
			method.SetResult(NewFileDataTypeFromStr(result))
			method.SetException(NewFileDataTypeFromStr(except))
			if argsMp != nil {
				var i int16 = 0
				argIds := make(map[int16]int)
				for argName, argV := range argsMp {
					argMp, argOk := argV.(map[string]interface{})
					if !argOk {
						return fmt.Errorf("服务接口%s.%s输入参数%s定义错误，未定义id、类型等字段", svcName, methodName, argName)
					}
					argId := t.Int16(argMp["id"])
					szArgType := t.String(argMp["type"])
					argSummary := t.String(argMp["summary"])
					argLimit := t.String(argMp["limit"])
					argRule := t.String(argMp["rule"])
					i++
					if argId < 1 {
						argId = i
					}
					if _, ok := argIds[argId]; ok {
						return fmt.Errorf("服务接口%s.%s输入参数%sID定义重复%d", svcName, methodName, argName, argId)
					}
					if szArgType == "" {
						return fmt.Errorf("服务接口%s.%s输入参数%s未定义参数类型名", svcName, methodName, argName)
					}
					nLimit := SPDLimitRequired
					if argLimit == "optional" {
						nLimit = SPDLimitOptional
					}
					argType := NewFileDataTypeFromStr(szArgType)
					argField := NewFileDataField(argId, argName, argType, argSummary)
					argField.SetRule(argRule)
					argField.SetLimit(nLimit)
					method.AddArg(argField)
				}
			}
			if method.Args != nil {
				sort.Slice(method.Args, func(i, j int) bool {
					return method.Args[i].Id-method.Args[j].Id < 0
				})
			}
			service.AddMethod(method)
		}
		service.SetSummary(serviceSummary)
		p.AddService(service)
	}
	return nil
}

func (p *FileData) OpenFile(fileName string) error {
	if !xfile.Exists(fileName) {
		return fmt.Errorf("文件%s不存在", fileName)
	}
	parser, err := xparser.Load(fileName)
	if err != nil {
		return err
	}
	namespace := parser.GetString("namespace")
	if namespace == "" {
		return fmt.Errorf("文件%s不是有效的协议文件", fileName)
	}
	p.Namespace = namespace
	p.fileName = xfile.Basename(fileName)
	p.filePath = xfile.Dir(fileName)
	if arr := parser.GetArray("imports"); arr != nil {
		if err := p.loadImports(arr); err != nil {
			return err
		}
	}
	if m := parser.GetMap("typeDefs"); m != nil {
		if err := p.loadTypedefs(m); err != nil {
			return err
		}
	}
	if m := parser.GetMap("types"); m != nil {
		if err := p.loadTypes(m); err != nil {
			return err
		}
	}
	if m := parser.GetMap("interfaces"); m != nil {
		if err := p.loadServices(m); err != nil {
			return err
		}
	}
	return nil
}
