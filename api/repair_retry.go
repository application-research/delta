package api

import (
	"delta/core"
	"delta/jobs"
	"delta/utils"
	"encoding/json"
	model "github.com/application-research/delta-db/db_models"
	"github.com/labstack/echo/v4"
	"strings"
	"time"
)

type RetryDealResponse struct {
	Status       string      `json:"status"`
	Message      string      `json:"message"`
	NewContentId int64       `json:"new_content_id,omitempty"`
	OldContentId interface{} `json:"old_content_id,omitempty"`
}

type MultipleImportRequest struct {
	ContentID   string      `json:"content_id"`
	DealRequest DealRequest `json:"metadata"`
}

type ImportRetryRequest struct {
	ContentIds []string `json:"content_ids"`
}
type ImportRetryResponse struct {
	Message     string        `json:"message"`
	Content     model.Content `json:"content"`
	DealRequest DealRequest   `json:"metadata"`
}

// ConfigureRepairRouter repair deals (re-create or re-try)
// It's a function that configures the repair router
func ConfigureRepairRouter(e *echo.Group, node *core.DeltaNode) {

	// repair with a different (miner, duration only)
	repair := e.Group("/repair")
	repair.GET("/deal/end-to-end/:contentId", handleRepairDealContent(node))
	repair.GET("/deal/import/:contentId", handleRepairImportContent(node))
	repair.POST("/deal/imports", handleRepairMultipleImport(node))

	// retry
	retry := e.Group("/retry")
	retry.GET("/deal/end-to-end/:contentId", handleRetryDealContent(node))
	retry.GET("/deal/import/:contentId", handleRetryDealImport(node))
	retry.POST("/deal/imports", handleRetryMultipleImport(node))

	// disable auto-retry
	autoRetry := e.Group("/auto-retry")
	autoRetry.GET("/deal/disable/:contentId", handleDisableAutoRetry(node))
	autoRetry.GET("/deal/enable/:contentId", handleEnableAutoRetry(node))
}

func handleDisableAutoRetry(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {

		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		var contentId = c.Param("contentId")
		var content model.Content

		node.DB.Model(&model.Content{}).Where("id = ? AND requesting_api_key = ?", contentId, authParts[1]).First(&content)
		content.AutoRetry = false
		node.DB.Model(&model.Content{}).Save(&content)
		return nil
	}
}

// The function handles enabling auto-retry for a specific content ID in a Go application.
func handleEnableAutoRetry(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {

		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		var contentId = c.Param("contentId")
		var content model.Content

		node.DB.Model(&model.Content{}).Where("id = ? AND requesting_api_key = ?", contentId, authParts[1]).First(&content)
		content.AutoRetry = true
		node.DB.Model(&model.Content{}).Save(&content)
		return nil
	}
}

// > This function handles the retry of a deal content
// This function handles retrying a content deal and returns a JSON response.
func handleRetryDealContent(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {

		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		paramContentId := c.Param("contentId")

		// if the deal is not in the right state, throw an error.
		var content model.Content
		node.DB.Model(&model.Content{}).Where("id = ? AND requesting_api_key = ?", paramContentId).First(&content)
		content.RequestingApiKey = ""

		if content.ConnectionMode != utils.CONNECTION_MODE_E2E {
			return c.JSON(200, map[string]interface{}{
				"message": "content is not in end-to-end mode",
			})
		}

		// get the content deal entry
		var contentDeal model.ContentDeal
		node.DB.Model(&model.ContentDeal{}).Where("content = ?", paramContentId, authParts[1]).First(&contentDeal)

		// if not content deal entry, throw an error.
		if contentDeal.ID == 0 {
			return c.JSON(200, map[string]interface{}{
				"message": "content deal not found",
			})
		}

		// retry it.
		processor := jobs.NewPieceCommpProcessor(node, content)
		node.Dispatcher.AddJobAndDispatch(processor, 1)

		return c.JSON(200, map[string]interface{}{
			"message": "retrying deal",
			"content": content,
		})
	}
}

