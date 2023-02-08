package core

type MinerReputation struct {
}

func NewMinerReputation() *MinerReputation {
	return &MinerReputation{}
}

type ReputableMinerParam struct {
	MaxPieceSize int64
	MinPieceSize int64
	MaxPrice     int64
	MinPrice     int64
	MaxDuration  int64
}

func (mr *MinerReputation) GetReputableMinerBasedOnParam(minerReputableParam ReputableMinerParam) int {

	return 0
}
