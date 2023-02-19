package jobs

import (
	"context"
	"delta/core"
	"delta/core/model"
	"delta/utils"
	"github.com/application-research/filclient"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	"github.com/labstack/gommon/log"
	"io"
	"time"
)

type PieceCommpProcessor struct {
	Context         context.Context
	LightNode       *core.DeltaNode
	Content         model.Content
	DealPieceConfig filclient.DealConfig
	CommpService    *core.CommpService
}

func NewPieceCommpProcessor(ln *core.DeltaNode, content model.Content) IProcessor {
	commpService := new(core.CommpService)
	return &PieceCommpProcessor{
		LightNode:    ln,
		Content:      content,
		Context:      context.Background(),
		CommpService: commpService,
	}
}

func (i PieceCommpProcessor) Run() error {

	i.LightNode.DB.Model(&model.Content{}).Where("id = ?", i.Content.ID).Updates(model.Content{Status: utils.CONTENT_PIECE_COMPUTING})
	payloadCid, err := cid.Decode(i.Content.Cid)
	if err != nil {
		i.LightNode.DB.Model(&model.Content{}).Where("id = ?", i.Content.ID).Updates(model.Content{Status: utils.CONTENT_PIECE_COMPUTING_FAILED, LastMessage: err.Error()})
	}

	// prepare the commp
	node, err := i.LightNode.Node.GetFile(context.Background(), payloadCid)
	nodeCopy := node

	bytesFromCar, err := io.ReadAll(nodeCopy)
	if err != nil {
		log.Error(err)
	}

	var pieceCid cid.Cid
	var payloadSize uint64
	var unPaddedPieceSize abi.UnpaddedPieceSize
	var paddedPieceSize abi.PaddedPieceSize

	pieceInfo, err := i.CommpService.GenerateCommPCarV2(node)
	if err != nil {
		pieceCid, payloadSize, unPaddedPieceSize, err = i.CommpService.GenerateCommPFile(i.Context, payloadCid, i.LightNode.Node.Blockstore)
		paddedPieceSize = unPaddedPieceSize.Padded()
		if err != nil {
			log.Error(err)
		}
	} else {

		pieceCid = pieceInfo.PieceCID
		paddedPieceSize = pieceInfo.Size
		unPaddedPieceSize = pieceInfo.Size.Unpadded()

		if err != nil {
			log.Error(err)
		}
		payloadSize = uint64(len(bytesFromCar))
	}

	if err != nil {
		// put this back to the queue
		i.LightNode.Dispatcher.AddJobAndDispatch(NewPieceCommpProcessor(i.LightNode, i.Content), 1)
		return err
	}

	// save the commp to the database
	commpRec := &model.PieceCommitment{
		Cid:               payloadCid.String(),
		Piece:             pieceCid.String(),
		Size:              int64(payloadSize),
		PaddedPieceSize:   uint64(paddedPieceSize),
		UnPaddedPieceSize: uint64(unPaddedPieceSize),
		Status:            "open",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	i.LightNode.DB.Create(commpRec)
	i.LightNode.DB.Model(&model.Content{}).Where("id = ?", i.Content.ID).Updates(model.Content{Status: utils.CONTENT_PIECE_ASSIGNED, PieceCommitmentId: commpRec.ID})

	// add this to the job queue
	item := NewStorageDealMakerProcessor(i.LightNode, i.Content, *commpRec)
	i.LightNode.Dispatcher.AddJobAndDispatch(item, 1)

	return nil
}
