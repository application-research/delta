package core

import (
	model "github.com/application-research/delta-db/db_models"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/filecoin-project/lotus/chain/wallet"
	"reflect"
	"testing"
)

func TestCleanUpContentAndPieceComm(t *testing.T) {
	type args struct {
		ln *DeltaNode
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanUpContentAndPieceComm(tt.args.ln)
		})
	}
}

func TestGetHostname(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHostname(); got != tt.want {
				t.Errorf("GetHostname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPublicIP(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAnnounceAddrIP()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnnounceAddrIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAnnounceAddrIP() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLotusConnection(t *testing.T) {
	type args struct {
		fullNodeApiInfo string
	}
	tests := []struct {
		name    string
		args    args
		want    v1api.FullNode
		want1   jsonrpc.ClientCloser
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := LotusConnection(tt.args.fullNodeApiInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("LotusConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LotusConnection() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LotusConnection() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewLightNode(t *testing.T) {
	type args struct {
		repo NewLightNodeParams
	}
	tests := []struct {
		name    string
		args    args
		want    *DeltaNode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLightNode(tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLightNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLightNode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanHostComputeResources(t *testing.T) {
	type args struct {
		ln   *DeltaNode
		repo string
	}
	tests := []struct {
		name string
		args args
		want *model.InstanceMeta
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ScanHostComputeResources(tt.args.ln, tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScanHostComputeResources() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetupWallet(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    *wallet.LocalWallet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetupWallet(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetupWallet() got = %v, want %v", got, tt.want)
			}
		})
	}
}
