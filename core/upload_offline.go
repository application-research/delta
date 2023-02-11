package core

import "time"

type UploadOfflineParam struct {
	File        []byte
	Filename    string
	Replication int
	Miner       string
	WalletID    string
}

type UploadOfflineCommpParam struct {
	UploadOnlineParam
	Duration time.Duration
	Size     int64
}

type UploadOfflineResult struct {
}

type UploadOfflineService struct {
	DeltaNode DeltaNode
}

func NewUploadOfflineService(dn DeltaNode) *UploadOfflineService {
	return &UploadOfflineService{
		DeltaNode: dn,
	}
}

func (u *UploadOfflineService) Add(param UploadOfflineParam) (UploadOfflineResult, error) {
	return UploadOfflineResult{}, nil
}

func (u *UploadOfflineService) List(param UploadOfflineParam) (UploadOfflineResult, error) {
	return UploadOfflineResult{}, nil
}

func (u *UploadOfflineService) Commp(param UploadOfflineCommpParam) (UploadOfflineResult, error) {
	return UploadOfflineResult{}, nil
}
