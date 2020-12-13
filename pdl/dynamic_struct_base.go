/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-16 11:03
* Description:
*****************************************************************/

package pdl

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/x/xf/ef/xqi"
)

type TDynamicStructBase struct {
	inst        DynamicStruct
	dsFieldMaps map[string]int
}

var _ DynamicStruct = (*TDynamicStructBase)(nil)

func NewBasicStruct(inst DynamicStruct) *TDynamicStructBase {
	return &TDynamicStructBase{
		inst: inst,
	}
}

func (p *TDynamicStructBase) NewInstance() DynamicStruct {
	return nil
}

func (p *TDynamicStructBase) AllFields() map[string]*TStructFieldInfo {
	return nil
}

func (p *TDynamicStructBase) FieldNameMaps() map[string]string {
	return nil
}

func (p *TDynamicStructBase) SetFieldValue(fdName string, val interface{}) bool {
	fields := p.inst.AllFields()
	if fields == nil {
		return false
	}
	if fn, ok := fields[fdName]; ok {
		return fn.Setter(p.inst, val)
	}
	return false
}

func (p *TDynamicStructBase) AssignFromMap(mp map[string]interface{}) bool {
	fields := p.inst.FieldNameMaps()
	if fields == nil {
		return false
	}
	for a, k := range fields {
		if v, ok := mp[k]; ok {
			p.inst.SetFieldValue(k, v)
		} else if v, ok := mp[a]; ok {
			p.inst.SetFieldValue(k, v)
		}
	}
	return true
}

func (p *TDynamicStructBase) SliceFromMaps(mps []map[string]interface{}) []thrift.TStruct {
	size := len(mps)
	result := make([]thrift.TStruct, 0)
	for i := 0; i < size; i++ {
		inst := p.inst.NewInstance()
		inst.AssignFromMap(mps[i])
		if d, ok := inst.(thrift.TStruct); ok {
			result = append(result, d)
		}
	}
	return result
}

func (p *TDynamicStructBase) makeDatasetFieldMaps(ds xqi.Dataset) {
	if p.dsFieldMaps != nil {
		return
	}
	p.dsFieldMaps = make(map[string]int)
	fields := p.inst.FieldNameMaps()
	if fields == nil || ds == nil || !ds.IsOpen() {
		return
	}
	for a, k := range fields {
		if fd := ds.DSFieldByName(k); fd != nil {
			p.dsFieldMaps[k] = fd.FieldIndex()
		} else if fd := ds.DSFieldByName(a); fd != nil {
			p.dsFieldMaps[k] = fd.FieldIndex()
		}
	}
}

func (p *TDynamicStructBase) AssignFromDataSet(ds xqi.Dataset) bool {
	if ds == nil || !ds.IsOpen() {
		return false
	}
	p.makeDatasetFieldMaps(ds)
	fields := p.dsFieldMaps
	for k, idx := range fields {
		v := ds.FieldValue(idx)
		if v != nil {
			p.inst.SetFieldValue(k, v)
		}
	}
	return true
}

func (p *TDynamicStructBase) SliceFromDataSet(ds xqi.Dataset) []thrift.TStruct {
	result := make([]thrift.TStruct, 0)
	if ds == nil || !ds.IsOpen() {
		return result
	}
	p.makeDatasetFieldMaps(ds)
	defer func() {
		p.dsFieldMaps = nil
	}()
	ds.MoveFirst()
	for ds.Next() {
		inst := p.inst.NewInstance()
		if inst.AssignFromDataSet(ds) {
			if d, ok := inst.(thrift.TStruct); ok {
				result = append(result, d)
			}
		}
	}
	return result
}
