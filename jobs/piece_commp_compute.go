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

type ItemContentProcessor struct {
	ContentProcessor
}

func NewItemContentProcessor(ln *core.LightNode, content core.Content) IProcessor {
	return &ItemContentProcessor{
		ContentProcessor{
			LightNode: ln,
			Content:   content,
			Context:   context.Background(),
		},
	}
}

func (i ItemContentProcessor) Run() error {

	i.LightNode.DB.Model(&core.Content{}).Where("id = ?", i.Content.ID).Updates(core.Content{Status: "piece-computing"})
	payloadCid, err := cid.Decode(i.Content.Cid)
	if err != nil {
		panic(err)
	}

	// prepare the commp
	commitment, u, a, err := filclient.GeneratePieceCommitmentFFI(i.Context, payloadCid, i.LightNode.Node.Blockstore)
	if err != nil {
		return err
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
	i.LightNode.DB.Create(commpRec)

	// update bucket status to commp-computed
	i.LightNode.DB.Model(&core.Content{}).Where("id = ?", i.Content.ID).Updates(core.Content{Status: "piece-assigned", PieceCommitmentId: commpRec.ID})
	return nil
}

func NewPieceCommpProcessor(ln *core.LightNode) IProcessor {
	return &PieceCommpProcessor{
		Processor{
			LightNode: ln,
		},
	}
}

func (r *PieceCommpProcessor) Run() error {
	// get the CID field of the bucket and generate a commp for it.
	var contents []core.Content
	r.LightNode.DB.Model(&core.Content{}).Where("status = ?", "pinned").Find(&contents)
	dispatcher := CreateNewDispatcher()
	for _, content := range contents {
		job := NewItemContentProcessor(r.LightNode, content)
		dispatcher.AddJob(job)
		dispatcher.Start(10)
	}
	return nil
}
