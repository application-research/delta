package api

import (
	"context"
	"delta/core"
	model "github.com/application-research/delta-db/db_models"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/labstack/echo/v4"
)

// TODO: OPTIMIZE!!
func ConfigureOpenInfoCheckRouter(e *echo.Group, node *core.DeltaNode) {

	e.GET("/info/wallet/balance/:address", func(c echo.Context) error {
		return handleOpenGetBalance(c, node)
	})

}

// It gets the balance of a wallet
func handleOpenGetBalance(c echo.Context, node *core.DeltaNode) error {
	var wallet model.Wallet
	node.DB.Model(&model.Wallet{}).Where("addr = ?", c.Param("address")).First(&wallet)

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
			"verified_client_balance": vcstatus,
		},
	})
}
