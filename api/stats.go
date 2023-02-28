package api

import (
	"delta/core"
	model "github.com/application-research/delta-db/db_models"
	"strings"

	"github.com/labstack/echo/v4"
)

type StatsCheckResponse struct {
	Content struct {
		ID      int64  `json:"id"`
		Name    string `json:"name"`
		Cid     string `json:"cid,omitempty"`
		Status  string `json:"status"`
		Message string `json:"message,omitempty"`
	} `json:"content"`
}

// ConfigureStatsCheckRouter Creating a new router and adding a route to it.
// It configures the router for the stats check API
func ConfigureStatsCheckRouter(e *echo.Group, node *core.DeltaNode) {

	e.GET("/stats/miner/:minerId/content", func(c echo.Context) error {
		return handleGetContentsByMiner(c, node)
	})

	e.GET("/stats/miner/:minerId/piece-commitment", func(c echo.Context) error {
		return handleGetCommitmentPiecesByMiner(c, node)
	})

	e.GET("/stats/miner/:minerId/deals", func(c echo.Context) error {
		return handleGetDealsByMiner(c, node)
	})

	e.GET("/stats/miner/:minerId/deal-proposals", func(c echo.Context) error {
		return handleGetDealProposalsByMiner(c, node)
	})

	e.GET("/stats/miner/:minerId/content/:contentId", func(c echo.Context) error {
		return handleGetContentByMiner(c, node)
	})

	e.GET("/stats/miner/:minerId/piece-commitment/:commitmentPieceId", func(c echo.Context) error {
		return handleGetCommitmentPieceByMiner(c, node)
	})

	e.GET("/stats/miner/:minerId/deals/:dealId", func(c echo.Context) error {
		return handleGetDealByMiner(c, node)
	})

	e.GET("/stats/miner/:minerId/deal-proposals/:dealProposalId", func(c echo.Context) error {
		return handleGetDealProposalByMiner(c, node)
	})

	e.GET("/stats/content/:contentId", func(c echo.Context) error {
		return handleGetStatsByContent(c, node)
	})

	e.GET("/stats/piece-commitment/:commitmentPieceId", func(c echo.Context) error {
		return handleGetCommitmentPiece(c, node)
	})

	e.GET("/stats/piece-commitments", func(c echo.Context) error {
		return handleGetCommitmentPieces(c, node)
	})

	e.GET("/stats/deal/:dealId", func(c echo.Context) error {
		return nil
	})

	e.GET("/stats/deal-proposal/:id", func(c echo.Context) error {
		return nil
	})

	e.GET("/stats/deals", func(c echo.Context) error {
		return handleGetDeals(c, node)
	})

	e.GET("/stats/contents", func(c echo.Context) error {
		return handleGetContents(c, node)
	})

	e.GET("/stats/miner/:minerId", handleMinerStats(node))

	e.GET("/stats", handleStats(node))

}

func handleStats(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {

		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		// select * from content_deals cd, content c where cd.content = c.id and c.requesting_api_key = ?;
		var content []model.Content
		node.DB.Raw("select c.* from content_deals cd, contents c where cd.content = c.id and c.requesting_api_key = ?", authParts[1]).Scan(&content)

		var contentDeal []model.ContentDeal
		node.DB.Raw("select cd.* from content_deals cd, contents c where cd.content = c.id and c.requesting_api_key = ?", authParts[1]).Scan(&contentDeal)

		// select * from piece_commitments pc, content c where c.piece_commitment_id = pc.id and c.requesting_api_key = ?;
		var pieceCommitments []model.PieceCommitment
		node.DB.Raw("select pc.* from piece_commitments pc, contents c where c.piece_commitment_id = pc.id and c.requesting_api_key = ?", authParts[1]).Scan(&pieceCommitments)

		return c.JSON(200, map[string]interface{}{
			"content":           content,
			"deals":             contentDeal,
			"piece_commitments": pieceCommitments,
		})
	}
}

