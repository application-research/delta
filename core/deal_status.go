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

//func (d DealStatusService) GetDealID(param DealStatusParam) (DealStatusResult, error) {
//	ctx := context.Background()
//	d.DeltaNode.LotusApi.StateSearchMsg(ctx, types.EmptyTSK, pubcid, 1000, false)
//	return DealStatusResult{}, nil
//}

//func (m *manager) GetProviderDealStatus(ctx context.Context, d *model.ContentDeal, maddr address.Address, dealUUID *uuid.UUID) (*storagemarket.ProviderDealState, bool, error) {
//	isPushTransfer := false
//	providerDealState, err := m.fc.DealStatus(ctx, maddr, d.PropCid.CID, dealUUID)
//	if err != nil && providerDealState == nil {
//		isPushTransfer = true
//		providerDealState, err = m.fc.DealStatus(ctx, maddr, d.PropCid.CID, nil)
//	}
//	return providerDealState, isPushTransfer, err
//}

//func (s *apiV1) handleGetDealInfo(c echo.Context) error {
//	dealid, err := strconv.ParseInt(c.Param("dealid"), 10, 64)
//	if err != nil {
//		return err
//	}
//
//	deal, err := s.api.StateMarketStorageDeal(c.Request().Context(), abi.DealID(dealid), types.EmptyTSK)
//	if err != nil {
//		return err
//	}
//
//	return c.JSON(http.StatusOK, deal)
//}
