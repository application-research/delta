package core

import (
	"context"
	"delta/models"
	"reflect"
	"testing"
)

func TestNewWalletService(t *testing.T) {
	type args struct {
		dn *DeltaNode
	}
	tests := []struct {
		name string
		args args
		want *WalletService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWalletService(tt.args.dn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWalletService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_Create(t *testing.T) {
	type fields struct {
		Context   context.Context
		DeltaNode *DeltaNode
	}
	type args struct {
		param CreateWalletParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    AddWalletResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WalletService{
				Context:   tt.fields.Context,
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := w.Create(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_Get(t *testing.T) {
	type fields struct {
		Context   context.Context
		DeltaNode *DeltaNode
	}
	type args struct {
		param GetWalletParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    db_models.Wallet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WalletService{
				Context:   tt.fields.Context,
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := w.Get(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_GetTokenHash(t *testing.T) {
	type fields struct {
		Context   context.Context
		DeltaNode *DeltaNode
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WalletService{
				Context:   tt.fields.Context,
				DeltaNode: tt.fields.DeltaNode,
			}
			if got := w.GetTokenHash(tt.args.key); got != tt.want {
				t.Errorf("GetTokenHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_Import(t *testing.T) {
	type fields struct {
		Context   context.Context
		DeltaNode *DeltaNode
	}
	type args struct {
		param ImportWalletParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ImportWalletResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WalletService{
				Context:   tt.fields.Context,
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := w.Import(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Import() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Import() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_ImportWithHex(t *testing.T) {
	type fields struct {
		Context   context.Context
		DeltaNode *DeltaNode
	}
	type args struct {
		hexKey string
		auth   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ImportWalletResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WalletService{
				Context:   tt.fields.Context,
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := w.ImportWithHex(tt.args.hexKey, tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImportWithHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImportWithHex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_List(t *testing.T) {
	type fields struct {
		Context   context.Context
		DeltaNode *DeltaNode
	}
	type args struct {
		param WalletParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []db_models.Wallet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WalletService{
				Context:   tt.fields.Context,
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := w.List(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_Remove(t *testing.T) {
	type fields struct {
		Context   context.Context
		DeltaNode *DeltaNode
	}
	type args struct {
		param RemoveWalletParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DeleteWalletResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WalletService{
				Context:   tt.fields.Context,
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := w.Remove(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remove() got = %v, want %v", got, tt.want)
			}
		})
	}
}