// `handleRepairDealContent` is a function that takes a `DeltaNode` and returns a function that takes a `Context` and
// returns an `error`
func handleRepairDealContent(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")
		paramContentId := c.Param("contentId")
		meta := c.FormValue("metadata") // only allow miner and durations

		// if the deal is not in the right state, throw an error.
		var content model.Content
		node.DB.Model(&model.Content{}).Where("id = ? and requesting_api_key = ?", paramContentId, authParts[1]).First(&content)
		content.RequestingApiKey = ""

		if content.ID == 0 {
			return c.JSON(200, map[string]interface{}{
				"message": "content not found",
			})
		}

		if content.ConnectionMode != utils.CONNECTION_MODE_E2E {
			return c.JSON(200, map[string]interface{}{
				"message": "content is not in end-to-end mode",
			})
		}

		var dealRequest DealRequest
		err := json.Unmarshal([]byte(meta), &dealRequest)
		if err != nil {
			return err
		}

		// validate the deal request
		if (DealRequest{} != dealRequest && dealRequest.StartEpochInDays > 14) {
			return c.JSON(200, map[string]interface{}{
				"message": "start epoch cannot be more than 14 days",
			})
		}

		if (DealRequest{} != dealRequest && dealRequest.DurationInDays > 540) {
			return c.JSON(200, map[string]interface{}{
				"message": "duration cannot be more than 540 days",
			})
		}

		if dealRequest.StartEpochInDays > dealRequest.DurationInDays {
			return c.JSON(200, map[string]interface{}{
				"message": "start epoch cannot be more than duration",
			})
		}

		// get the content deal entry
		var contentDeal model.ContentDeal
		node.DB.Model(&model.ContentDeal{}).Where("content = ?", paramContentId).First(&contentDeal)

		// if not content deal entry, throw an error.
		if contentDeal.ID == 0 {
			return c.JSON(200, map[string]interface{}{
				"message": "content deal not found",
			})
		}

		// get the proposal
		var dealProposalParam model.ContentDealProposalParameters
		node.DB.Model(&model.ContentDealProposalParameters{}).Where("content = ?", paramContentId).First(&dealProposalParam)

		if dealProposalParam.ID == 0 {
			return c.JSON(200, map[string]interface{}{
				"message": "content deal proposal not found",
			})
		}

		// get the miner entry
		var dealContentMiner model.ContentMiner
		node.DB.Model(&model.ContentMiner{}).Where("content = ?", paramContentId).First(&dealContentMiner)

		// only change the miner and duration
		// create new content miner record.
		dealContentMiner.Miner = dealRequest.Miner
		if dealRequest.StartEpochInDays != 0 && dealRequest.DurationInDays != 0 {
			startEpochTime := time.Now().AddDate(0, 0, int(dealRequest.StartEpochInDays))
			dealProposalParam.StartEpoch = utils.DateToHeight(startEpochTime)
			dealProposalParam.EndEpoch = dealProposalParam.StartEpoch + (utils.EPOCH_PER_DAY * (dealRequest.DurationInDays - dealRequest.StartEpochInDays))
			dealProposalParam.Duration = dealProposalParam.EndEpoch - dealProposalParam.StartEpoch
		} else {
			dealProposalParam.StartEpoch = 0
			dealProposalParam.Duration = utils.DEFAULT_DURATION
		}

		var pieceComm model.PieceCommitment
		node.DB.Model(&model.PieceCommitment{}).Where("id = ?", content.PieceCommitmentId).First(&pieceComm)

		// retry it.
		processor := jobs.NewStorageDealMakerProcessor(node, content, pieceComm)
		node.Dispatcher.AddJobAndDispatch(processor, 1)

		return c.JSON(200, map[string]interface{}{
			"message": "retrying deal",
			"content": content,
		})
	}
}

