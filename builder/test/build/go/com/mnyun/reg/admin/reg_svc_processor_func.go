package admin

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
)

type RegSvcRemoveResult struct {
	*xthrift.TBaseProcessorFunction
	handler RegSvc
}

func newRegSvcRemoveResult(handler RegSvc) *RegSvcRemoveResult {
	inst := &RegSvcRemoveResult{handler: handler}
	inst.TBaseProcessorFunction = xthrift.NewBaseProcessorFunction(inst)
	return inst
}

func (p *RegSvcRemoveResult) Invoke(args thrift.TStruct) (thrift.TStruct, error) {
	if input, ok := args.(*RegSvcRemoveResultArgs); ok {
		result := NewRegSvcRemoveResultResult()

		success := p.handler.RemoveResult(input.RegId)
		result.Success = &success

		return result, nil
	}
	return nil, thrift.NewTApplicationException(thrift.INVALID_DATA, "输入参数错误")
}

func (p *RegSvcRemoveResult) GetInputArgsInstance() thrift.TStruct {
	return NewRegSvcRemoveResultArgs()

}

type RegSvcUpdateResult struct {
	*xthrift.TBaseProcessorFunction
	handler RegSvc
}

func newRegSvcUpdateResult(handler RegSvc) *RegSvcUpdateResult {
	inst := &RegSvcUpdateResult{handler: handler}
	inst.TBaseProcessorFunction = xthrift.NewBaseProcessorFunction(inst)
	return inst
}

func (p *RegSvcUpdateResult) Invoke(args thrift.TStruct) (thrift.TStruct, error) {
	if input, ok := args.(*RegSvcUpdateResultArgs); ok {
		result := NewRegSvcUpdateResultResult()

		success := p.handler.UpdateResult(input.RegId, input.ParId, input.Name)
		result.Success = &success

		return result, nil
	}
	return nil, thrift.NewTApplicationException(thrift.INVALID_DATA, "输入参数错误")
}

func (p *RegSvcUpdateResult) GetInputArgsInstance() thrift.TStruct {
	return NewRegSvcUpdateResultArgs()

}

type RegSvcGetTownList struct {
	*xthrift.TBaseProcessorFunction
	handler RegSvc
}

func newRegSvcGetTownList(handler RegSvc) *RegSvcGetTownList {
	inst := &RegSvcGetTownList{handler: handler}
	inst.TBaseProcessorFunction = xthrift.NewBaseProcessorFunction(inst)
	return inst
}

func (p *RegSvcGetTownList) Invoke(args thrift.TStruct) (thrift.TStruct, error) {
	if input, ok := args.(*RegSvcGetTownListArgs); ok {
		result := NewRegSvcGetTownListResult()

		success := p.handler.GetTownList(input.CountyId)
		result.Success = success

		return result, nil
	}
	return nil, thrift.NewTApplicationException(thrift.INVALID_DATA, "输入参数错误")
}

func (p *RegSvcGetTownList) GetInputArgsInstance() thrift.TStruct {
	return NewRegSvcGetTownListArgs()

}

type RegSvcGetCityList struct {
	*xthrift.TBaseProcessorFunction
	handler RegSvc
}

func newRegSvcGetCityList(handler RegSvc) *RegSvcGetCityList {
	inst := &RegSvcGetCityList{handler: handler}
	inst.TBaseProcessorFunction = xthrift.NewBaseProcessorFunction(inst)
	return inst
}

func (p *RegSvcGetCityList) Invoke(args thrift.TStruct) (thrift.TStruct, error) {
	if input, ok := args.(*RegSvcGetCityListArgs); ok {
		result := NewRegSvcGetCityListResult()

		success := p.handler.GetCityList(input.ProvinceId)
		result.Success = success

		return result, nil
	}
	return nil, thrift.NewTApplicationException(thrift.INVALID_DATA, "输入参数错误")
}

func (p *RegSvcGetCityList) GetInputArgsInstance() thrift.TStruct {
	return NewRegSvcGetCityListArgs()

}

type RegSvcGetRegTreeResult struct {
	*xthrift.TBaseProcessorFunction
	handler RegSvc
}

func newRegSvcGetRegTreeResult(handler RegSvc) *RegSvcGetRegTreeResult {
	inst := &RegSvcGetRegTreeResult{handler: handler}
	inst.TBaseProcessorFunction = xthrift.NewBaseProcessorFunction(inst)
	return inst
}

func (p *RegSvcGetRegTreeResult) Invoke(args thrift.TStruct) (thrift.TStruct, error) {
	if input, ok := args.(*RegSvcGetRegTreeResultArgs); ok {
		result := NewRegSvcGetRegTreeResultResult()

		success := p.handler.GetRegTreeResult(input.ParId)
		result.Success = success

		return result, nil
	}
	return nil, thrift.NewTApplicationException(thrift.INVALID_DATA, "输入参数错误")
}

