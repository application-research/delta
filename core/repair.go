package core

type RepairParam struct {
}

type RepairResult struct {
}

type RepairService struct {
	DeltaNode DeltaNode
}

func NewRepairService(dn DeltaNode) *RepairService {
	return &RepairService{
		DeltaNode: dn,
	}
}

func (r RepairService) RecreateDeal(param RepairParam) (RepairResult, error) {
	return RepairResult{}, nil
}

func (r RepairService) RestartDataTransfer(param RepairParam) (RepairResult, error) {
	return RepairResult{}, nil
}

func (r RepairService) RecreateCommp(param RepairParam) (RepairResult, error) {
	return RepairResult{}, nil
}

func (r RepairService) RecreateContent(param RepairParam) (RepairResult, error) {
	return RepairResult{}, nil
}