// `handleMinerStats` is a function that takes a `*core.DeltaNode` and returns a function that takes an `echo.Context` and
// returns an `error`
// `handleMinerStats` is a function that takes a `*core.DeltaNode` and returns a function that takes an `echo.Context` and
// returns an `error`
func handleMinerStats(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {

		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		var contents []model.Content
		node.DB.Raw("select c.* from content_deals cd, contents c where cd.content = c.id and cd.miner = ? and c.requesting_api_key = ?", c.Param("minerId"), authParts[1]).Scan(&contents)

		var contentMinerAssignment []model.ContentMiner
		node.DB.Raw("select cma.* from content_miners cma, contents c where cma.content = c.id and cma.miner = ? and c.requesting_api_key = ?", c.Param("minerId"), authParts[1]).Scan(&contentMinerAssignment)

		return c.JSON(200, map[string]interface{}{
			"content": contents,
			"cmas":    contentMinerAssignment,
		})
		return nil
	}
}

// A function that takes in a commitment and a piece number and returns the piece of the commitment.
func handleGetCommitmentPiece(c echo.Context, node *core.DeltaNode) error {
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	// select * from piece_commitments pc, content c where c.piece_commitment_id = pc.id and c.requesting_api_key = ?;
	var pieceCommitments []model.PieceCommitment
	node.DB.Raw("select pc.* from piece_commitments pc, contents c where c.piece_commitment_id = ? and c.requesting_api_key = ?", c.Param("commitmentPieceId"), authParts[1]).Scan(&pieceCommitments)

	return c.JSON(200, map[string]interface{}{
		"piece_commitments": pieceCommitments,
	})
	return nil
}

// function to get all stats given a content id and user api key
func handleGetStatsByContent(c echo.Context, node *core.DeltaNode) error {
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	var content model.Content
	node.DB.Raw("select c.* from contents c where c.id = ? and c.requesting_api_key = ?", c.Param("contentId"), authParts[1]).Scan(&content)

	var contentDeal []model.ContentDeal
	node.DB.Raw("select cd.* from content_deals cd, contents c where cd.content = c.id and c.id = ? and c.requesting_api_key = ?", c.Param("contentId"), authParts[1]).Scan(&contentDeal)

	var pieceCommitments []model.PieceCommitment
	node.DB.Raw("select pc.* from piece_commitments pc, contents c where c.piece_commitment_id = pc.id and c.id = ? and c.requesting_api_key = ?", c.Param("contentId"), authParts[1]).Scan(&pieceCommitments)

	var contentDealProposal []model.ContentDealProposalParameters
	node.DB.Raw("select cdp.* from content_deal_proposal_parameters cdp, contents c where cdp.content = c.id and c.id = ? and c.requesting_api_key = ?", c.Param("contentId"), authParts[1]).Scan(&contentDealProposal)

	return c.JSON(200, map[string]interface{}{
		"content":           content,
		"deals":             contentDeal,
		"piece_commitments": pieceCommitments,
		"deal_proposals":    contentDealProposal,
	})
}

// function to get all contents of a given a miner
func handleGetContentsByMiner(c echo.Context, node *core.DeltaNode) error {
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	var contents []model.Content
	node.DB.Raw("select c.* from content_deals cd, contents c where cd.content = c.id and cd.miner = ? and c.requesting_api_key = ?", c.Param("minerId"), authParts[1]).Scan(&contents)

	c.JSON(200, map[string]interface{}{
		"content": contents,
	})

	return nil
}

// function to get all piece-commitment of a given miner
func handleGetCommitmentPiecesByMiner(c echo.Context, node *core.DeltaNode) error {
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	var pieceCommitments []model.PieceCommitment
	node.DB.Raw("select pc.* from piece_commitments pc, contents c where c.piece_commitment_id = pc.id and c.requesting_api_key = ?", authParts[1]).Scan(&pieceCommitments)

	c.JSON(200, map[string]interface{}{
		"piece_commitments": pieceCommitments,
	})

	return nil
}

// function to get all deals of a given miner
func handleGetDealsByMiner(c echo.Context, node *core.DeltaNode) error {
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	var contentDeal []model.ContentDeal
	node.DB.Raw("select cd.* from content_deals cd, contents c where cd.content = c.id and c.requesting_api_key = ? and cd.miner = ?", authParts[1], c.Param("minerId")).Scan(&contentDeal)

	c.JSON(200, map[string]interface{}{
		"deals": contentDeal,
	})

	return nil
}

