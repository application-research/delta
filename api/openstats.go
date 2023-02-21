package api

import (
	"delta/core"
	"delta/core/model"
	"github.com/labstack/echo/v4"
)

func ConfigureOpenStatsCheckRouter(e *echo.Group, node *core.DeltaNode) {

	e.GET("/stats/miner/:minerId", func(c echo.Context) error {
		return handleOpenStatsByMiner(c, node)
	})

	e.GET("/stats/miner/:minerId/deals", func(c echo.Context) error {
		return handleOpenGetDealsByMiner(c, node)
	})

	e.GET("/stats/totals/info", func(c echo.Context) error {
		return handleOpenGetTotalsInfo(c, node)
	})
}

func handleOpenStatsByMiner(c echo.Context, node *core.DeltaNode) error {

	// get content consumed by miner
	var content []model.Content
	node.DB.Raw("select c.* from contents c, content_miner cma where c.id = cma.content and cma.miner = ?", c.Param("minerId")).Scan(&content)

	// get content deals by miner
	var contentDeal []model.ContentDeal
	node.DB.Raw("select cd.* from content_deals cd, content_miners cma where cd.content = cma.content and cma.miner = ?", c.Param("minerId")).Scan(&contentDeal)

	var pieceCommitments []model.PieceCommitment
	node.DB.Raw("select pc.* from piece_commitments pc, content_deals cd, content_miners cma where pc.content_deal = cd.id and cd.content = cma.content and cma.miner = ?", c.Param("minerId")).Scan(&pieceCommitments)

	var contentDealProposal []model.ContentDealProposalParameters
	node.DB.Raw("select cdp.* from content_deal_proposal_parameters cdp, content_deals cd, content_miners cma where cdp.content_deal = cd.id and cd.content = cma.content and cma.miner = ?", c.Param("minerId")).Scan(&contentDealProposal)

	return c.JSON(200, map[string]interface{}{
		"content":          content,
		"deals":            contentDeal,
		"commitment_piece": pieceCommitments,
		"deal_proposals":   contentDealProposal,
	})

	return nil
}

// function to get all deals of a given miner
func handleOpenGetDealsByMiner(c echo.Context, node *core.DeltaNode) error {

	var contentDeal []model.ContentDeal
	node.DB.Raw("select cd.* from content_deals cd, contents c where cd.content = c.id and c.miner = ?", c.Param("minerId")).Scan(&contentDeal)

	c.JSON(200, map[string]interface{}{
		"deals": contentDeal,
	})

	return nil
}

// function to get all totals info
func handleOpenGetTotalsInfo(c echo.Context, node *core.DeltaNode) error {

	var totalContentConsumed int64
	node.DB.Raw("select count(*) from contents").Scan(&totalContentConsumed)

	var totalTransferStarted int64
	node.DB.Raw("select count(*) from contents where status = 'transfer-started'").Scan(&totalTransferStarted)

	var totalTransferFinished int64
	node.DB.Raw("select count(*) from contents where status = 'transfer-finished'").Scan(&totalTransferFinished)

	var totalProposalMade int64
	node.DB.Raw("select count(*) from deal_proposals").Scan(&totalProposalMade)

	var totalPiece int64
	node.DB.Raw("select count(*) from piece_commitments").Scan(&totalPiece)

	var totalPieceCommitted int64
	node.DB.Raw("select count(*) from piece_commitments where status = 'committed'").Scan(&totalPieceCommitted)

	var totalMiners int64
	node.DB.Raw("select distinct(miner) from content_miners").Count(&totalMiners)

	var totalStorageAllocated int64
	node.DB.Raw("select sum(size) from contents").Scan(&totalStorageAllocated)

	var totalProposalSent int64
	node.DB.Raw("select count(*) from contents where status = 'deal-proposal-sent'").Scan(&totalProposalSent)

	c.JSON(200, map[string]interface{}{
		"total_content_consumed":  totalContentConsumed,
		"total_transfer_started":  totalTransferStarted,
		"total_transfer_finished": totalTransferFinished,
		"total_proposal_made":     totalProposalMade,
		"total_piece":             totalPiece,
		"total_piece_committed":   totalPieceCommitted,
		"total_miners":            totalMiners,
		"total_storage_allocated": totalStorageAllocated,
		"total_proposal_sent":     totalProposalSent,
	})
	return nil
}
