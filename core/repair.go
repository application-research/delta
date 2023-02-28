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

// RecreateDeal A method of RepairService.
// A method of RepairService.
func (r RepairService) RecreateDeal(param RepairParam) (RepairResult, error) {
	return RepairResult{}, nil
}

// A method of RepairService.
func (r RepairService) RestartDataTransfer(param RepairParam) (RepairResult, error) {
	return RepairResult{}, nil
}

// A method of RepairService.
func (r RepairService) RecreateCommp(param RepairParam) (RepairResult, error) {
	return RepairResult{}, nil
}

// A method of RepairService.
func (r RepairService) RecreateContent(param RepairParam) (RepairResult, error) {
	return RepairResult{}, nil
}
