package jobs

import (
	"context"
	"fc-deal-making-service/core"
	"github.com/application-research/filclient"
	"github.com/ipfs/go-cid"
	"time"
)

// workers
// jobs

// this processors are independent. we want it to run on it's own without waiting
// for other groups.

type PieceCommpProcessor struct {
	Processor
}

func NewPieceCommpProcessor(ln *core.LightNode) PieceCommpProcessor {
	return PieceCommpProcessor{
		Processor{
			LightNode: ln,
		},
	}
}

func (r *PieceCommpProcessor) Run() {

	// get the CID field of the bucket and generate a commp for it.
	var contents []core.Content
	r.LightNode.DB.Model(&core.Content{}).Where("status = ?", "open").Find(&contents)

	// for each bucket, generate a commp
	for _, content := range contents {

		payloadCid, err := cid.Decode(content.Cid)
		if err != nil {
			panic(err)
		}

		// prepare the commp
		commitment, u, a, err := filclient.GeneratePieceCommitmentFFI(context.Background(), payloadCid, r.LightNode.Node.Blockstore)
		if err != nil {
			return
		}

		// save the commp to the database
		commpRec := &core.PieceCommitment{
			Cid:             payloadCid.String(),
			Piece:           commitment.String(),
			Size:            int64(u),
			PaddedPieceSize: uint64(a),
			Status:          "open",
			Created_at:      time.Now(),
			Updated_at:      time.Now(),
		}
		r.LightNode.DB.Create(commpRec)

		// update bucket status to commp-computed
		r.LightNode.DB.Model(&core.Content{}).Where("id = ?", content.ID).Update("piece_commitment_id", commpRec.ID).Update("status", "piece-assigned")
	}
}
