package core

type DealStatusParam struct {
	DealUuid  string
	MinerAddr string
}

type DealStatusResult struct {
}

type DealStatusService struct {
	DeltaNode DeltaNode
	BoostApi  string
}

func NewDealStatusService(dn DeltaNode) *DealStatusService {
	return &DealStatusService{
		DeltaNode: dn,
	}
}

func (d DealStatusService) GetDealStatus(param DealStatusParam) (DealStatusResult, error) {
	return DealStatusResult{}, nil
}

// get dealid
// check status
