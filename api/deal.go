package api

import (
	"delta/core"
	"delta/core/model"
	"delta/jobs"
	"delta/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type CidRequest struct {
	Cids []string `json:"cids"`
}

type ContentMakeDealResponse struct {
	Status            string                 `json:"status"`
	Message           string                 `json:"message"`
	ContentID         int64                  `json:"content_id,omitempty"`
	PieceCommitmentId int64                  `json:"piece_commitment_id,omitempty"`
	Miner             string                 `json:"miner,omitempty"`
	Meta              ContentMakeDealRequest `json:"meta,omitempty"`
}

type WalletRequest struct {
	KeyType    string `json:"key_type,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
}

type PieceCommitmentRequest struct {
	Piece             string `json:"piece,omitempty"`
	PaddedPieceSize   uint64 `json:"padded_piece_size,omitempty"`
	UnPaddedPieceSize uint64 `json:"unpadded_piece_size,omitempty"`
}

type ContentMakeDealRequest struct {
	Cid             string                 `json:"cid,omitempty"`
	Miner           string                 `json:"miner,omitempty"`
	Duration        int64                  `json:"duration,omitempty"`
	Wallet          WalletRequest          `json:"wallet,omitempty"`
	PieceCommitment PieceCommitmentRequest `json:"commp,omitempty"`
	ConnectionMode  string                 `json:"connection_mode,omitempty"`
	Size            int64                  `json:"size,omitempty"`
	StartEpoch      int64                  `json:"start_epoch,omitempty"`
	Replication     int64                  `json:"replication,omitempty"`
}

func DealRouter(e *echo.Group, node *core.DeltaNode) {

	//	inject the stats service
	statsService := core.NewStatsStatsService(node)

	dealMake := e.Group("/deal/make")
	dealStatus := e.Group("/deal/status")

	// Save content to the database
	// Save default values for the content (default miner and wallet)
	dealMake.POST("/content", func(c echo.Context) error {
		return handleContentAdd(c, node, *statsService)
	})
	dealMake.POST("/commitment-piece", func(c echo.Context) error {

		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		var contentMakeDealRequest ContentMakeDealRequest
		//err := json.NewDecoder(c.Request().Body).Decode(&contentMakeDealRequest)
		err := (&echo.DefaultBinder{}).BindBody(c, &contentMakeDealRequest)
		if err != nil {
			c.JSON(500, ContentMakeDealResponse{
				Status:  "error",
				Message: "Error pinning the file" + err.Error(),
			})
		}
		var connMode = contentMakeDealRequest.ConnectionMode
		if connMode == "" || (connMode != utils.CONNECTION_MODE_ONLINE && connMode != utils.CONNECTION_MODE_OFFLINE) {
			connMode = "online"
		}

		var pieceCommp model.PieceCommitment
		if contentMakeDealRequest.PieceCommitment.Piece != "" {
			// if commp is there, make sure the piece and size are there. Use default duration.
			pieceCommp.Cid = contentMakeDealRequest.Cid
			pieceCommp.Piece = contentMakeDealRequest.PieceCommitment.Piece
			pieceCommp.Size = contentMakeDealRequest.Size
			pieceCommp.UnPaddedPieceSize = contentMakeDealRequest.PieceCommitment.UnPaddedPieceSize
			pieceCommp.PaddedPieceSize = contentMakeDealRequest.PieceCommitment.PaddedPieceSize
			pieceCommp.CreatedAt = time.Now()
			pieceCommp.UpdatedAt = time.Now()
			pieceCommp.Status = utils.COMMP_STATUS_OPEN
			node.DB.Create(&pieceCommp)
		}

		// save the file to the database.
		fileSize := contentMakeDealRequest.Size

		content := model.Content{
			Name:              contentMakeDealRequest.Cid,
			Size:              int64(fileSize),
			Cid:               contentMakeDealRequest.Cid,
			RequestingApiKey:  authParts[1],
			PieceCommitmentId: pieceCommp.ID,
			Status:            utils.CONTENT_PINNED,
			ConnectionMode:    connMode,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		node.DB.Create(&content)

		//	assign a miner
		if contentMakeDealRequest.Miner != "" {
			contentMinerAssignment := model.ContentMiner{
				Miner:     contentMakeDealRequest.Miner,
				Content:   content.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			node.DB.Create(&contentMinerAssignment)
		}

		// 	assign a wallet
		if contentMakeDealRequest.Wallet.KeyType != "" {

			var hexedWallet WalletRequest
			hexedWallet.KeyType = contentMakeDealRequest.Wallet.KeyType
			hexedWallet.PrivateKey = hex.EncodeToString([]byte(contentMakeDealRequest.Wallet.PrivateKey))
			walletByteArr, err := json.Marshal(hexedWallet)

			if err != nil {
				c.JSON(500, ContentMakeDealResponse{
					Status:  "error",
					Message: "Error pinning the file" + err.Error(),
				})
			}
			contentWalletAssignment := model.ContentWallet{
				Wallet:    string(walletByteArr),
				Content:   content.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			node.DB.Create(&contentWalletAssignment)
		}

		var dealProposalParam model.ContentDealProposalParameters
		dealProposalParam.CreatedAt = time.Now()
		dealProposalParam.UpdatedAt = time.Now()
		dealProposalParam.Content = content.ID
		dealProposalParam.StartEpoch = contentMakeDealRequest.StartEpoch
		dealProposalParam.Duration = contentMakeDealRequest.Duration
		node.DB.Create(&dealProposalParam)

		var dispatchJobs core.IProcessor
		if pieceCommp.ID != 0 {
			dispatchJobs = jobs.NewStorageDealMakerProcessor(node, content, pieceCommp) // straight to storage deal making
			// add the job so we can process it later
		} else {
			dispatchJobs = jobs.NewPieceCommpProcessor(node, content) // straight to pieceCommp
			// add the job so we can process it later
		}

		node.Dispatcher.AddJobAndDispatch(dispatchJobs, 1)

		c.JSON(200, ContentMakeDealResponse{
			Status:            "success",
			Message:           "File uploaded and pinned successfully",
			ContentID:         content.ID,
			PieceCommitmentId: pieceCommp.ID,
			Miner:             contentMakeDealRequest.Miner,
			Meta:              contentMakeDealRequest,
		})

		return nil
	})
	dealStatus.POST("/content/:contentId", func(c echo.Context) error {
		return handleContentStats(c, node, *statsService)
	})
	dealStatus.POST("/commitment-piece/:piece-commitmentId", func(c echo.Context) error {
		return handleCommitmentPieceStats(c, node, *statsService)

	})
}

func handleContentAdd(c echo.Context, node *core.DeltaNode, stats core.StatsService) error {
	var contentMakeDealRequest ContentMakeDealRequest

	// lets record this.
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")
	file, err := c.FormFile("data") // file
	meta := c.FormValue("metadata")

	//	validate the meta
	json.Unmarshal([]byte(meta), &contentMakeDealRequest)
	var connMode = contentMakeDealRequest.ConnectionMode
	if connMode == "" || (connMode != utils.CONNECTION_MODE_ONLINE && connMode != utils.CONNECTION_MODE_OFFLINE) {
		connMode = "online"
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, ContentMakeDealResponse{
			Status:  "error",
			Message: "Error pinning the file" + err.Error(),
		})
		return err
	}

	addNode, err := node.Node.AddPinFile(c.Request().Context(), src, nil)

	// if size is given, let's create a commp record for it.
	var pieceCommp model.PieceCommitment
	if contentMakeDealRequest.PieceCommitment.Piece != "" {
		// if commp is there, make sure the piece and size are there. Use default duration.
		pieceCommp.Cid = addNode.Cid().String()
		pieceCommp.Piece = contentMakeDealRequest.PieceCommitment.Piece
		pieceCommp.Size = file.Size
		pieceCommp.UnPaddedPieceSize = contentMakeDealRequest.PieceCommitment.UnPaddedPieceSize
		pieceCommp.PaddedPieceSize = contentMakeDealRequest.PieceCommitment.PaddedPieceSize
		pieceCommp.CreatedAt = time.Now()
		pieceCommp.UpdatedAt = time.Now()
		pieceCommp.Status = utils.COMMP_STATUS_OPEN
		node.DB.Create(&pieceCommp)
	}

	// save the file to the database.
	content := model.Content{
		Name:              file.Filename,
		Size:              file.Size,
		Cid:               addNode.Cid().String(),
		RequestingApiKey:  authParts[1],
		PieceCommitmentId: pieceCommp.ID,
		Status:            utils.CONTENT_PINNED,
		ConnectionMode:    connMode,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	node.DB.Create(&content)

	//	assign a miner
	if contentMakeDealRequest.Miner != "" {
		contentMinerAssignment := model.ContentMiner{
			Miner:     contentMakeDealRequest.Miner,
			Content:   content.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		node.DB.Create(&contentMinerAssignment)
	}

	// 	assign a wallet
	if contentMakeDealRequest.Wallet.KeyType != "" {
		var hexedWallet WalletRequest
		hexedWallet.KeyType = contentMakeDealRequest.Wallet.KeyType
		hexedWallet.PrivateKey = hex.EncodeToString([]byte(contentMakeDealRequest.Wallet.PrivateKey))
		walletByteArr, err := json.Marshal(hexedWallet)

		if err != nil {
			c.JSON(500, ContentMakeDealResponse{
				Status:  "error",
				Message: "Error pinning the file" + err.Error(),
			})
		}
		contentWalletAssignment := model.ContentWallet{
			Wallet:    string(walletByteArr),
			Content:   content.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		node.DB.Create(&contentWalletAssignment)
	}

	var dealProposalParam model.ContentDealProposalParameters
	dealProposalParam.CreatedAt = time.Now()
	dealProposalParam.UpdatedAt = time.Now()
	dealProposalParam.Content = content.ID

	fmt.Println("contentMakeDealRequest.StartEpoch", contentMakeDealRequest.StartEpoch)

	fmt.Println("duration", contentMakeDealRequest.Duration)
	if dealProposalParam.Duration == 0 {
		dealProposalParam.Duration = utils.DEFAULT_DURATION
	} else {
		dealProposalParam.Duration = contentMakeDealRequest.Duration
	}

	node.DB.Create(&dealProposalParam)

	//	error
	if err != nil {
		c.JSON(500, ContentMakeDealResponse{
			Status:  "error",
			Message: "Error pinning the file" + err.Error(),
		})
	}

	var dispatchJobs core.IProcessor
	if pieceCommp.ID != 0 {
		dispatchJobs = jobs.NewStorageDealMakerProcessor(node, content, pieceCommp) // straight to storage deal making
	} else {
		dispatchJobs = jobs.NewPieceCommpProcessor(node, content) // straight to pieceCommp
	}

	node.Dispatcher.AddJobAndDispatch(dispatchJobs, 1)

	c.JSON(200, ContentMakeDealResponse{
		Status:            "success",
		Message:           "File uploaded and pinned successfully",
		ContentID:         content.ID,
		PieceCommitmentId: pieceCommp.ID,
		Miner:             contentMakeDealRequest.Miner,
		Meta:              contentMakeDealRequest,
	})

	return nil
}

func handleContentStats(c echo.Context, node *core.DeltaNode, statsService core.StatsService) error {
	contentIdParam := c.Param("contentId")
	contentId, err := strconv.Atoi(contentIdParam)
	if err != nil {
		c.JSON(500, ContentMakeDealResponse{
			Status:  "error",
			Message: "Error looking up the status of the content" + err.Error(),
		})
	}

	status, err := statsService.ContentStatus(core.ContentStatsParam{
		ContentId: int64(contentId),
	})

	if err != nil {
		c.JSON(500, ContentMakeDealResponse{
			Status:  "error",
			Message: "Error looking up the status of the content" + err.Error(),
		})

	}

	return c.JSON(200, status)
}

func handleCommitmentPieceStats(c echo.Context, node *core.DeltaNode, statsService core.StatsService) error {
	pieceCommitmentIdParam := c.Param("piece-commitmentId")
	pieceCommitmentId, err := strconv.Atoi(pieceCommitmentIdParam)
	if err != nil {
		c.JSON(500, ContentMakeDealResponse{
			Status:  "error",
			Message: "Error looking up the status of the piece commitment" + err.Error(),
		})
	}

	status, err := statsService.PieceCommitmentStatus(core.PieceCommitmentStatsParam{
		PieceCommpId: int64(pieceCommitmentId),
	})

	if err != nil {
		c.JSON(500, ContentMakeDealResponse{
			Status:  "error",
			Message: "Error looking up the status of the piece commitment" + err.Error(),
		})
	}

	return c.JSON(200, status)
}

type ValidateMetaResult struct {
	IsValid bool
	Message string
}

func ValidateMeta(meta string) ValidateMetaResult {
	var makeDealMeta ContentMakeDealRequest
	err := json.Unmarshal([]byte(meta), &makeDealMeta)
	if err != nil {
		return ValidateMetaResult{
			IsValid: false,
			Message: "Invalid meta data",
		}
	}

	return ValidateMetaResult{
		IsValid: true,
		Message: "Meta is valid",
	}
}
