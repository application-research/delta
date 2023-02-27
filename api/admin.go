package api

import (
	"delta/core"
	"encoding/hex"
	model "github.com/application-research/delta-db/db_models"
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

// It creates a new wallet and saves it to the database
func ConfigureAdminRouter(e *echo.Group, node *core.DeltaNode) {

	adminWallet := e.Group("/wallet")
	adminStats := e.Group("/stats")
	adminStats.GET("/miner/:minerId", handleAdminStatsMiner(node))
	adminWallet.POST("/register", handleAdminRegisterWallet(node))
	adminWallet.POST("/create", handleAdminCreateWallet(node))
}

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
			"message":     "Successfully imported a wallet address. Please take note of the UUID.",
			"wallet_uuid": newWallet.UuId,
		})
	}
}

func handleAdminRegisterWallet(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")
		var addWalletRequest AddWalletRequest
		c.Bind(&addWalletRequest)

		hexedWallet := hex.EncodeToString([]byte(addWalletRequest.PrivateKey))

		if len(authParts) != 2 {
			return c.JSON(401, map[string]interface{}{
				"message": "unauthorized",
			})
		}

		// validate, owner, keytype, address and private key are required
		if addWalletRequest.Address == "" || addWalletRequest.KeyType == "" || addWalletRequest.PrivateKey == "" {
			return c.JSON(400, map[string]interface{}{
				"message": "address, key_type and private_key are required",
			})
		}

		walletUuid, err := uuid.NewUUID()
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to generate uuid",
			})
		}
		newWallet := &model.Wallet{
			UuId:       walletUuid.String(),
			Addr:       addWalletRequest.Address,
			Owner:      authParts[1],
			KeyType:    addWalletRequest.KeyType,
			PrivateKey: hexedWallet,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		// save on wallet table
		node.DB.Model(&model.Wallet{}).Create(newWallet)

		return c.JSON(200, map[string]interface{}{
			"message":     "Successfully imported a wallet address. Please take note of the UUID.",
			"wallet_addr": addWalletRequest.Address,
			"wallet_uuid": newWallet.UuId,
		})

	}
}

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