func (p *RegSvcGetRegTreeResult) GetInputArgsInstance() thrift.TStruct {
	return NewRegSvcGetRegTreeResultArgs()

}

type RegSvcRegDetailResult struct {
	*xthrift.TBaseProcessorFunction
	handler RegSvc
}

func newRegSvcRegDetailResult(handler RegSvc) *RegSvcRegDetailResult {
	inst := &RegSvcRegDetailResult{handler: handler}
	inst.TBaseProcessorFunction = xthrift.NewBaseProcessorFunction(inst)
	return inst
}

func (p *RegSvcRegDetailResult) Invoke(args thrift.TStruct) (thrift.TStruct, error) {
	if input, ok := args.(*RegSvcRegDetailResultArgs); ok {
		result := NewRegSvcRegDetailResultResult()

		success := p.handler.RegDetailResult(input.RegId)
		result.Success = success

		return result, nil
	}
	return nil, thrift.NewTApplicationException(thrift.INVALID_DATA, "输入参数错误")
}

func (p *RegSvcRegDetailResult) GetInputArgsInstance() thrift.TStruct {
	return NewRegSvcRegDetailResultArgs()

}

type RegSvcGetProvincesList struct {
	*xthrift.TBaseProcessorFunction
	handler RegSvc
}

func newRegSvcGetProvincesList(handler RegSvc) *RegSvcGetProvincesList {
	inst := &RegSvcGetProvincesList{handler: handler}
	inst.TBaseProcessorFunction = xthrift.NewBaseProcessorFunction(inst)
	return inst
}

func (p *RegSvcGetProvincesList) Invoke(args thrift.TStruct) (thrift.TStruct, error) {
	if _, ok := args.(*RegSvcGetProvincesListArgs); ok {
		result := NewRegSvcGetProvincesListResult()

		success := p.handler.GetProvincesList()
		result.Success = success

		return result, nil
	}
	return nil, thrift.NewTApplicationException(thrift.INVALID_DATA, "输入参数错误")
}

func (p *RegSvcGetProvincesList) GetInputArgsInstance() thrift.TStruct {
	return NewRegSvcGetProvincesListArgs()

}

type RegSvcGetCountList struct {
	*xthrift.TBaseProcessorFunction
	handler RegSvc
}

func newRegSvcGetCountList(handler RegSvc) *RegSvcGetCountList {
	inst := &RegSvcGetCountList{handler: handler}
	inst.TBaseProcessorFunction = xthrift.NewBaseProcessorFunction(inst)
	return inst
}

func (p *RegSvcGetCountList) Invoke(args thrift.TStruct) (thrift.TStruct, error) {
	if input, ok := args.(*RegSvcGetCountListArgs); ok {
		result := NewRegSvcGetCountListResult()

		success := p.handler.GetCountList(input.CityId)
		result.Success = success

		return result, nil
	}
	return nil, thrift.NewTApplicationException(thrift.INVALID_DATA, "输入参数错误")
}

func (p *RegSvcGetCountList) GetInputArgsInstance() thrift.TStruct {
	return NewRegSvcGetCountListArgs()

}

type RegSvcAddResult struct {
	*xthrift.TBaseProcessorFunction
	handler RegSvc
}

func newRegSvcAddResult(handler RegSvc) *RegSvcAddResult {
	inst := &RegSvcAddResult{handler: handler}
	inst.TBaseProcessorFunction = xthrift.NewBaseProcessorFunction(inst)
	return inst
}

func (p *RegSvcAddResult) Invoke(args thrift.TStruct) (thrift.TStruct, error) {
	if input, ok := args.(*RegSvcAddResultArgs); ok {
		result := NewRegSvcAddResultResult()

		success := p.handler.AddResult(input.ParId, input.Name)
		result.Success = &success

		return result, nil
	}
	return nil, thrift.NewTApplicationException(thrift.INVALID_DATA, "输入参数错误")
}

func (p *RegSvcAddResult) GetInputArgsInstance() thrift.TStruct {
	return NewRegSvcAddResultArgs()

}

type RegSvcGetChildListResult struct {
	*xthrift.TBaseProcessorFunction
	handler RegSvc
}

func newRegSvcGetChildListResult(handler RegSvc) *RegSvcGetChildListResult {
	inst := &RegSvcGetChildListResult{handler: handler}
	inst.TBaseProcessorFunction = xthrift.NewBaseProcessorFunction(inst)
	return inst
}

func (p *RegSvcGetChildListResult) Invoke(args thrift.TStruct) (thrift.TStruct, error) {
	if input, ok := args.(*RegSvcGetChildListResultArgs); ok {
		result := NewRegSvcGetChildListResultResult()

		success := p.handler.GetChildListResult(input.ParId)
		result.Success = success

		return result, nil
	}
	return nil, thrift.NewTApplicationException(thrift.INVALID_DATA, "输入参数错误")
}

func (p *RegSvcGetChildListResult) GetInputArgsInstance() thrift.TStruct {
	return NewRegSvcGetChildListResultArgs()

}
