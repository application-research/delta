package core

import (
	"github.com/filecoin-project/go-data-segment/datasegment"
	"github.com/filecoin-project/go-data-segment/merkletree"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
)

type InclusionProofService struct {
	DeltaNode *DeltaNode
}

func NewInclusionProofService(deltaNode DeltaNode) *InclusionProofService {
	return &InclusionProofService{
		DeltaNode: &deltaNode,
	}
}

func (i *InclusionProofService) GenerateInclusionProof(pieceInfo abi.PieceInfo, paddedPieceSize abi.PaddedPieceSize) (*datasegment.InclusionProof, error) {
	aggregate, err := datasegment.NewAggregate(paddedPieceSize, []abi.PieceInfo{pieceInfo})
	if err != nil {
		panic(err)
	}
	inclusionProof, err := aggregate.ProofForPieceInfo(pieceInfo)
	if err != nil {
		return nil, err
	}
	// store on the database
	return inclusionProof, nil
}

func (i *InclusionProofService) ValidateInclusionProof(cidFromDb string, pieceInfoInput abi.PieceInfo) {

	// load from database
	//
	inclusionProof := datasegment.InclusionProof{
		ProofSubtree: merkletree.ProofData{},
		ProofIndex:   merkletree.ProofData{},
	}

	// from database, generate the expected aux data
	pieceCid, _ := cid.Decode("piececid....")
	pieceInfoFromDB := abi.PieceInfo{
		PieceCID: pieceCid,
		Size:     0,
	}
	inclusionVerifiedData := datasegment.VerifierDataForPieceInfo(pieceInfoFromDB)
	aux, _ := inclusionProof.ComputeExpectedAuxData(inclusionVerifiedData)

	// create a custom aux data to validate against
	validateAgainst := datasegment.InclusionAuxData{CommPa: pieceInfoInput.PieceCID, SizePa: pieceInfoInput.Size}

	if *aux == validateAgainst {
		// verified nice!
	} else {
		// not verified!
	}
}
