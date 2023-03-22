package core

import (
	"reflect"
	"testing"
)

func TestNewRepairService(t *testing.T) {
	type args struct {
		dn DeltaNode
	}
	tests := []struct {
		name string
		args args
		want *RepairService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRepairService(tt.args.dn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRepairService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepairService_RecreateCommp(t *testing.T) {
	type fields struct {
		DeltaNode DeltaNode
	}
	type args struct {
		param RepairParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    RepairResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RepairService{
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := r.RecreateCommp(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecreateCommp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecreateCommp() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepairService_RecreateContent(t *testing.T) {
	type fields struct {
		DeltaNode DeltaNode
	}
	type args struct {
		param RepairParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    RepairResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RepairService{
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := r.RecreateContent(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecreateContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecreateContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepairService_RecreateDeal(t *testing.T) {
	type fields struct {
		DeltaNode DeltaNode
	}
	type args struct {
		param RepairParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    RepairResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RepairService{
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := r.RecreateDeal(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecreateDeal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecreateDeal() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepairService_RestartDataTransfer(t *testing.T) {
	type fields struct {
		DeltaNode DeltaNode
	}
	type args struct {
		param RepairParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    RepairResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RepairService{
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := r.RestartDataTransfer(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("RestartDataTransfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RestartDataTransfer() got = %v, want %v", got, tt.want)
			}
		})
	}
}
