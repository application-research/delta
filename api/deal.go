package api

import (
	"delta/core"
	"delta/core/model"
	"delta/jobs"
	"delta/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type CidRequest struct {
	Cids []string `json:"cids"`
}

type WalletRequest struct {
	KeyType    string `json:"key_type,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
}

type PieceCommitmentRequest struct {
	Piece             string `json:"piece_cid,omitempty"`
	PaddedPieceSize   uint64 `json:"padded_piece_size,omitempty"`
	UnPaddedPieceSize uint64 `json:"unpadded_piece_size,omitempty"`
}

type DealRequest struct {
	Cid                  string                 `json:"cid,omitempty"`
	Miner                string                 `json:"miner,omitempty"`
	Duration             int64                  `json:"duration,omitempty"`
	Wallet               WalletRequest          `json:"wallet,omitempty"`
	PieceCommitment      PieceCommitmentRequest `json:"piece_commitment,omitempty"`
	ConnectionMode       string                 `json:"connection_mode,omitempty"`
	Size                 int64                  `json:"size,omitempty"`
	StartEpoch           int64                  `json:"start_epoch,omitempty"`
	Replication          int64                  `json:"replication,omitempty"`
	RemoveUnsealedCopies bool                   `json:"remove_unsealed_copies,omitempty"`
	SkipIPNIAnnounce     bool                   `json:"skip_ipni_announce,omitempty"`
}

type DealResponse struct {
	Status      string      `json:"status"`
	Message     string      `json:"message"`
	ContentId   int64       `json:"content_id,omitempty"`
	DealRequest interface{} `json:"request_meta,omitempty"`
}

func DealRouter(e *echo.Group, node *core.DeltaNode) {

	//	inject the stats service
	statsService := core.NewStatsStatsService(node)

	dealMake := e.Group("/deal")
	dealPrepare := dealMake.Group("/prepare")
	dealAnnounce := dealMake.Group("/announce")
	dealStatus := dealMake.Group("/status")

	dealMake.POST("/content", func(c echo.Context) error {
		return handleContentAdd(c, node, *statsService)
	})
	dealMake.POST("/piece-commitment", func(c echo.Context) error {
		return handleCommPieceAdd(c, node, *statsService)
	})

	// make piece-commitments
	dealMake.POST("/piece-commitments", func(c echo.Context) error {
		return handleCommPiecesAdd(c, node, *statsService)
	})

	dealPrepare.POST("/content", func(c echo.Context) error {
		return nil
	})

	dealPrepare.POST("/piece-commitment", func(c echo.Context) error {
		return nil
	})

	dealPrepare.POST("/piece-commitments", func(c echo.Context) error {
		return nil
	})

	dealAnnounce.POST("/content", func(c echo.Context) error {
		return nil
	})

	dealAnnounce.POST("/piece-commitment", func(c echo.Context) error {
		return nil
	})

	dealAnnounce.POST("/piece-commitments", func(c echo.Context) error {
		return nil
	})

	dealStatus.POST("/content/:contentId", func(c echo.Context) error {
		return handleContentStats(c, node, *statsService)
	})
	dealStatus.POST("/piece-commitment/:piece-commitmentId", func(c echo.Context) error {
		return handleCommitmentPieceStats(c, node, *statsService)

	})
}

// handleContentStats returns the status of a content
// @Summary returns the status of a content
// @Description returns the status of a content
// @Tags deal
// @Accept  json
// @Produce  json
// @Param contentId path int true "Content ID"
// @Success 200 {object} ContentMakeDealResponse
// @Failure 500 {object} ContentMakeDealResponse
// @Router /deal/content/{contentId} [post]
func handleContentAdd(c echo.Context, node *core.DeltaNode, stats core.StatsService) error {
	var dealRequest DealRequest

	// lets record this.
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")
	file, err := c.FormFile("data") // file
	meta := c.FormValue("metadata")

	//	validate the meta
	json.Unmarshal([]byte(meta), &dealRequest)

	err = ValidateMeta(dealRequest)
	if err != nil {
		// return the error from the validation
		return err
	}

	// process the file
	src, err := file.Open()
	if err != nil {
		return errors.New("Error opening the file")
	}

	addNode, err := node.Node.AddPinFile(c.Request().Context(), src, nil)
	if err != nil {
		return errors.New("Error pinning the file")
	}

	// specify the connection mode
	var connMode = dealRequest.ConnectionMode
	if connMode == "" || (connMode != utils.CONNECTION_MODE_E2E && connMode != utils.CONNECTION_MODE_IMPORT) {
		connMode = "e2e"
	}

	// let's create a commp but only if we have
	// a cid, a piece_cid, a padded_piece_size, size
	var pieceCommp model.PieceCommitment
	if (PieceCommitmentRequest{} != dealRequest.PieceCommitment && dealRequest.PieceCommitment.Piece != "") &&
		(dealRequest.PieceCommitment.PaddedPieceSize != 0 && dealRequest.PieceCommitment.UnPaddedPieceSize != 0) &&
		(dealRequest.Size != 0) {

		// if commp is there, make sure the piece and size are there. Use default duration.
		pieceCommp.Cid = addNode.Cid().String()
		pieceCommp.Piece = dealRequest.PieceCommitment.Piece
		pieceCommp.Size = file.Size
		pieceCommp.UnPaddedPieceSize = dealRequest.PieceCommitment.UnPaddedPieceSize
		pieceCommp.PaddedPieceSize = dealRequest.PieceCommitment.PaddedPieceSize
		pieceCommp.CreatedAt = time.Now()
		pieceCommp.UpdatedAt = time.Now()
		pieceCommp.Status = utils.COMMP_STATUS_OPEN
		node.DB.Create(&pieceCommp)

		dealRequest.PieceCommitment = PieceCommitmentRequest{
			Piece:             pieceCommp.Piece,
			PaddedPieceSize:   pieceCommp.PaddedPieceSize,
			UnPaddedPieceSize: pieceCommp.UnPaddedPieceSize,
		}
	}

	// save the content to the DB with the piece_commitment_id
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
	dealRequest.Cid = content.Cid

	//	assign a miner
	if dealRequest.Miner != "" {
		contentMinerAssignment := model.ContentMiner{
			Miner:     dealRequest.Miner,
			Content:   content.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		node.DB.Create(&contentMinerAssignment)
		dealRequest.Miner = contentMinerAssignment.Miner
	}

	// 	assign a wallet_estuary
	if (WalletRequest{} != dealRequest.Wallet && dealRequest.Wallet.KeyType != "") {
		var hexedWallet WalletRequest
		hexedWallet.KeyType = dealRequest.Wallet.KeyType
		hexedWallet.PrivateKey = hex.EncodeToString([]byte(dealRequest.Wallet.PrivateKey))
		walletByteArr, err := json.Marshal(hexedWallet)

		if err != nil {
			return errors.New("Error encoding the wallet")
		}
		contentWalletAssignment := model.ContentWallet{
			Wallet:    string(walletByteArr),
			Content:   content.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		node.DB.Create(&contentWalletAssignment)
		dealRequest.Wallet = WalletRequest{
			KeyType: contentWalletAssignment.Wallet,
		}
	}

	var dealProposalParam model.ContentDealProposalParameters
	dealProposalParam.CreatedAt = time.Now()
	dealProposalParam.UpdatedAt = time.Now()
	dealProposalParam.Content = content.ID
	dealProposalParam.Label = content.Cid

	// duration
	if dealRequest.Duration == 0 {
		dealProposalParam.Duration = utils.DEFAULT_DURATION
	} else {
		dealProposalParam.Duration = dealRequest.Duration
	}

	// start epoch
	if dealRequest.StartEpoch != 0 {
		dealProposalParam.StartEpoch = dealRequest.StartEpoch
	}

	// remove unsealed copy
	if dealRequest.RemoveUnsealedCopies == false {
		dealProposalParam.RemoveUnsealedCopy = false
	} else {
		dealProposalParam.RemoveUnsealedCopy = true
	}

	// deal proposal parameters
	node.DB.Create(&dealProposalParam)

	if err != nil {
		return errors.New("Error pinning the file")
	}

	var dispatchJobs core.IProcessor
	if pieceCommp.ID != 0 {
		dispatchJobs = jobs.NewStorageDealMakerProcessor(node, content, pieceCommp) // straight to storage deal making
	} else {
		dispatchJobs = jobs.NewPieceCommpProcessor(node, content) // straight to pieceCommp
	}

	node.Dispatcher.AddJobAndDispatch(dispatchJobs, 1)

	c.JSON(200, DealResponse{
		Status:      "success",
		Message:     "File uploaded and pinned successfully",
		ContentId:   content.ID,
		DealRequest: dealRequest,
	})

	return nil
}

// handleCommPieceAdd handles the request to add a commp record.
// @Summary Add a commp record
// @Description Add a commp record
// @Tags deals
// @Accept  json
// @Produce  json
func handleCommPieceAdd(c echo.Context, node *core.DeltaNode, statsService core.StatsService) error {
	var dealRequest DealRequest

	// lets record this.
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")
	err := c.Bind(&dealRequest)

	if err != nil {
		return errors.New("Error parsing the request, please check the request body if it complies with the spec")
	}

	err = ValidateMeta(dealRequest)
	if err != nil {
		return err
	}

	// specify the connection mode
	var connMode = dealRequest.ConnectionMode
	if connMode == "" || (connMode != utils.CONNECTION_MODE_E2E && connMode != utils.CONNECTION_MODE_IMPORT) {
		connMode = "e2e"
	}

	// let's create a commp but only if we have
	// a cid, a piece_cid, a padded_piece_size, size
	var pieceCommp model.PieceCommitment
	if (PieceCommitmentRequest{} != dealRequest.PieceCommitment && dealRequest.PieceCommitment.Piece != "") &&
		(dealRequest.PieceCommitment.PaddedPieceSize != 0) &&
		(dealRequest.Size != 0) {

		// if commp is there, make sure the piece and size are there. Use default duration.
		pieceCommp.Cid = dealRequest.Cid
		pieceCommp.Piece = dealRequest.PieceCommitment.Piece
		pieceCommp.Size = dealRequest.Size
		pieceCommp.UnPaddedPieceSize = dealRequest.PieceCommitment.UnPaddedPieceSize
		pieceCommp.PaddedPieceSize = dealRequest.PieceCommitment.PaddedPieceSize
		pieceCommp.CreatedAt = time.Now()
		pieceCommp.UpdatedAt = time.Now()
		pieceCommp.Status = utils.COMMP_STATUS_OPEN
		node.DB.Create(&pieceCommp)

		dealRequest.PieceCommitment = PieceCommitmentRequest{
			Piece:             pieceCommp.Piece,
			PaddedPieceSize:   pieceCommp.PaddedPieceSize,
			UnPaddedPieceSize: pieceCommp.UnPaddedPieceSize,
		}
	}

	// save the content to the DB with the piece_commitment_id
	content := model.Content{
		Name:              dealRequest.Cid,
		Size:              dealRequest.Size,
		Cid:               dealRequest.Cid,
		RequestingApiKey:  authParts[1],
		PieceCommitmentId: pieceCommp.ID,
		Status:            utils.CONTENT_PINNED,
		ConnectionMode:    connMode,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	node.DB.Create(&content)
	dealRequest.Cid = content.Cid

	//	assign a miner
	if dealRequest.Miner != "" {
		contentMinerAssignment := model.ContentMiner{
			Miner:     dealRequest.Miner,
			Content:   content.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		node.DB.Create(&contentMinerAssignment)
		dealRequest.Miner = contentMinerAssignment.Miner
	}

	// 	assign a wallet_estuary
	if (WalletRequest{} != dealRequest.Wallet && dealRequest.Wallet.KeyType != "") {
		var hexedWallet WalletRequest
		hexedWallet.KeyType = dealRequest.Wallet.KeyType
		hexedWallet.PrivateKey = hex.EncodeToString([]byte(dealRequest.Wallet.PrivateKey))
		walletByteArr, err := json.Marshal(hexedWallet)

		if err != nil {
			return errors.New("Error parsing the wallet")
		}
		contentWalletAssignment := model.ContentWallet{
			Wallet:    string(walletByteArr),
			Content:   content.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		node.DB.Create(&contentWalletAssignment)
		dealRequest.Wallet = WalletRequest{
			KeyType: contentWalletAssignment.Wallet,
		}
	}

	var dealProposalParam model.ContentDealProposalParameters
	dealProposalParam.CreatedAt = time.Now()
	dealProposalParam.UpdatedAt = time.Now()
	dealProposalParam.Content = content.ID
	dealProposalParam.Label = content.Cid

	// duration
	if dealRequest.Duration == 0 {
		dealProposalParam.Duration = utils.DEFAULT_DURATION
	} else {
		dealProposalParam.Duration = dealRequest.Duration
	}

	// start epoch
	if dealRequest.StartEpoch != 0 {
		dealProposalParam.StartEpoch = dealRequest.StartEpoch
	}

	// remove unsealed copy
	if dealRequest.RemoveUnsealedCopies == false {
		dealProposalParam.RemoveUnsealedCopy = false
	} else {
		dealProposalParam.RemoveUnsealedCopy = true
	}

	// deal proposal parameters
	node.DB.Create(&dealProposalParam)

	if err != nil {
		return errors.New("Error parsing the request, please check the request body if it complies with the spec")
	}

	var dispatchJobs core.IProcessor
	if pieceCommp.ID != 0 {
		dispatchJobs = jobs.NewStorageDealMakerProcessor(node, content, pieceCommp) // straight to storage deal making
	}

	node.Dispatcher.AddJobAndDispatch(dispatchJobs, 1)

	c.JSON(200, DealResponse{
		Status:      "success",
		Message:     "File uploaded and pinned successfully",
		ContentId:   content.ID,
		DealRequest: dealRequest,
	})
	return nil
}

// handleCommPiecesAdd handles the request to add a commp record.
// @Summary Add a commp record
// @Description Add a commp record
// @Tags CommP
// @Accept  json
// @Produce  json
func handleCommPiecesAdd(c echo.Context, node *core.DeltaNode, statsService core.StatsService) error {
	var dealRequests []DealRequest

	// lets record this.
	authorizationString := c.Request().Header.Get("Authorization")
	authParts := strings.Split(authorizationString, " ")

	//	validate the meta
	err := c.Bind(&dealRequests)

	if err != nil {
		return errors.New("Error parsing the request, please check the request body if it complies with the spec")
	}

	var dealResponses []DealResponse
	for _, dealRequest := range dealRequests {

		err = ValidateMeta(dealRequest)
		if err != nil {
			return err
		}

		// specify the connection mode
		var connMode = dealRequest.ConnectionMode
		if connMode == "" || (connMode != utils.CONNECTION_MODE_E2E && connMode != utils.CONNECTION_MODE_IMPORT) {
			connMode = "e2e"
		}

		// let's create a commp but only if we have
		// a cid, a piece_cid, a padded_piece_size, size
		var pieceCommp model.PieceCommitment
		if (PieceCommitmentRequest{} != dealRequest.PieceCommitment && dealRequest.PieceCommitment.Piece != "") &&
			(dealRequest.PieceCommitment.PaddedPieceSize != 0) &&
			(dealRequest.Size != 0) {

			// if commp is there, make sure the piece and size are there. Use default duration.
			pieceCommp.Cid = dealRequest.Cid
			pieceCommp.Piece = dealRequest.PieceCommitment.Piece
			pieceCommp.Size = dealRequest.Size
			pieceCommp.UnPaddedPieceSize = dealRequest.PieceCommitment.UnPaddedPieceSize
			pieceCommp.PaddedPieceSize = dealRequest.PieceCommitment.PaddedPieceSize
			pieceCommp.CreatedAt = time.Now()
			pieceCommp.UpdatedAt = time.Now()
			pieceCommp.Status = utils.COMMP_STATUS_OPEN
			node.DB.Create(&pieceCommp)

			dealRequest.PieceCommitment = PieceCommitmentRequest{
				Piece:             pieceCommp.Piece,
				PaddedPieceSize:   pieceCommp.PaddedPieceSize,
				UnPaddedPieceSize: pieceCommp.UnPaddedPieceSize,
			}
		}

		// save the content to the DB with the piece_commitment_id
		content := model.Content{
			Name:              dealRequest.Cid,
			Size:              dealRequest.Size,
			Cid:               dealRequest.Cid,
			RequestingApiKey:  authParts[1],
			PieceCommitmentId: pieceCommp.ID,
			Status:            utils.CONTENT_PINNED,
			ConnectionMode:    connMode,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		node.DB.Create(&content)
		dealRequest.Cid = content.Cid

		//	assign a miner
		if dealRequest.Miner != "" {
			contentMinerAssignment := model.ContentMiner{
				Miner:     dealRequest.Miner,
				Content:   content.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			node.DB.Create(&contentMinerAssignment)
			dealRequest.Miner = contentMinerAssignment.Miner
		}

		// 	assign a wallet_estuary
		if (WalletRequest{} != dealRequest.Wallet && dealRequest.Wallet.KeyType != "") {
			var hexedWallet WalletRequest
			hexedWallet.KeyType = dealRequest.Wallet.KeyType
			hexedWallet.PrivateKey = hex.EncodeToString([]byte(dealRequest.Wallet.PrivateKey))
			walletByteArr, err := json.Marshal(hexedWallet)

			if err != nil {
				return errors.New("Wallet could not be encoded")
			}
			contentWalletAssignment := model.ContentWallet{
				Wallet:    string(walletByteArr),
				Content:   content.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			node.DB.Create(&contentWalletAssignment)
			dealRequest.Wallet = WalletRequest{
				KeyType: contentWalletAssignment.Wallet,
			}
		}

		var dealProposalParam model.ContentDealProposalParameters
		dealProposalParam.CreatedAt = time.Now()
		dealProposalParam.UpdatedAt = time.Now()
		dealProposalParam.Content = content.ID
		dealProposalParam.Label = content.Cid

		// duration
		if dealRequest.Duration == 0 {
			dealProposalParam.Duration = utils.DEFAULT_DURATION
		} else {
			dealProposalParam.Duration = dealRequest.Duration
		}

		// start epoch
		if dealRequest.StartEpoch != 0 {
			dealProposalParam.StartEpoch = dealRequest.StartEpoch
		}

		// remove unsealed copy
		if dealRequest.RemoveUnsealedCopies == false {
			dealProposalParam.RemoveUnsealedCopy = false
		} else {
			dealProposalParam.RemoveUnsealedCopy = true
		}

		// deal proposal parameters
		node.DB.Create(&dealProposalParam)

		var dispatchJobs core.IProcessor
		fmt.Println(pieceCommp.ID)
		if pieceCommp.ID != 0 {
			dispatchJobs = jobs.NewStorageDealMakerProcessor(node, content, pieceCommp) // straight to storage deal making
		}

		node.Dispatcher.AddJob(dispatchJobs)

		dealResponses = append(dealResponses, DealResponse{
			Status:      "success",
			Message:     "File uploaded and pinned successfully",
			DealRequest: dealRequest,
		})

	}
	node.Dispatcher.Start(len(dealRequests))
	c.JSON(http.StatusOK, dealResponses)

	return nil
}

func handleContentStats(c echo.Context, node *core.DeltaNode, statsService core.StatsService) error {
	contentIdParam := c.Param("contentId")
	contentId, err := strconv.Atoi(contentIdParam)
	if err != nil {
		return errors.New("Error looking up the status of the content" + err.Error())
	}

	status, err := statsService.ContentStatus(core.ContentStatsParam{
		ContentId: int64(contentId),
	})

	if err != nil {
		return errors.New("Error looking up the status of the content" + err.Error())

	}

	return c.JSON(200, status)
}

func handleCommitmentPieceStats(c echo.Context, node *core.DeltaNode, statsService core.StatsService) error {
	pieceCommitmentIdParam := c.Param("piece-commitmentId")
	pieceCommitmentId, err := strconv.Atoi(pieceCommitmentIdParam)
	if err != nil {
		return errors.New("Error looking up the status of the piece commitment" + err.Error())
	}

	status, err := statsService.PieceCommitmentStatus(core.PieceCommitmentStatsParam{
		PieceCommpId: int64(pieceCommitmentId),
	})

	if err != nil {
		return errors.New("Error looking up the status of the piece commitment" + err.Error())
	}

	return c.JSON(200, status)
}

type ValidateMetaResult struct {
	IsValid bool
	Message string
}

func ValidateMeta(dealRequest DealRequest) error {

	if (DealRequest{} == dealRequest) {
		return errors.New("invalid request")
	}
	// miner is required
	if (DealRequest{} != dealRequest && dealRequest.Miner == "") {
		return errors.New("miner is required")
	}

	if (DealRequest{} != dealRequest && (dealRequest.ConnectionMode != utils.CONNECTION_MODE_E2E && dealRequest.ConnectionMode != utils.CONNECTION_MODE_IMPORT)) {
		return errors.New("connection mode can only be e2e or import")
	}

	if (PieceCommitmentRequest{} != dealRequest.PieceCommitment && dealRequest.PieceCommitment.Piece != "") &&
		(dealRequest.PieceCommitment.PaddedPieceSize == 0 && dealRequest.PieceCommitment.UnPaddedPieceSize == 0) &&
		(dealRequest.Size == 0) {
		return errors.New("piece commitment is invalid, make sure you have the cid, piece_cid, size and padded_piece_size or unpadded_piece_size")

	}
	return nil
}

func handlePrepareContent() {

}
func handlePrepareCommitmentPiece() {

}
func handlePrepareCommitmentPieces() {

}

func handleAnnounceContent() {

}

func handleAnnounceCommitmentPiece() {

}

func handleAnnounceCommitmentPieces() {

}

//
//
//{
//"cid": "bafybeidty2dovweduzsne3kkeeg3tllvxd6nc2ifh6ztexvy4krc5pe7om",
//"wallet": {},
//"commp": {
//"piece": "baga6ea4seaqhfvwbdypebhffobtxjyp4gunwgwy2ydanlvbe6uizm5hlccxqmeq",
//"padded_piece_size": 4294967296
//},
//"connection_mode": "offline",
//"size": 2500366291,
//"online_sign":false
//}
// "{ uuid:'','hexed_unsigned_deal_proposal'}"
// "{ uuid:'','hexed_signed_deal_proposal'}"
//
///deal/prepare/content - prepare a proposal and send back the unsigned network proposal
///deal/prepare/piece-commitment - prepare a proposal and send back the unsigned network proposal
///deal/prepare/piece-commitments - prepare a proposal and send back the unsigned network proposal
//
///deal/announce/content (accepts a signed HEX)
///deal/announce/piece-commitment (accepts a signed HEX)
///deal/announce/piece-commitments (accepts a collection of signed HEX)

// commp-sign --wallet=wallet.json --proposal-meta=// "{ uuid:'','hexed_unsigned_deal_proposal'}"
