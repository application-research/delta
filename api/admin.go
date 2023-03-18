package api

import (
	"context"
	"delta/core"
	"encoding/base64"
	model "github.com/application-research/delta-db/db_models"
	"github.com/filecoin-project/go-address"
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

type ImportWalletWithHexRequest struct {
	HexKey string `json:"hex_key"`
}

// ConfigureAdminRouter It creates a new wallet and saves it to the database
// It configures the admin router
func ConfigureAdminRouter(e *echo.Group, node *core.DeltaNode) {
	adminWallet := e.Group("/wallet")
	adminWallet.POST("/register", handleAdminRegisterWallet(node))
	adminWallet.POST("/register-hex", handleAdminRegisterWalletWithHex(node))
	adminWallet.POST("/create", handleAdminCreateWallet(node))
	adminWallet.GET("/list", handleAdminListWallets(node))
	adminWallet.GET("/balance/:address", handleAdminGetBalance(node))
}

// handleAdminRegisterWallet It creates a new wallet and saves it to the database
// @Summary It creates a new wallet and saves it to the database
// @Description It creates a new wallet and saves it to the database
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param address path string true "address"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/wallet/balance/:address [post]
func handleAdminGetBalance(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		// check if the address is registered in the database
		var wallet model.Wallet
		node.DB.Model(&model.Wallet{}).Where("addr = ? and owner =?", c.Param("address"), authParts[1]).First(&wallet)

		if wallet.ID == 0 {
			return c.JSON(400, map[string]interface{}{
				"message": "wallet not found, register the wallet first",
			})
		}

		addressFromParam := c.Param("address")
		address, err := address.NewFromString(addressFromParam)
		if err != nil {
			return c.JSON(400, map[string]interface{}{
				"message": "invalid address",
			})
		}
		bigIntBalance, err := node.LotusApi.WalletBalance(context.Background(), address)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to get balance",
			})
		}

		act, err := node.LotusApi.StateGetActor(context.Background(), address, types.EmptyTSK)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to get actor",
			})
		}

		market, err := node.LotusApi.StateMarketBalance(context.Background(), address, types.EmptyTSK)
		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to get market balance",
			})
		}

		vcstatus, err := node.LotusApi.StateVerifiedClientStatus(context.Background(), address, types.EmptyTSK)

		if err != nil {
			return c.JSON(500, map[string]interface{}{
				"message": "failed to get verified client status",
			})
		}

		avail := types.BigSub(market.Escrow, market.Locked)

		return c.JSON(200, map[string]interface{}{
			"message": "success",
			"balance": map[string]interface{}{
				"account":                 address.String(),
				"wallet_balance":          types.FIL(bigIntBalance),
				"balance":                 types.FIL(act.Balance),
				"market_escrow":           types.FIL(market.Escrow),
				"market_locked":           types.FIL(market.Locked),
				"market_available":        types.FIL(avail),
				"verified_client_balance": vcstatus.Int64(),
			},
		})

	}
}

// handleAdminRegisterWallet It creates a new wallet and saves it to the database
// @Summary It creates a new wallet and saves it to the database
// @Description It creates a new wallet and saves it to the database
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param address path string true "address"
// @Param key_type path string true "key_type"
// @Param private_key path string true "private_key"
// @Success 200 {object} AddWalletRequest
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/wallet/register [post]
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

func handleAdminRegisterWalletWithHex(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")
		walletService := core.NewWalletService(node)

		if len(authParts) != 2 {
			return c.JSON(401, map[string]interface{}{
				"message": "unauthorized",
			})
		}
		var hexedKey ImportWalletWithHexRequest
		c.Bind(&hexedKey)

		importedWallet, err := walletService.ImportWithHex(hexedKey.HexKey)

		if err != nil {
			return c.JSON(400, map[string]interface{}{
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

// Creating a new wallet address and saving it to the database.
// `handleAdminRegisterWallet` is a function that takes a `DeltaNode` and returns a function that takes an `echo.Context`
// and returns an `error`
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
// `handleAdminListWallets` is a function that takes a `DeltaNode` and returns a function that takes an `echo.Context` and
// returns an `error`
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
