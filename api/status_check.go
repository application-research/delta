package api

import (
	"delta/core"
	"delta/core/model"
	"strings"

	"github.com/labstack/echo/v4"
)

type StatusCheckResponse struct {
	Content struct {
		ID      int64  `json:"id"`
		Name    string `json:"name"`
		Cid     string `json:"cid,omitempty"`
		Status  string `json:"status"`
		Message string `json:"message,omitempty"`
	} `json:"content"`
}

func ConfigureStatusCheckRouter(e *echo.Group, node *core.DeltaNode) {

	e.GET("/status/:id", func(c echo.Context) error {

		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		var content model.Content
		node.DB.Raw("select c.id, c.cid, c.status from contents as c where c.id = ? and c.requesting_api_key = ?", c.Param("id"), authParts[1]).Scan(&content)

		return c.JSON(200, StatusCheckResponse{
			Content: struct {
				ID      int64  `json:"id"`
				Name    string `json:"name"`
				Cid     string `json:"cid,omitempty"`
				Status  string `json:"status"`
				Message string `json:"message,omitempty"`
			}{ID: content.ID, Name: content.Name, Status: content.Status, Cid: content.Cid},
		})
	})

	e.GET("/list-all-cids", func(c echo.Context) error {

		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		var content []model.Content
		node.DB.Raw("select c.name, c.id, c.cid, c.status,c.created_at,c.updated_at from contents as c where requesting_api_key = ?", authParts[1]).Scan(&content)

		return c.JSON(200, content)

	})

	e.GET("/stats/commps", func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		// select * from piece_commitments pc, content c where c.piece_commitment_id = pc.id and c.requesting_api_key = ?;
		var pieceCommitments []model.PieceCommitment
		node.DB.Raw("select pc.* from piece_commitments pc, contents c where c.piece_commitment_id = pc.id and c.requesting_api_key = ?", authParts[1]).Scan(&pieceCommitments)

		return c.JSON(200, map[string]interface{}{
			"piece_commitments": pieceCommitments,
		})
		return nil
	})

	e.GET("/stats/content/:id", func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		var content model.Content
		node.DB.Raw("select c.* from contents c where c.id = ? and c.requesting_api_key = ?", c.Param("id"), authParts[1]).Scan(&content)

		var contentDeal []model.ContentDeal
		node.DB.Raw("select cd.* from content_deals cd, contents c where cd.content = c.id and c.id = ? and c.requesting_api_key = ?", c.Param("id"), authParts[1]).Scan(&contentDeal)

		var pieceCommitments []model.PieceCommitment
		node.DB.Raw("select pc.* from piece_commitments pc, contents c where c.piece_commitment_id = pc.id and c.id = ? and c.requesting_api_key = ?", c.Param("id"), authParts[1]).Scan(&pieceCommitments)

		var contentDealProposal []model.ContentDealProposalParameters
		node.DB.Raw("select cdp.* from content_deal_proposal_parameters cdp, contents c where cdp.content = c.id and c.id = ? and c.requesting_api_key = ?", c.Param("id"), authParts[1]).Scan(&contentDealProposal)

		var contentMinerAssignment []model.ContentMiner
		node.DB.Raw("select cma.* from content_miner_assignments cma, contents c where cma.content = c.id and c.id = ? and c.requesting_api_key = ?", c.Param("id"), authParts[1]).Scan(&contentMinerAssignment)

		var contentWalletAssignment []model.ContentWallet
		node.DB.Raw("select cwa.* from content_wallet_assignments cwa, contents c where cwa.content = c.id and c.id = ? and c.requesting_api_key = ?", c.Param("id"), authParts[1]).Scan(&contentWalletAssignment)

		return c.JSON(200, map[string]interface{}{
			"content": content,
			"deals":   contentDeal,
			"commps":  pieceCommitments,
			"cdps":    contentDealProposal,
			"cmas":    contentMinerAssignment,
			"cwas":    contentWalletAssignment,
		})
		return nil
	})

	e.GET("/stats/deals", func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		var contentDeal []model.ContentDeal
		node.DB.Raw("select cd.* from content_deals cd, contents c where cd.content = c.id and c.requesting_api_key = ?", authParts[1]).Scan(&contentDeal)

		return c.JSON(200, map[string]interface{}{
			"deals": contentDeal,
		})
		return nil
	})

	e.GET("/stats/contents", func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		var content []model.Content
		node.DB.Raw("select c.* from content_deals cd, contents c where cd.content = c.id and c.requesting_api_key = ?", authParts[1]).Scan(&content)

		return c.JSON(200, map[string]interface{}{
			"content": content,
		})
		return nil
	})

	e.GET("/stats/miner/:minerId", func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		var contents []model.Content
		node.DB.Raw("select c.* from content_deals cd, contents c where cd.content = c.id and cd.miner = ? and c.requesting_api_key = ?", c.Param("minerId"), authParts[1]).Scan(&contents)

		var contentMinerAssignment []model.ContentMiner
		node.DB.Raw("select cma.* from content_miner_assignments cma, contents c where cma.content = c.id and cma.miner = ? and c.requesting_api_key = ?", c.Param("minerId"), authParts[1]).Scan(&contentMinerAssignment)

		return c.JSON(200, map[string]interface{}{
			"content": contents,
			"cmas":    contentMinerAssignment,
		})
		return nil
	})

	e.GET("/stats", func(c echo.Context) error {

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
	})
}
