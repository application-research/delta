package api

import (
	"delta/core"
	"encoding/hex"
	model "github.com/application-research/delta-db/db_models"
	"github.com/labstack/echo/v4"
	"strings"
	"time"
)

// {
//   "address":"",
//   "key_type":"",
//   "private_key:"",
//

type AddWalletRequest struct {
	Address    string `json:"address"`
	KeyType    string `json:"key_type"`
	PrivateKey string `json:"private_key"`
}

// ConfigureAdminRouter configures the admin router
// This is the router that is used to administer the node
func ConfigureAdminRouter(e *echo.Group, node *core.DeltaNode) {

	//walletService := core.NewWalletService(node)

	adminRepair := e.Group("/repair")
	adminWallet := e.Group("/wallet")
	adminDashboard := e.Group("/dashboard")
	adminStats := e.Group("/stats")

	adminStats.GET("/miner/:minerId", func(c echo.Context) error {

		var contents []model.Content
		node.DB.Raw("select c.* from content_deals cd, contents c where cd.content = c.id and cd.miner = ?", c.Param("minerId")).Scan(&contents)

		var contentMinerAssignment []model.ContentMiner
		node.DB.Raw("select cma.* from content_miners cma, contents c where cma.content = c.id and cma.miner = ?", c.Param("minerId")).Scan(&contentMinerAssignment)

		return c.JSON(200, map[string]interface{}{
			"content": contents,
			"cmas":    contentMinerAssignment,
		})
	})

	// repair endpoints
	adminRepair.GET("/deal", func(c echo.Context) error {
		return nil
	})

	adminRepair.GET("/commp", func(c echo.Context) error {
		return nil
	})

	adminRepair.GET("/run-cleanup", func(c echo.Context) error {
		return nil
	})

	adminRepair.GET("/retry-deal-making-content", func(c echo.Context) error {
		return nil
	})

	// import wallet_estuary endpoint
	adminWallet.POST("/import", func(c echo.Context) error {
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

		newWallet := &model.Wallet{
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
			"message":   "Successfully imported a wallet address",
			"wallet_id": newWallet.ID,
		})

	})

	adminWallet.POST("/add", func(c echo.Context) error {
		return nil
	})

	// list wallet_estuary endpoint
	adminWallet.GET("/list", func(c echo.Context) error {
		return nil
	})

	adminWallet.GET("/info", func(c echo.Context) error {
		return nil
	})

	adminDashboard.GET("/index", func(c echo.Context) error {
		return nil
	})
}
