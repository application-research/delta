package core

import (
	"bytes"
	"context"
	"fmt"
	"github.com/application-research/filclient"
	"github.com/filecoin-project/go-commp-utils/writer"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	carv2 "github.com/ipld/go-car/v2"
	"github.com/labstack/gommon/log"
	"io"
)

type CommpService struct {
	DeltaNode *DeltaNode
}

// Generating a CommP file from a payload file.
func (c CommpService) GenerateCommPFile(context context.Context, payloadCid cid.Cid, blockstore blockstore.Blockstore) (pieceCid cid.Cid, payloadSize uint64, unPaddedPieceSize abi.UnpaddedPieceSize, err error) {
	return filclient.GeneratePieceCommitment(context, payloadCid, blockstore)
}

// Generating a CommP file from a CARv2 file.
func (c CommpService) GenerateCommPCarV2(readerFromFile io.Reader) (*abi.PieceInfo, error) {
	bytesFromCar, err := io.ReadAll(readerFromFile)
	rd, err := carv2.NewReader(bytes.NewReader(bytesFromCar))
	if err != nil {
		return nil, fmt.Errorf("failed to get CARv2 reader: %w", err)
	}

	defer func() {
		if err := rd.Close(); err != nil {
			log.Warnf("failed to close CARv2 reader: %w", err)
		}
	}()

	// dump the CARv1 payload of the CARv2 file to the Commp Writer and get back the CommP.
	w := &writer.Writer{}
	r, err := rd.DataReader()
	if err != nil {
		return nil, fmt.Errorf("getting data reader for CAR v1 from CAR v2: %w", err)
	}

	written, err := io.Copy(w, r)
	if err != nil {
		return nil, fmt.Errorf("writing to commp writer: %w", err)
	}
	size, err := c.GetCarSize(readerFromFile, rd)
	if err != nil {
		return nil, err
	}
	if size == 0 {
		size = int64(len(bytesFromCar))
	}

	if written != size {
		return nil, fmt.Errorf("number of bytes written to CommP writer %d not equal to the CARv1 payload size %d", written, rd.Header.DataSize)
	}

	pi, err := w.Sum()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate CommP: %w", err)
	}

	return &abi.PieceInfo{
		Size:     pi.PieceSize,
		PieceCID: pi.PieceCID,
	}, nil
}

// Getting the size of the file.
func (c CommpService) GetSize(stream io.Reader) int {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Len()
}

// Getting the size of the CARv2 file.
func (c CommpService) GetCarSize(stream io.Reader, rd *carv2.Reader) (int64, error) {
	var size int64
	switch rd.Version {
	case 2:
		size = int64(rd.Header.DataSize)
	case 1:
		bytes, err := io.ReadAll(stream)
		if err != nil {
			return 0, err
		}
		size = int64(len(bytes))
	}
	return size, nil
}
