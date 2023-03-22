package core

import "testing"

func TestSetDataTransferEventsSubscribe(t *testing.T) {
	type args struct {
		i *DeltaNode
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDataTransferEventsSubscribe(tt.args.i)
		})
	}
}

func TestSetLibp2pManagerSubscribe(t *testing.T) {
	type args struct {
		i *DeltaNode
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLibp2pManagerSubscribe(tt.args.i)
		})
	}
}