// function to get all deal-proposal of a given miner
func handleGetDealProposalsByMiner(c echo.Context, node *core.DeltaNode) error {
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	var contentMinerAssignment []model.ContentMiner
	node.DB.Raw("select cma.* from content_miners cma, contents c where cma.content = c.id and cma.miner = ? and c.requesting_api_key = ?", c.Param("minerId"), authParts[1]).Scan(&contentMinerAssignment)

	c.JSON(200, map[string]interface{}{
		"cmas": contentMinerAssignment,
	})

	return nil
}

// function to get all content of a given api key
func handleGetContents(c echo.Context, node *core.DeltaNode) error {
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	var content []model.Content
	node.DB.Raw("select c.* from content_deals cd, contents c where cd.content = c.id and c.requesting_api_key = ?", authParts[1]).Scan(&content)

	c.JSON(200, map[string]interface{}{
		"content": content,
	})

	return nil
}

// function to get all piece-commitment of a given api key
func handleGetCommitmentPieces(c echo.Context, node *core.DeltaNode) error {
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	var pieceCommitments []model.PieceCommitment
	node.DB.Raw("select pc.* from piece_commitments pc, contents c where c.piece_commitment_id = pc.id and c.requesting_api_key = ?", authParts[1]).Scan(&pieceCommitments)

	c.JSON(200, map[string]interface{}{
		"piece_commitments": pieceCommitments,
	})

	return nil
}

// function to get all deals of a given api key
func handleGetDeals(c echo.Context, node *core.DeltaNode) error {
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	var contentDeal []model.ContentDeal
	node.DB.Raw("select cd.* from content_deals cd, contents c where cd.content = c.id and c.requesting_api_key = ?", authParts[1]).Scan(&contentDeal)

	c.JSON(200, map[string]interface{}{
		"deals": contentDeal,
	})

	return nil
}

// function to get all deal-proposal of a given api key
func handleGetDealProposals(c echo.Context, node *core.DeltaNode) error {
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	var contentMinerAssignment []model.ContentMiner
	node.DB.Raw("select cma.* from content_miners cma, contents c where cma.content = c.id and c.requesting_api_key = ?", authParts[1]).Scan(&contentMinerAssignment)

	c.JSON(200, map[string]interface{}{
		"cmas": contentMinerAssignment,
	})

	return nil
}

// function to get a specific content with a given miner
func handleGetContentByMiner(c echo.Context, node *core.DeltaNode) error {
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	var content model.Content
	node.DB.Raw("select c.* from content_deals cd, contents c where cd.content = ? and cd.miner = ? and c.requesting_api_key = ?", c.Param("contentId"), c.Param("minerId"), authParts[1]).Scan(&content)

	c.JSON(200, map[string]interface{}{
		"content": content,
	})

	return nil
}

// function to get a specific piece-commitment with a given miner
func handleGetCommitmentPieceByMiner(c echo.Context, node *core.DeltaNode) error {
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	var pieceCommitment model.PieceCommitment
	node.DB.Raw("select pc.* from piece_commitments pc, contents c where c.piece_commitment_id = ? and c.miner = ? and c.id = ? and c.requesting_api_key = ?", c.Param("commitmentPieceId"), c.Param("minerId"), authParts[1]).Scan(&pieceCommitment)

	c.JSON(200, map[string]interface{}{
		"piece_commitment": pieceCommitment,
	})

	return nil
}

// function to get a specific deal with a given miner
func handleGetDealByMiner(c echo.Context, node *core.DeltaNode) error {
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	var contentDeal model.ContentDeal
	node.DB.Raw("select cd.* from content_deals cd, contents c where cd.id = ? and cd.miner = ? and c.requesting_api_key = ?", c.Param("dealId"), c.Param("minerId"), authParts[1]).Scan(&contentDeal)

	c.JSON(200, map[string]interface{}{
		"deal": contentDeal,
	})

	return nil
}

// function to get a specific deal-proposal with a given miner
func handleGetDealProposalByMiner(c echo.Context, node *core.DeltaNode) error {
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	var contentMinerAssignment model.ContentMiner
	node.DB.Raw("select cma.* from content_deal_proposals cdp, contents c where cdp.id = ? and cdp.miner = ? and c.requesting_api_key = ?", c.Param("dealProposalId"), c.Param("minerId"), authParts[1]).Scan(&contentMinerAssignment)

	c.JSON(200, map[string]interface{}{
		"cmas": contentMinerAssignment,
	})

	return nil
}
