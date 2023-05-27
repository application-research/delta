package jobs

import (
	"context"
	"delta/core"
	model "github.com/application-research/delta-db/db_models"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
)

type DealStatusCheck struct {
	LightNode *core.DeltaNode
	Content   *model.Content
}

func (d DealStatusCheck) Run() error {
	var contentDeal model.ContentDeal
	// get the latest content deal of the content
	// select * from content_deals where content = content.id order by created_at desc limit 1
	d.LightNode.DB.Raw("select * from content_deals where content = ? order by created_at desc limit 1", d.Content.ID).Scan(&contentDeal)

	if contentDeal.DealUUID == "" {
		return nil
	}

	miner, err := address.NewFromString(contentDeal.Miner)
	if err != nil {
		return err
	}

	cidProp, err := cid.Decode(contentDeal.PropCid)
	if err != nil {
		return err

	}
	dealUuid, err := uuid.Parse(contentDeal.DealUUID)
	if err != nil {
		return err
	}

	// get the status
	status, err := d.LightNode.FilClient.DealStatus(context.Background(), miner, cidProp, &dealUuid)
	contentDeal.DealID = int64(status.DealID)
	d.Content.Status = storagemarket.DealStates[status.State]
	d.Content.LastMessage = storagemarket.DealStatesDescriptions[status.State]
	contentDeal.LastMessage = storagemarket.DealStatesDescriptions[status.State]
	d.LightNode.DB.Save(&contentDeal)
	d.LightNode.DB.Save(&d.Content)
	return nil
}

func NewDealStatusCheck(ln *core.DeltaNode, content *model.Content) IProcessor {
	return &DealStatusCheck{
		LightNode: ln,
		Content:   content,
	}
}
