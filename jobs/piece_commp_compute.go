package jobs

import (
	"context"
	"delta/core"
	"delta/utils"
	"io"
	"time"

	model "github.com/application-research/delta-db/db_models"
	"github.com/application-research/filclient"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	"github.com/labstack/gommon/log"
)

// PieceCommpProcessor `PieceCommpProcessor` is a struct that contains a `context.Context`, a `*core.DeltaNode`, a `model.Content`, a
// `filclient.DealConfig`, and a `*core.CommpService`.
// @property Context - The context of the current request
// @property LightNode - The light node that is currently being processed
// @property Content - The content of the piece of data
// @property DealPieceConfig - The configuration of the piece of information, including the piece of information, the piece
// of information, the piece of information, the piece of information, the piece of information, the piece of information,
// the piece of information, the piece of information, the piece of information, the piece of information, the
// @property CommpService - The CommpService object is used to communicate with the Commp protocol.
type PieceCommpProcessor struct {
	Context         context.Context
	LightNode       *core.DeltaNode
	Content         model.Content
	DealPieceConfig filclient.DealConfig
	CommpService    *core.CommpService
}

// NewPieceCommpProcessor `NewPieceCommpProcessor` is a function that returns a `PieceCommpProcessor` struct
func NewPieceCommpProcessor(ln *core.DeltaNode, content model.Content) IProcessor {
	commpService := new(core.CommpService)
	return &PieceCommpProcessor{
		LightNode:    ln,
		Content:      content,
		Context:      context.Background(),
		CommpService: commpService,
	}
}

// Run The process of generating the commp.
func (i PieceCommpProcessor) Run() error {
	var content model.Content

	i.LightNode.DB.Model(&i.Content).Where("id = ?", i.Content.ID).Find(&content)
	content.Status = utils.CONTENT_PIECE_COMPUTING
	content.UpdatedAt = time.Now()
	i.LightNode.DB.Save(&content)

	payloadCid, err := cid.Decode(i.Content.Cid)
	if err != nil {
		content.Status = utils.CONTENT_PIECE_COMPUTING_FAILED
		content.LastMessage = err.Error()
		content.UpdatedAt = time.Now()
		i.LightNode.DB.Save(&content)
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

	if i.LightNode.Config.Common.CommpMode == utils.COMMP_MODE_FAST {

		pieceInfo, err := i.CommpService.GenerateCommp(node)
		if err != nil {
			i.LightNode.DB.Model(&i.Content).Where("id = ?", i.Content.ID).Updates(model.Content{
				Status:      utils.CONTENT_FAILED_TO_PROCESS,
				LastMessage: err.Error(),
				UpdatedAt:   time.Now(),
			})
			return err
		}
		pieceCid = pieceInfo.PieceCID
		paddedPieceSize = abi.PaddedPieceSize(pieceInfo.PayloadSize)
		unPaddedPieceSize = pieceInfo.PieceSize.Unpadded()

		payloadSize = uint64(len(bytesFromCar))

	} else {

		if i.Content.ConnectionMode == utils.CONNECTION_MODE_IMPORT {
			pieceCid, payloadSize, unPaddedPieceSize, err = filclient.GeneratePieceCommitment(context.Background(), payloadCid, i.LightNode.Node.Blockstore)
			if err != nil {
				i.LightNode.DB.Model(&i.Content).Where("id = ?", i.Content.ID).Updates(model.Content{
					Status:      utils.CONTENT_FAILED_TO_PROCESS,
					LastMessage: err.Error(),
					UpdatedAt:   time.Now(),
				})
				return err
			}
			paddedPieceSize = abi.PaddedPieceSize(payloadSize)

		} else {
			pieceCid, payloadSize, unPaddedPieceSize, err = i.CommpService.GenerateCommPFile(i.Context, payloadCid, i.LightNode.Node.Blockstore)
			paddedPieceSize = unPaddedPieceSize.Padded()
			if err != nil {
				i.LightNode.DB.Model(&i.Content).Where("id = ?", i.Content.ID).Updates(model.Content{
					Status:      utils.CONTENT_FAILED_TO_PROCESS,
					LastMessage: err.Error(),
					UpdatedAt:   time.Now(),
				})
				return err
			}
		}
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
	i.LightNode.DB.Model(&i.Content).Where("id = ?", i.Content.ID).Updates(model.Content{
		Status:            utils.CONTENT_PIECE_ASSIGNED,
		PieceCommitmentId: commpRec.ID,
		UpdatedAt:         time.Now(),
	})

	// add this to the job queue
	item := NewStorageDealMakerProcessor(i.LightNode, i.Content, *commpRec)
	i.LightNode.Dispatcher.AddJobAndDispatch(item, 1)

	return nil
}