// This function handles a request to retry multiple failed import deals and returns a JSON response.
func handleRepairMultipleImport(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")
		var multipleImportRequest []MultipleImportRequest
		err := c.Bind(&multipleImportRequest)
		if err != nil {
			return c.JSON(200, map[string]interface{}{
				"message": "invalid request",
			})
		}
		var importResponse []ImportRetryResponse
		for _, request := range multipleImportRequest {
			paramContentId := request.ContentID
			dealRequest := request.DealRequest

			// validate the deal request
			if (DealRequest{} != dealRequest && dealRequest.StartEpochInDays > 14) {
				importResponse = append(importResponse, ImportRetryResponse{
					Message: "start epoch cannot be more than 14 days",
				})
			}

			if (DealRequest{} != dealRequest && dealRequest.DurationInDays > 540) {
				return c.JSON(200, map[string]interface{}{
					"message": "duration cannot be more than 540 days",
				})
			}

			if dealRequest.StartEpochInDays > dealRequest.DurationInDays {
				return c.JSON(200, map[string]interface{}{
					"message": "start epoch cannot be more than duration",
				})
			}

			// if the deal is not in the right state, throw an error.
			var content model.Content
			node.DB.Model(&model.Content{}).Where("id = ? and requesting_api_key = ?", paramContentId, authParts[1]).First(&content)
			content.RequestingApiKey = ""

			if content.ConnectionMode != utils.CONNECTION_MODE_IMPORT {
				importResponse = append(importResponse, ImportRetryResponse{
					Message: "content is not in import mode",
					Content: content,
				})
			}

			// get the content deal entry
			var contentDeal model.ContentDeal
			node.DB.Model(&model.ContentDeal{}).Where("content = ?", paramContentId).First(&contentDeal)

			// if not content deal entry, throw an error.
			if contentDeal.ID == 0 {
				importResponse = append(importResponse, ImportRetryResponse{
					Message: "content deal not found",
					Content: content,
				})
			}

			// get the proposal
			var dealProposalParam model.ContentDealProposalParameters
			node.DB.Model(&model.ContentDealProposalParameters{}).Where("content = ?", paramContentId).First(&dealProposalParam)

			if dealProposalParam.ID == 0 {
				importResponse = append(importResponse, ImportRetryResponse{
					Message: "content deal proposal not found",
				})
			}

			// get the miner entry
			var dealContentMiner model.ContentMiner
			node.DB.Model(&model.ContentMiner{}).Where("content = ?", paramContentId).First(&dealContentMiner)

			// only change the miner and duration
			// create new content miner record.
			dealContentMiner.Miner = dealRequest.Miner
			if dealRequest.StartEpochInDays != 0 && dealRequest.DurationInDays != 0 {
				startEpochTime := time.Now().AddDate(0, 0, int(dealRequest.StartEpochInDays))
				dealProposalParam.StartEpoch = utils.DateToHeight(startEpochTime)
				dealProposalParam.EndEpoch = dealProposalParam.StartEpoch + (utils.EPOCH_PER_DAY * (dealRequest.DurationInDays - dealRequest.StartEpochInDays))
				dealProposalParam.Duration = dealProposalParam.EndEpoch - dealProposalParam.StartEpoch
			} else {
				dealProposalParam.StartEpoch = 0
				dealProposalParam.Duration = utils.DEFAULT_DURATION
			}

			var pieceComm model.PieceCommitment
			node.DB.Model(&model.PieceCommitment{}).Where("id = ?", content.PieceCommitmentId).First(&pieceComm)

			// retry it.
			processor := jobs.NewStorageDealMakerProcessor(node, content, pieceComm)
			node.Dispatcher.AddJobAndDispatch(processor, 1)

			importResponse = append(importResponse, ImportRetryResponse{
				Message: "retrying deal",
				Content: content,
			})
		}
		return c.JSON(200, importResponse)
	}
}

