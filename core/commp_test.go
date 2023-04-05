package core

import (
	"context"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/filecoin-project/go-commp-utils/writer"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	carv2 "github.com/ipld/go-car/v2"
)

func TestCommpService_GenerateCommPCarV2(t *testing.T) {
	type fields struct {
		DeltaNode *DeltaNode
	}
	type args struct {
		readerFromFile io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *abi.PieceInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CommpService{
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := c.GenerateCommPCarV2(tt.args.readerFromFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateCommPCarV2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateCommPCarV2() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommpService_GenerateCommPFile(t *testing.T) {
	type fields struct {
		DeltaNode *DeltaNode
	}
	type args struct {
		context    context.Context
		payloadCid cid.Cid
		blockstore blockstore.Blockstore
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantPieceCid          cid.Cid
		wantPayloadSize       uint64
		wantUnPaddedPieceSize abi.UnpaddedPieceSize
		wantErr               bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CommpService{
				DeltaNode: tt.fields.DeltaNode,
			}
			gotPieceCid, gotPayloadSize, gotUnPaddedPieceSize, err := c.GenerateCommPFile(tt.args.context, tt.args.payloadCid, tt.args.blockstore)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateCommPFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPieceCid, tt.wantPieceCid) {
				t.Errorf("GenerateCommPFile() gotPieceCid = %v, want %v", gotPieceCid, tt.wantPieceCid)
			}
			if gotPayloadSize != tt.wantPayloadSize {
				t.Errorf("GenerateCommPFile() gotPayloadSize = %v, want %v", gotPayloadSize, tt.wantPayloadSize)
			}
			if gotUnPaddedPieceSize != tt.wantUnPaddedPieceSize {
				t.Errorf("GenerateCommPFile() gotUnPaddedPieceSize = %v, want %v", gotUnPaddedPieceSize, tt.wantUnPaddedPieceSize)
			}
		})
	}
}

func TestCommpService_GenerateParallelCommp(t *testing.T) {
	type fields struct {
		DeltaNode *DeltaNode
	}
	type args struct {
		readerFromFile *os.File
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    writer.DataCIDSize
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CommpService{
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := c.GenerateCommp(tt.args.readerFromFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateCommp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateCommp() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommpService_GetCarSize(t *testing.T) {
	type fields struct {
		DeltaNode *DeltaNode
	}
	type args struct {
		stream io.Reader
		rd     *carv2.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CommpService{
				DeltaNode: tt.fields.DeltaNode,
			}
			got, err := c.GetCarSize(tt.args.stream, tt.args.rd)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCarSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCarSize() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommpService_GetSize(t *testing.T) {
	type fields struct {
		DeltaNode *DeltaNode
	}
	type args struct {
		stream io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CommpService{
				DeltaNode: tt.fields.DeltaNode,
			}
			if got := c.GetSize(tt.args.stream); got != tt.want {
				t.Errorf("GetSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
