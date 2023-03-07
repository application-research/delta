package api

import (
	"delta/core"
	model "github.com/application-research/delta-db/db_models"
	"github.com/labstack/echo/v4"
)

// TODO: OPTIMIZE!!

// ConfigureOpenStatsCheckRouter It configures the router to handle the following routes:
//
// - `GET /stats/miner/:minerId`
// - `GET /stats/miner/:minerId/deals`
// - `GET /stats/totals/info`
//
// The first two routes are handled by the `handleOpenStatsByMiner` and `handleOpenGetDealsByMiner` functions,
// respectively. The last route is handled by the `handleOpenGetTotalsInfo` function
func ConfigureOpenStatsCheckRouter(e *echo.Group, node *core.DeltaNode) {

	e.GET("/stats/miner/:minerId", func(c echo.Context) error {
		return handleOpenStatsByMiner(c, node)
	})

	e.GET("/stats/miner/:minerId/deals", func(c echo.Context) error {
		return handleOpenGetDealsByMiner(c, node)
	})

	e.GET("/stats/content/:contentId", func(c echo.Context) error {
		return handleOpenGetStatsByContent(c, node)
	})
	e.GET("/stats/all-contents", func(c echo.Context) error {
		return handleOpenGetStatsByAllContents(c, node)
	})
	e.POST("/stats/contents", func(c echo.Context) error {
		return handleOpenGetStatsByContents(c, node)
	})

	e.GET("/stats/totals/info", func(c echo.Context) error {
		return handleOpenGetTotalsInfo(c, node)
	})

	e.GET("/stats/deal/by-cid/:cid", func(c echo.Context) error {
		return handleOpenGetDealByCid(c, node)
	})

	e.GET("/stats/deal/by-uuid/:uuid", func(c echo.Context) error {
		return handleOpenGetDealByUuid(c, node)
	})

	e.GET("/stats/deal/by-deal-id/:dealId", func(c echo.Context) error {
		return handleOpenGetDealByDealId(c, node)
	})
}

// It gets the content deal, content, content deal proposal, and piece commitment from the database and returns them as
// JSON
func handleOpenGetDealByCid(c echo.Context, node *core.DeltaNode) error {

	var content model.Content
	node.DB.Raw("select * from contents where cid = ?", c.Param("cid")).Scan(&content)
	content.RequestingApiKey = ""

	var contentDeal model.ContentDeal
	node.DB.Raw("select * from content_deals where content = ?", content.ID).Scan(&contentDeal)

	var contentDealProposal model.ContentDealProposal
	node.DB.Raw("select * from content_deal_proposals where content = ?", content.ID).Scan(&contentDealProposal)

	var pieceCommitment model.PieceCommitment
	node.DB.Raw("select * from piece_commitments where id = ?", content.PieceCommitmentId).Scan(&pieceCommitment)

	return c.JSON(200, map[string]interface{}{
		"deal":             contentDeal,
		"content":          content,
		"deal_proposal":    contentDealProposal,
		"piece_commitment": pieceCommitment,
	})
}

// It gets the content deal, content, content deal proposal, and piece commitment from the database and returns them as
// JSON
func handleOpenGetDealByDealId(c echo.Context, node *core.DeltaNode) error {
	var contentDeal model.ContentDeal
	node.DB.Raw("select * from content_deals where deal_id = ?", c.Param("dealId")).Scan(&contentDeal)

	var content model.Content
	node.DB.Raw("select * from contents where id = ?", contentDeal.Content).Scan(&content)
	content.RequestingApiKey = ""

	var contentDealProposal model.ContentDealProposal
	node.DB.Raw("select * from content_deal_proposals where content = ?", content.ID).Scan(&contentDealProposal)

	var pieceCommitment model.PieceCommitment
	node.DB.Raw("select * from piece_commitments where id = ?", content.PieceCommitmentId).Scan(&pieceCommitment)

	return c.JSON(200, map[string]interface{}{
		"deal":             contentDeal,
		"content":          content,
		"deal_proposal":    contentDealProposal,
		"piece_commitment": pieceCommitment,
	})
}

