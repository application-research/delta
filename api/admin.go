package api

import (
	"delta/core"
	"encoding/base64"
	model "github.com/application-research/delta-db/db_models"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"strings"
	"time"
)

type AddWalletRequest struct {
	Address    string `json:"address"`
	KeyType    string `json:"key_type"`
	PrivateKey string `json:"private_key"`
}

type ImportWalletRequest struct {
	KeyType    string `json:"key_type"`
	PrivateKey string `json:"private_key"`
}

// It creates a new wallet and saves it to the database
func ConfigureAdminRouter(e *echo.Group, node *core.DeltaNode) {

	adminWallet := e.Group("/wallet")
	adminStats := e.Group("/stats")
	adminStats.GET("/miner/:minerId", handleAdminStatsMiner(node))
	adminWallet.POST("/register", handleAdminRegisterWallet(node))
	adminWallet.POST("/create", handleAdminCreateWallet(node))
	adminWallet.GET("/list", handleAdminListWallets(node))
}

// It creates a new wallet address and saves it to the database
func handleAdminCreateWallet(node *core.DeltaNode) func(c echo.Context) error {

	return func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		if len(authParts) != 2 {
			return c.JSON(401, map[string]interface{}{
				"message": "unauthorized",
			})
		}

		var createWalletParam core.CreateWalletParam
		walletService := core.NewWalletService(node)
		createWalletParam.RequestingApiKey = authParts[1]
		create, err := walletService.Create(createWalletParam)
		if err != nil {
			return err
		}

		walletUuid, err := uuid.NewUUID()
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to generate uuid",
			})
		}
		newWallet := &model.Wallet{
			UuId:       walletUuid.String(),
			Addr:       create.WalletAddress.String(),
			Owner:      authParts[1],
			KeyType:    create.Wallet.KeyType,
			PrivateKey: create.Wallet.PrivateKey,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		// save on wallet table
		node.DB.Model(&model.Wallet{}).Create(newWallet)

		return c.JSON(200, map[string]interface{}{
			"message":     "Successfully imported a wallet address. Please take note of the following information.",
			"wallet_uuid": newWallet.UuId,
			"wallet_addr": newWallet.Addr,
		})
	}
}

// Creating a new wallet address and saving it to the database.
func handleAdminRegisterWallet(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")
		walletService := core.NewWalletService(node)
		var importWalletRequest ImportWalletRequest
		c.Bind(&importWalletRequest)

		if len(authParts) != 2 {
			return c.JSON(401, map[string]interface{}{
				"message": "unauthorized",
			})
		}

		// validate, owner, keytype, address and private key are required
		if importWalletRequest.KeyType == "" || string(importWalletRequest.PrivateKey) == "" {
			return c.JSON(400, map[string]interface{}{
				"message": "key_type and private_key are required",
			})
		}
		decodedPrivateKey, err := base64.StdEncoding.DecodeString(importWalletRequest.PrivateKey)
		if err != nil {
			return c.JSON(400, map[string]interface{}{
				"message": "failed to decode private key",
				"error":   err.Error(),
			})
		}
		importedWallet, err := walletService.Import(core.ImportWalletParam{
			WalletParam: core.WalletParam{
				RequestingApiKey: authParts[1],
			},
			KeyType:    types.KeyType(importWalletRequest.KeyType),
			PrivateKey: decodedPrivateKey,
		})
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to import wallet",
				"error":   err.Error(),
			})
		}

		return c.JSON(200, map[string]interface{}{
			"message":     "Successfully imported a wallet address. Please take note of the following information.",
			"wallet_addr": importedWallet.WalletAddress.String(),
			"wallet_uuid": importedWallet.Wallet.UuId,
		})

	}
}

// It takes the authorization header from the request, splits it into two parts, and then uses the second part to find all
// wallets owned by the user
func handleAdminListWallets(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		if len(authParts) != 2 {
			return c.JSON(401, map[string]interface{}{
				"message": "unauthorized",
			})
		}

		var wallets []model.Wallet
		node.DB.Where("owner = ?", authParts[1]).Find(&wallets)

		return c.JSON(200, map[string]interface{}{
			"wallets": wallets,
		})
	}
}

// A function that returns a function that returns an error.
func handleAdminStatsMiner(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {

		var contents []model.Content
		node.DB.Raw("select c.* from content_deals cd, contents c where cd.content = c.id and cd.miner = ?", c.Param("minerId")).Scan(&contents)

		var contentMinerAssignment []model.ContentMiner
		node.DB.Raw("select cma.* from content_miners cma, contents c where cma.content = c.id and cma.miner = ?", c.Param("minerId")).Scan(&contentMinerAssignment)

		return c.JSON(200, map[string]interface{}{
			"content": contents,
			"cmas":    contentMinerAssignment,
		})
	}
}
