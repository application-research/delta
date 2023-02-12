package core

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/wallet"
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

type ImportWalletParam struct {
	WalletParam
	KeyType    types.KeyType
	PrivateKey []byte
}

type AddWalletResult struct {
	Wallet        Wallet
	WalletAddress address.Address
}

type DeleteWalletResult struct {
	WalletAddress string
	Message       string
}

type WalletResult struct {
	Wallet Wallet
}

type ImportWalletResult struct {
	Wallet        Wallet
	WalletAddress address.Address
}

type WalletService struct {
	Context   context.Context
	DeltaNode DeltaNode
}

func NewWalletService(dn DeltaNode) *WalletService {
	return &WalletService{
		DeltaNode: dn,
	}
}

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
	walletToDb := &Wallet{
		Addr:       address.String(),
		Owner:      param.RequestingApiKey,
		KeyType:    string(param.KeyType),
		PrivateKey: string(address.Payload()),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}
	w.DeltaNode.DB.Create(walletToDb)

	return AddWalletResult{
		Wallet:        *walletToDb,
		WalletAddress: address,
	}, nil
}

func (w WalletService) Import(param ImportWalletParam) (ImportWalletResult, error) {
	newWallet, err := wallet.NewWallet(wallet.NewMemKeyStore())
	if err != nil {
		return ImportWalletResult{}, err
	}
	address, err := newWallet.WalletImport(w.Context, &types.KeyInfo{
		Type:       param.KeyType,
		PrivateKey: param.PrivateKey,
	})
	if err != nil {
		return ImportWalletResult{}, err
	}

	// save it on the DB
	walletToDb := &Wallet{
		Addr:       address.String(),
		Owner:      param.RequestingApiKey,
		KeyType:    string(param.KeyType),
		PrivateKey: string(param.PrivateKey),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}
	w.DeltaNode.DB.Create(walletToDb)

	return ImportWalletResult{
		Wallet:        *walletToDb,
		WalletAddress: address,
	}, nil
}

func (w WalletService) Remove(param RemoveWalletParam) (DeleteWalletResult, error) {
	err := w.DeltaNode.DB.Delete(&Wallet{}).Where("owner = ? and addr = ?", param.RequestingApiKey, param.Address).Error
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

func (w WalletService) List(param WalletParam) ([]Wallet, error) {
	var wallets []Wallet
	w.DeltaNode.DB.Model(&Wallet{}).Where("owner = ?", param.RequestingApiKey).Find(&wallets)
	return wallets, nil
}

func (w WalletService) Get(param GetWalletParam) (Wallet, error) {
	var wallet Wallet
	w.DeltaNode.DB.Model(&Wallet{}).Where("owner = ? and addr = ?", param.RequestingApiKey, param.Address).Find(&wallet)
	return wallet, nil
}
