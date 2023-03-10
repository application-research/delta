package core

import (
	"bytes"
	"context"
	"fmt"
	"github.com/application-research/filclient"
	"github.com/filecoin-project/go-commp-utils/nonffi"
	"github.com/filecoin-project/go-commp-utils/writer"
	"github.com/filecoin-project/go-commp-utils/zerocomm"
	commcid "github.com/filecoin-project/go-fil-commcid"
	commp "github.com/filecoin-project/go-fil-commp-hashhash"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	carv2 "github.com/ipld/go-car/v2"
	"github.com/labstack/gommon/log"
	"golang.org/x/xerrors"
	"io"
	"math/bits"
	"runtime"
)

type CommpService struct {
	DeltaNode *DeltaNode
}

type DataCIDSize struct {
	PayloadSize int64
	PieceSize   abi.PaddedPieceSize
	PieceCID    cid.Cid
}

const commPBufPad = abi.PaddedPieceSize(8 << 20)
const CommPBuf = abi.UnpaddedPieceSize(commPBufPad - (commPBufPad / 128)) // can't use .Unpadded() for const

type ciderr struct {
	c   cid.Cid
	err error
}

// GenerateCommPFile Generating a CommP file from a payload file.
// Generating a CommP file from a payload file.
func (c CommpService) GenerateCommPFile(context context.Context, payloadCid cid.Cid, blockstore blockstore.Blockstore) (pieceCid cid.Cid, payloadSize uint64, unPaddedPieceSize abi.UnpaddedPieceSize, err error) {
	return filclient.GeneratePieceCommitment(context, payloadCid, blockstore)
}

// GenerateCommPCarV2 Generating a CommP file from a CARv2 file.
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

// Generating a CommP file from a CARv2 file.
func (c CommpService) GenerateParallelCommp(readerFromFile io.Reader) (DataCIDSize, error) {
	bytesFromCar, err := io.ReadAll(readerFromFile)
	if err != nil {
		return DataCIDSize{}, err
	}

	writer := &DataCidWriter{}
	writer.Write(bytesFromCar)
	return writer.Sum()
}

// GetSize Getting the size of the file.
// Getting the size of the file.
func (c CommpService) GetSize(stream io.Reader) int {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Len()
}

// GetCarSize Getting the size of the CARv2 file.
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

type DataCidWriter struct {
	len    int64
	buf    [CommPBuf]byte
	leaves []chan ciderr

	tbufs    [][CommPBuf]byte
	throttle chan int
}

func (w *DataCidWriter) Write(p []byte) (int, error) {
	if w.throttle == nil {
		w.throttle = make(chan int, runtime.NumCPU())
		for i := 0; i < cap(w.throttle); i++ {
			w.throttle <- i
		}
	}
	if w.tbufs == nil {
		w.tbufs = make([][CommPBuf]byte, cap(w.throttle))
	}

	n := len(p)
	for len(p) > 0 {
		buffered := int(w.len % int64(len(w.buf)))
		toBuffer := len(w.buf) - buffered
		if toBuffer > len(p) {
			toBuffer = len(p)
		}

		copied := copy(w.buf[buffered:], p[:toBuffer])
		p = p[copied:]
		w.len += int64(copied)

		if copied > 0 && w.len%int64(len(w.buf)) == 0 {
			leaf := make(chan ciderr, 1)
			bufIdx := <-w.throttle
			copy(w.tbufs[bufIdx][:], w.buf[:])

			go func() {
				defer func() {
					w.throttle <- bufIdx
				}()

				cc := new(commp.Calc)
				_, _ = cc.Write(w.tbufs[bufIdx][:])
				p, _, _ := cc.Digest()
				l, _ := commcid.PieceCommitmentV1ToCID(p)
				leaf <- ciderr{
					c:   l,
					err: nil,
				}
			}()

			w.leaves = append(w.leaves, leaf)
		}
	}
	return n, nil
}

func (w *DataCidWriter) Sum() (DataCIDSize, error) {
	// process last non-zero leaf if exists
	lastLen := w.len % int64(len(w.buf))
	rawLen := w.len

	leaves := make([]cid.Cid, len(w.leaves))
	for i, leaf := range w.leaves {
		r := <-leaf
		if r.err != nil {
			return DataCIDSize{}, xerrors.Errorf("processing leaf %d: %w", i, r.err)
		}
		leaves[i] = r.c
	}

	// process remaining bit of data
	if lastLen != 0 {
		if len(leaves) != 0 {
			copy(w.buf[lastLen:], make([]byte, int(int64(CommPBuf)-lastLen)))
			lastLen = int64(CommPBuf)
		}

		cc := new(commp.Calc)
		_, _ = cc.Write(w.buf[:lastLen])
		pb, pps, _ := cc.Digest()
		p, _ := commcid.PieceCommitmentV1ToCID(pb)

		if abi.PaddedPieceSize(pps).Unpadded() < CommPBuf { // special case for pieces smaller than 16MiB
			return DataCIDSize{
				PayloadSize: w.len,
				PieceSize:   abi.PaddedPieceSize(pps),
				PieceCID:    p,
			}, nil
		}

		leaves = append(leaves, p)
	}

	// pad with zero pieces to power-of-two size
	fillerLeaves := (1 << (bits.Len(uint(len(leaves) - 1)))) - len(leaves)
	for i := 0; i < fillerLeaves; i++ {
		leaves = append(leaves, zerocomm.ZeroPieceCommitment(CommPBuf))
	}

	if len(leaves) == 1 {
		return DataCIDSize{
			PayloadSize: rawLen,
			PieceSize:   abi.PaddedPieceSize(len(leaves)) * commPBufPad,
			PieceCID:    leaves[0],
		}, nil
	}

	pieces := make([]abi.PieceInfo, len(leaves))
	for i, leaf := range leaves {
		pieces[i] = abi.PieceInfo{
			Size:     commPBufPad,
			PieceCID: leaf,
		}
	}

	p, err := nonffi.GenerateUnsealedCID(abi.RegisteredSealProof_StackedDrg32GiBV1, pieces)
	if err != nil {
		return DataCIDSize{}, xerrors.Errorf("generating unsealed CID: %w", err)
	}

	return DataCIDSize{
		PayloadSize: rawLen,
		PieceSize:   abi.PaddedPieceSize(len(leaves)) * commPBufPad,
		PieceCID:    p,
	}, nil
}
