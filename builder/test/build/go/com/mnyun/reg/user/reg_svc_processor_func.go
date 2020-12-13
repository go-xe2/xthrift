package user

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
)

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

		success := p.handler.RegDetailResult(input.Id)
		result.Success = success

		return result, nil
	}
	return nil, thrift.NewTApplicationException(thrift.INVALID_DATA, "输入参数错误")
}

func (p *RegSvcRegDetailResult) GetInputArgsInstance() thrift.TStruct {
	return NewRegSvcRegDetailResultArgs()

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