// This function handles the repair of import content by updating the miner and duration of a content deal and retrying the
// deal.
func handleRepairImportContent(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")
		paramContentId := c.Param("contentId")
		meta := c.FormValue("metadata") // only allow miner and durations

		var dealRequest DealRequest
		err := json.Unmarshal([]byte(meta), &dealRequest)
		if err != nil {
			return err
		}

		// if the deal is not in the right state, throw an error.
		var content model.Content
		node.DB.Model(&model.Content{}).Where("id = ? and requesting_api_key = ?", paramContentId, authParts[1]).First(&content)
		content.RequestingApiKey = ""

		if content.ConnectionMode != utils.CONNECTION_MODE_IMPORT {
			return c.JSON(200, map[string]interface{}{
				"message": "content is not in import mode",
			})
		}

		// get the content deal entry
		var contentDeal model.ContentDeal
		node.DB.Model(&model.ContentDeal{}).Where("content = ?", paramContentId).First(&contentDeal)

		// if not content deal entry, throw an error.
		if contentDeal.ID == 0 {
			return c.JSON(200, map[string]interface{}{
				"message": "content deal not found",
			})
		}

		// get the proposal
		var dealProposalParam model.ContentDealProposalParameters
		node.DB.Model(&model.ContentDealProposalParameters{}).Where("content = ?", paramContentId).First(&dealProposalParam)

		if dealProposalParam.ID == 0 {
			return c.JSON(200, map[string]interface{}{
				"message": "content deal proposal not found",
			})
		}

		// get the miner entry
		var dealContentMiner model.ContentMiner
		node.DB.Model(&model.ContentMiner{}).Where("content = ?", paramContentId).First(&dealContentMiner)

		// only change the miner and duration
		// create new content miner record.
		dealContentMiner.Miner = dealRequest.Miner
		if dealRequest.StartEpochInDays != 0 && dealRequest.DurationInDays != 0 {
			startEpochTime := time.Now().AddDate(0, 0, int(dealRequest.StartEpochInDays))
			dealProposalParam.StartEpoch = utils.DateToHeight(startEpochTime)
			dealProposalParam.EndEpoch = dealProposalParam.StartEpoch + (utils.EPOCH_PER_DAY * (dealRequest.DurationInDays - dealRequest.StartEpochInDays))
			dealProposalParam.Duration = dealProposalParam.EndEpoch - dealProposalParam.StartEpoch
		} else {
			dealProposalParam.StartEpoch = 0
			dealProposalParam.Duration = utils.DEFAULT_DURATION
		}

		var pieceComm model.PieceCommitment
		node.DB.Model(&model.PieceCommitment{}).Where("id = ?", content.PieceCommitmentId).First(&pieceComm)

		// retry it.
		processor := jobs.NewStorageDealMakerProcessor(node, content, pieceComm)
		node.Dispatcher.AddJobAndDispatch(processor, 1)

		return c.JSON(200, map[string]interface{}{
			"message": "retrying deal",
			"content": content,
		})
	}
}

// It takes a content ID, finds the content deal entry for that content, and then retries the deal
// This function handles retrying multiple content deals for import.
func handleRetryMultipleImport(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")
		var importRetryRequest ImportRetryRequest
		err := c.Bind(&importRetryRequest)
		if err != nil {
			return c.JSON(200, map[string]interface{}{
				"message": "error parsing request",
			})
		}

		var importRetryResponse []ImportRetryResponse
		for _, paramContentId := range importRetryRequest.ContentIds {

			// if the deal is not in the right state, throw an error.
			var content model.Content
			node.DB.Model(&model.Content{}).Where("id = ? and requesting_api_key = ?", paramContentId, authParts[1]).First(&content)
			content.RequestingApiKey = ""

			if content.ConnectionMode != utils.CONNECTION_MODE_IMPORT {
				return c.JSON(200, map[string]interface{}{
					"message": "content is not in import mode",
				})
			}

			// get the content deal entry
			var contentDeal model.ContentDeal
			node.DB.Model(&model.ContentDeal{}).Where("content = ?", paramContentId).First(&contentDeal)

			// if not content deal entry, throw an error.
			if contentDeal.ID == 0 {
				return c.JSON(200, map[string]interface{}{
					"message": "content deal not found",
				})
			}

			// retry it.
			processor := jobs.NewPieceCommpProcessor(node, content)
			node.Dispatcher.AddJobAndDispatch(processor, 1)

			importRetryResponse = append(importRetryResponse, ImportRetryResponse{
				Message: "retrying deal",
				Content: content,
			})
		}
		return c.JSON(200, importRetryResponse)
	}
}

// This function handles retrying a content deal import and returns a JSON response.
func handleRetryDealImport(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")
		paramContentId := c.Param("contentId")

		// get the content deal entry
		var contentDeal model.ContentDeal
		node.DB.Model(&model.ContentDeal{}).Where("content = ? and requesting_api_key = ?", paramContentId, authParts[1]).First(&contentDeal)

		// if not content deal entry, throw an error.
		if contentDeal.ID == 0 {
			return c.JSON(200, map[string]interface{}{
				"message": "content deal not found",
			})
		}

		// if the deal is not in the right state, throw an error.
		var content model.Content
		node.DB.Model(&model.Content{}).Where("id = ?", paramContentId).First(&content)
		content.RequestingApiKey = ""

		if content.ConnectionMode != utils.CONNECTION_MODE_IMPORT {
			return c.JSON(200, map[string]interface{}{
				"message": "content is not in import mode",
			})
		}

		// retry it.
		processor := jobs.NewPieceCommpProcessor(node, content)
		node.Dispatcher.AddJobAndDispatch(processor, 1)

		return c.JSON(200, map[string]interface{}{
			"message": "retrying deal",
			"content": content,
		})
	}
}
