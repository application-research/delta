package core

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	model "github.com/application-research/delta-db/db_models"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/wallet"
	"github.com/google/uuid"
	"time"
)

type WalletParam struct {
	RequestingApiKey string
}

type GetWalletParam struct {
	WalletParam
	Address string
}

type RemoveWalletParam struct {
	WalletParam
	Address string
}

type CreateWalletParam struct {
	WalletParam
	KeyType types.KeyType
}

type ImportWithHexKey struct {
	KeyType    types.KeyType `json:"Type"`
	PrivateKey string        `json:"PrivateKey"`
}

type ImportWalletParam struct {
	WalletParam
	KeyType    types.KeyType `json:"key_type"`
	PrivateKey []byte
}

type AddWalletResult struct {
	Wallet        model.Wallet
	WalletAddress address.Address
}

type DeleteWalletResult struct {
	WalletAddress string
	Message       string
}

type WalletResult struct {
	Wallet model.Wallet
}

type ImportWalletResult struct {
	Wallet        model.Wallet
	WalletAddress address.Address
}

type WalletService struct {
	Context   context.Context
	DeltaNode *DeltaNode
}

// NewWalletService Creating a new wallet service.
func NewWalletService(dn *DeltaNode) *WalletService {
	return &WalletService{
		DeltaNode: dn,
	}
}

// Create Creating a new wallet and saving it to the database.
func (w WalletService) Create(param CreateWalletParam) (AddWalletResult, error) {
	newWallet, err := wallet.NewWallet(wallet.NewMemKeyStore())
	if err != nil {
		return AddWalletResult{}, err
	}
	address, err := newWallet.WalletNew(w.Context, param.KeyType)

	if err != nil {
		return AddWalletResult{}, err
	}

	// save it on the DB
	hexedKey := hex.EncodeToString(address.Payload())
	walletToDb := &model.Wallet{
		Addr:       address.String(),
		Owner:      param.RequestingApiKey,
		KeyType:    string(param.KeyType),
		PrivateKey: hexedKey,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}
	w.DeltaNode.DB.Create(walletToDb)

	return AddWalletResult{
		Wallet:        *walletToDb,
		WalletAddress: address,
	}, nil
}

func (w WalletService) ImportWithHex(hexKey string) (ImportWalletResult, error) {
	fmt.Println(hexKey)
	hexString, err := hex.DecodeString(hexKey)
	if err != nil {
		panic(err)
	}
	var importWithHexKey ImportWithHexKey
	json.Unmarshal(hexString, &importWithHexKey)
	fmt.Println(importWithHexKey.KeyType)
	bKey, err := base64.StdEncoding.DecodeString(importWithHexKey.PrivateKey)
	if err != nil {
		panic(err)

	}
	result, err := w.Import(ImportWalletParam{
		KeyType:    importWithHexKey.KeyType,
		PrivateKey: bKey,
	})

	return result, err

}

// Import Importing a wallet.
func (w WalletService) Import(param ImportWalletParam) (ImportWalletResult, error) {
	newWallet, err := wallet.NewWallet(wallet.NewMemKeyStore())
	if err != nil {
		return ImportWalletResult{}, err
	}

	hexedWallet := hex.EncodeToString(param.PrivateKey)

	address, err := newWallet.WalletImport(w.Context, &types.KeyInfo{
		Type:       param.KeyType,
		PrivateKey: param.PrivateKey,
	})
	if err != nil {
		return ImportWalletResult{}, err
	}

	// save it on the DB
	walletUuid, err := uuid.NewUUID()
	if err != nil {
		return ImportWalletResult{}, err
	}
	walletToDb := &model.Wallet{
		UuId:       walletUuid.String(),
		Addr:       address.String(),
		Owner:      param.RequestingApiKey,
		KeyType:    string(param.KeyType),
		PrivateKey: hexedWallet,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}
	w.DeltaNode.DB.Create(walletToDb)

	return ImportWalletResult{
		Wallet:        *walletToDb,
		WalletAddress: address,
	}, nil
}

// Remove Deleting the wallet from the database.
func (w WalletService) Remove(param RemoveWalletParam) (DeleteWalletResult, error) {
	err := w.DeltaNode.DB.Delete(&model.Wallet{}).Where("owner = ? and addr = ?", param.RequestingApiKey, param.Address).Error
	if err != nil {
		return DeleteWalletResult{
			Message:       "Wallet not found",
			WalletAddress: param.Address,
		}, err
	}
	return DeleteWalletResult{
		Message:       "Wallet deleted",
		WalletAddress: param.Address,
	}, nil
}

// List A function that takes a WalletParam and returns a list of model.Wallet and an error.
func (w WalletService) List(param WalletParam) ([]model.Wallet, error) {
	var wallets []model.Wallet
	w.DeltaNode.DB.Model(&model.Wallet{}).Where("owner = ?", param.RequestingApiKey).Find(&wallets)
	return wallets, nil
	// Getting the wallet from the database.
}

// Getting the wallet from the database.
func (w WalletService) Get(param GetWalletParam) (model.Wallet, error) {
	var wallet model.Wallet
	w.DeltaNode.DB.Model(&model.Wallet{}).Where("owner = ? and addr = ?", param.RequestingApiKey, param.Address).Find(&wallet)
	return wallet, nil
}
