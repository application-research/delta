package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/bits"
	"runtime"
	"time"

	"github.com/filecoin-project/go-commp-utils/nonffi"
	"github.com/filecoin-project/go-commp-utils/writer"
	"github.com/filecoin-project/go-commp-utils/zerocomm"
	commcid "github.com/filecoin-project/go-fil-commcid"
	commp "github.com/filecoin-project/go-fil-commp-hashhash"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

const commPBufPad = abi.PaddedPieceSize(8 << 20)
const CommPBuf = abi.UnpaddedPieceSize(commPBufPad - (commPBufPad / 128)) // can't use .Unpadded() for const

type ciderr struct {
	c   cid.Cid
	err error
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

	}
	for i := 0; i < cap(w.throttle); i++ {
		w.throttle <- i
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

func (w *DataCidWriter) Sum() (writer.DataCIDSize, error) {
	// process last non-zero leaf if exists
	lastLen := w.len % int64(len(w.buf))
	rawLen := w.len

	leaves := make([]cid.Cid, len(w.leaves))
	for i, leaf := range w.leaves {
		r := <-leaf
		if r.err != nil {
			return writer.DataCIDSize{}, xerrors.Errorf("processing leaf %d: %w", i, r.err)
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
			return writer.DataCIDSize{
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
		return writer.DataCIDSize{
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
		return writer.DataCIDSize{}, xerrors.Errorf("generating unsealed CID: %w", err)
	}

	return writer.DataCIDSize{
		PayloadSize: rawLen,
		PieceSize:   abi.PaddedPieceSize(len(leaves)) * commPBufPad,
		PieceCID:    p,
	}, nil
}

// fastCommp calculates the commp of a CARv1 or CARv2 file
func fastCommp(reader io.ReadSeekCloser) (writer.DataCIDSize, error) {
	start := time.Now()
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return writer.DataCIDSize{}, err
	}

	elapsed := time.Since(start)
	fmt.Printf("Elapsed load time: %s\n", elapsed)

	start = time.Now()
	cc := new(DataCidWriter)
	cc.Write(data)
	sum, err := cc.Sum()
	if err != nil {
		panic(err)
	}

	elapsed = time.Since(start)
	fmt.Printf("commP: %s\n", sum.PieceCID.String())
	fmt.Printf("Elapsed commP time: %s\n", elapsed)

	return sum, nil
}
