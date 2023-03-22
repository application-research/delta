package core

import (
	"reflect"
	"testing"
)

func TestNewStatsStatsService(t *testing.T) {
	type args struct {
		deltaNode *DeltaNode
	}
	tests := []struct {
		name string
		args args
		want *StatsService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStatsStatsService(tt.args.deltaNode); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStatsStatsService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatsService_ContentStatus(t *testing.T) {
	type fields struct {
		DeltaNode *DeltaNode
	}
	type args struct {
		param ContentStatsParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    StatsContentResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StatsService{
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := s.ContentStatus(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContentStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatsService_DealStatus(t *testing.T) {
	type fields struct {
		DeltaNode *DeltaNode
	}
	type args struct {
		param DealStatsParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    StatsDealResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StatsService{
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := s.DealStatus(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("DealStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DealStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatsService_PieceCommitmentStatus(t *testing.T) {
	type fields struct {
		DeltaNode *DeltaNode
	}
	type args struct {
		param PieceCommitmentStatsParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    StatsPieceCommitmentResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StatsService{
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := s.PieceCommitmentStatus(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("PieceCommitmentStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PieceCommitmentStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatsService_Status(t *testing.T) {
	type fields struct {
		DeltaNode *DeltaNode
	}
	type args struct {
		param StatsParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    StatsResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StatsService{
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := s.Status(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Status() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Status() got = %v, want %v", got, tt.want)
			}
		})
	}
}
