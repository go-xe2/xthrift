package admin

import (
	"github.com/go-xe2/xthrift/lib/go/xthrift"
)

type RegSvcProcessor struct {
	*xthrift.TBaseProcessor
	handler RegSvc
}

func NewRegSvcProcessor(handler RegSvc) *RegSvcProcessor {
	inst := &RegSvcProcessor{
		handler: handler,
	}

	inst.TBaseProcessor = xthrift.NewBaseProcessor(inst)
	return inst.registerFunctions()
}

func (p *RegSvcProcessor) registerFunctions() *RegSvcProcessor {
	p.RegisterFunction("AddResult", newRegSvcAddResult(p.handler))
	p.RegisterFunction("GetChildListResult", newRegSvcGetChildListResult(p.handler))
	p.RegisterFunction("GetCityList", newRegSvcGetCityList(p.handler))
	p.RegisterFunction("GetCountList", newRegSvcGetCountList(p.handler))
	p.RegisterFunction("GetProvincesList", newRegSvcGetProvincesList(p.handler))
	p.RegisterFunction("GetRegTreeResult", newRegSvcGetRegTreeResult(p.handler))
	p.RegisterFunction("GetTownList", newRegSvcGetTownList(p.handler))
	p.RegisterFunction("RegDetailResult", newRegSvcRegDetailResult(p.handler))
	p.RegisterFunction("RemoveResult", newRegSvcRemoveResult(p.handler))
	p.RegisterFunction("UpdateResult", newRegSvcUpdateResult(p.handler))
	return p
}
