package user

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-xe2/xthrift/builder/test/build/go/com/mnyun/reg/types"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
	"golang.org/x/net/context"
)

type RegSvcClient struct {
	*xthrift.TXClient
}

func NewRegSvcClient(trans thrift.TTransport, in, out thrift.TProtocolFactory) *RegSvcClient {
	inst := &RegSvcClient{
		TXClient: xthrift.NewClient(trans, in, out),
	}
	return inst
}

// 下级地区
func (p *RegSvcClient) GetChildListResult(cxt context.Context, parId *int32) ([]*types.RegItem, error) {
	var args = NewRegSvcGetChildListResultArgs()
	args.ParId = parId
	result := NewRegSvcGetChildListResultResult()
	err := p.Call(cxt, "GetChildListResult", args, result)
	if err != nil {
		return nil, err
	}
	return result.GetSuccess(), nil
}

// 州市列表
func (p *RegSvcClient) GetCityList(cxt context.Context, provinceId int32) ([]*types.RegItem, error) {
	var args = NewRegSvcGetCityListArgs()
	args.ProvinceId = provinceId
	result := NewRegSvcGetCityListResult()
	err := p.Call(cxt, "GetCityList", args, result)
	if err != nil {
		return nil, err
	}
	return result.GetSuccess(), nil
}

// 区县列表
func (p *RegSvcClient) GetCountList(cxt context.Context, cityId int32) ([]*types.RegItem, error) {
	var args = NewRegSvcGetCountListArgs()
	args.CityId = cityId
	result := NewRegSvcGetCountListResult()
	err := p.Call(cxt, "GetCountList", args, result)
	if err != nil {
		return nil, err
	}
	return result.GetSuccess(), nil
}

// 省份列表
func (p *RegSvcClient) GetProvincesList(cxt context.Context) ([]*types.RegItem, error) {
	var args = NewRegSvcGetProvincesListArgs()
	result := NewRegSvcGetProvincesListResult()
	err := p.Call(cxt, "GetProvincesList", args, result)
	if err != nil {
		return nil, err
	}
	return result.GetSuccess(), nil
}

// 地区目录树
func (p *RegSvcClient) GetRegTreeResult(cxt context.Context, parId *int32) ([]*types.RegItem, error) {
	var args = NewRegSvcGetRegTreeResultArgs()
	args.ParId = parId
	result := NewRegSvcGetRegTreeResultResult()
	err := p.Call(cxt, "GetRegTreeResult", args, result)
	if err != nil {
		return nil, err
	}
	return result.GetSuccess(), nil
}

// 乡镇列表
func (p *RegSvcClient) GetTownList(cxt context.Context, countyId int32) ([]*types.RegItem, error) {
	var args = NewRegSvcGetTownListArgs()
	args.CountyId = countyId
	result := NewRegSvcGetTownListResult()
	err := p.Call(cxt, "GetTownList", args, result)
	if err != nil {
		return nil, err
	}
	return result.GetSuccess(), nil
}

// 详情
func (p *RegSvcClient) RegDetailResult(cxt context.Context, id int32) (*types.RegItem, error) {
	var args = NewRegSvcRegDetailResultArgs()
	args.Id = id
	result := NewRegSvcRegDetailResultResult()
	err := p.Call(cxt, "RegDetailResult", args, result)
	if err != nil {
		return nil, err
	}
	return result.GetSuccess(), nil
}
