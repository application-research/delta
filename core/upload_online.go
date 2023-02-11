package core

import "time"

type UploadOnlineParam struct {
	File        []byte
	Filename    string
	Replication int
	Miner       string
	WalletID    string
}

type UploadOnlineCommpParam struct {
	UploadOnlineParam
	Duration time.Duration
	Size     int64
}

type UploadOnlineResult struct {
}

type UploadOnlineService struct {
	DeltaNode DeltaNode
}

func NewUploadOnlineService(dn DeltaNode) *UploadOnlineService {
	return &UploadOnlineService{
		DeltaNode: dn,
	}
}

func (u UploadOnlineService) Add(param UploadOnlineParam) (UploadOnlineResult, error) {
	return UploadOnlineResult{}, nil
}

func (u UploadOnlineService) List(param UploadOnlineParam) (UploadOnlineResult, error) {
	return UploadOnlineResult{}, nil
}

func (u UploadOnlineService) Commp(param UploadOnlineCommpParam) (UploadOnlineResult, error) {
	return UploadOnlineResult{}, nil
}