func handleOpenGetDealByUuid(c echo.Context, node *core.DeltaNode) error {
	var contentDeal model.ContentDeal
	node.DB.Raw("select * from content_deals where deal_uuid = ?", c.Param("uuid")).Scan(&contentDeal)

	var content model.Content
	node.DB.Raw("select * from contents where id = ?", contentDeal.Content).Scan(&content)
	content.RequestingApiKey = ""

	var contentDealProposal model.ContentDealProposal
	node.DB.Raw("select * from content_deal_proposals where content = ?", content.ID).Scan(&contentDealProposal)

	var pieceCommitment model.PieceCommitment
	node.DB.Raw("select * from piece_commitments where id = ?", content.PieceCommitmentId).Scan(&pieceCommitment)

	return c.JSON(200, map[string]interface{}{
		"deal":             contentDeal,
		"content":          content,
		"deal_proposal":    contentDealProposal,
		"piece_commitment": pieceCommitment,
	})
}

// Getting the content consumed by a miner, the content deals by a miner, the piece commitments by a miner, and the content
// deal proposals by a miner.
func handleOpenStatsByMiner(c echo.Context, node *core.DeltaNode) error {

	// get content consumed by miner
	var content []model.Content
	node.DB.Raw("select c.* from contents c, content_miner cma where c.id = cma.content and cma.miner = ?", c.Param("minerId")).Scan(&content)
	// TODO: remove api key
	for i := range content {
		content[i].RequestingApiKey = ""
	}

	// get content deals by miner
	var contentDeal []model.ContentDeal
	node.DB.Raw("select cd.* from content_deals cd, content_miners cma where cd.content = cma.content and cma.miner = ?", c.Param("minerId")).Scan(&contentDeal)

	var pieceCommitments []model.PieceCommitment
	node.DB.Raw("select pc.* from piece_commitments pc, content_deals cd, content_miners cma where pc.content_deal = cd.id and cd.content = cma.content and cma.miner = ?", c.Param("minerId")).Scan(&pieceCommitments)

	var contentDealProposal []model.ContentDealProposal
	node.DB.Raw("select cdp.* from content_deal_proposals cdp, content_deals cd, content_miners cma where cdp.content = cd.content and cd.content = cma.content and cma.miner = ?", c.Param("minerId")).Scan(&contentDealProposal)

	return c.JSON(200, map[string]interface{}{
		"content":           content,
		"deals":             contentDeal,
		"piece_commitments": pieceCommitments,
		"deal_proposals":    contentDealProposal,
	})

	return nil
}

// function to get all deals of a given miner
func handleOpenGetDealsByMiner(c echo.Context, node *core.DeltaNode) error {

	var contentDeal []model.ContentDeal
	node.DB.Raw("select cd.* from content_deals cd, contents c where cd.content = c.id and cd.miner = ?", c.Param("minerId")).Scan(&contentDeal)

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
	node.DB.Raw("select count(*) from content_deal_proposals").Scan(&totalProposalMade)

	var totalCommitmentPiece int64
	node.DB.Raw("select count(*) from piece_commitments").Scan(&totalCommitmentPiece)

	var totalPieceCommitted int64
	node.DB.Raw("select count(*) from piece_commitments where status = 'committed'").Scan(&totalPieceCommitted)

	var totalMiners int64
	rows, err := node.DB.Raw("select distinct(miner) from content_miners").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var miner string
		rows.Scan(&miner)
		totalMiners++
	}

	var totalStorageAllocated int64
	node.DB.Raw("select sum(size) from contents").Scan(&totalStorageAllocated)

	var totalProposalSent int64
	node.DB.Raw("select count(*) from contents where status = 'deal-proposal-sent'").Scan(&totalProposalSent)

	var totalSealedDealInBytes int64
	node.DB.Raw("select sum(size) from contents where status in ('transfer-started','transfer-finished','deal-proposal-sent')").Scan(&totalSealedDealInBytes)

	var totalImportDeals int64
	node.DB.Raw("select count(*) from contents where connection_mode = 'import'").Scan(&totalImportDeals)

	var totalE2EDeals int64
	node.DB.Raw("select count(*) from contents where connection_mode = 'e2e'").Scan(&totalE2EDeals)

	var totalE2EDealsInBytes int64
	node.DB.Raw("select sum(size) from contents where connection_mode = 'e2e'").Scan(&totalE2EDealsInBytes)

	var totalImportDealsInBytes int64
	node.DB.Raw("select sum(size) from contents where connection_mode = 'import'").Scan(&totalImportDealsInBytes)

	c.JSON(200, map[string]interface{}{
		"total_content_consumed":      totalContentConsumed,
		"total_transfer_started":      totalTransferStarted,
		"total_transfer_finished":     totalTransferFinished,
		"total_piece_commitment_made": totalCommitmentPiece,
		"total_piece_committed":       totalPieceCommitted,
		"total_miners":                totalMiners,
		"total_storage_allocated":     totalStorageAllocated,
		"total_proposal_made":         totalProposalMade,
		"total_proposal_sent":         totalProposalSent,
		"total_sealed_deal_in_bytes":  totalSealedDealInBytes,
		"total_import_deals":          totalImportDeals,
		"total_e2e_deals":             totalE2EDeals,
		"total_e2e_deals_in_bytes":    totalE2EDealsInBytes,
		"total_import_deals_in_bytes": totalImportDealsInBytes,
	})
	return nil
}

func handleOpenGetStatsByAllContents(c echo.Context, node *core.DeltaNode) error {

	var contentIds []int64
	node.DB.Raw("select id from contents").Scan(&contentIds)

	var contentResponse []map[string]interface{}
	for _, contentId := range contentIds {

		var content model.Content
		node.DB.Raw("select c.* from contents c where c.id = ?", contentId).Scan(&content)
		content.RequestingApiKey = ""

		var contentDeal []model.ContentDeal
		node.DB.Raw("select cd.* from content_deals cd, contents c where cd.content = c.id and c.id = ?", contentId).Scan(&contentDeal)

		var pieceCommitments []model.PieceCommitment
		node.DB.Raw("select pc.* from piece_commitments pc, contents c where c.piece_commitment_id = pc.id and c.id = ?", contentId).Scan(&pieceCommitments)

		var contentDealProposal []model.ContentDealProposal
		node.DB.Raw("select cdp.* from content_deal_proposals cdp, contents c where cdp.content = c.id and c.id = ?", contentId).Scan(&contentDealProposal)

		contentResponse = append(contentResponse, map[string]interface{}{
			"content":           content,
			"deals":             contentDeal,
			"piece_commitments": pieceCommitments,
			"deal_proposals":    contentDealProposal,
		})
	}
	return c.JSON(200, contentResponse)

}

// A function that is called when a GET request is made to the /open/get_stats_by_contents endpoint.
func handleOpenGetStatsByContents(c echo.Context, node *core.DeltaNode) error {

	var contentIds []int64
	c.Bind(&contentIds)

	var contentResponse []map[string]interface{}
	for _, contentId := range contentIds {

		var content model.Content
		node.DB.Raw("select c.* from contents c where c.id = ?", contentId).Scan(&content)
		content.RequestingApiKey = ""

		var contentDeal []model.ContentDeal
		node.DB.Raw("select cd.* from content_deals cd, contents c where cd.content = c.id and c.id = ?", contentId).Scan(&contentDeal)

		var pieceCommitments []model.PieceCommitment
		node.DB.Raw("select pc.* from piece_commitments pc, contents c where c.piece_commitment_id = pc.id and c.id = ?", contentId).Scan(&pieceCommitments)

		var contentDealProposal []model.ContentDealProposal
		node.DB.Raw("select cdp.* from content_deal_proposals cdp, contents c where cdp.content = c.id and c.id = ?", contentId).Scan(&contentDealProposal)

		contentResponse = append(contentResponse, map[string]interface{}{
			"content":           content,
			"deals":             contentDeal,
			"piece_commitments": pieceCommitments,
			"deal_proposals":    contentDealProposal,
		})
	}
	return c.JSON(200, contentResponse)

}

// function to get all stats given a content id and user api key
func handleOpenGetStatsByContent(c echo.Context, node *core.DeltaNode) error {
	var content model.Content
	node.DB.Raw("select c.* from contents c where c.id = ?", c.Param("contentId")).Scan(&content)
	content.RequestingApiKey = ""

	var contentDeal []model.ContentDeal
	node.DB.Raw("select cd.* from content_deals cd, contents c where cd.content = c.id and c.id = ?", c.Param("contentId")).Scan(&contentDeal)

	var pieceCommitments []model.PieceCommitment
	node.DB.Raw("select pc.* from piece_commitments pc, contents c where c.piece_commitment_id = pc.id and c.id = ?", c.Param("contentId")).Scan(&pieceCommitments)

	var contentDealProposal []model.ContentDealProposalParameters
	node.DB.Raw("select cdp.* from content_deal_proposal_parameters cdp, contents c where cdp.content = c.id and c.id = ?", c.Param("contentId")).Scan(&contentDealProposal)

	return c.JSON(200, map[string]interface{}{
		"content":           content,
		"deals":             contentDeal,
		"piece_commitments": pieceCommitments,
		"deal_proposals":    contentDealProposal,
	})
}
