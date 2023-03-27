package core

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/filecoin-project/go-commp-utils/writer"
	commcid "github.com/filecoin-project/go-fil-commcid"
	commp "github.com/filecoin-project/go-fil-commp-hashhash"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
)

type CarHeader struct {
	Roots   []cid.Cid
	Version uint64
}

func init() {
	cbor.RegisterCborType(CarHeader{})
}

const BufSize = (4 << 20) / 128 * 127

func process(streamBuf *bufio.Reader, streamLen int64) (strLen int64, err error) {
	for {
		maybeNextFrameLen, err := streamBuf.Peek(10)
		if err == io.EOF {
			break
		}
		if err != nil && err != bufio.ErrBufferFull {
			log.Fatalf("unexpected error at offset %d: %s", streamLen, err)
		}
		if len(maybeNextFrameLen) == 0 {
			log.Fatalf("impossible 0-length peek without io.EOF at offset %d", streamLen)
		}

		frameLen, viLen := binary.Uvarint(maybeNextFrameLen)
		if viLen <= 0 {
			// car file with trailing garbage behind it
			return streamLen, fmt.Errorf("aborting car stream parse: undecodeable varint at offset %d", streamLen)
		}
		if frameLen > 2<<20 {
			// anything over ~2MiB got to be a mistake
			return streamLen, fmt.Errorf("aborting car stream parse: unexpectedly large frame length of %d bytes at offset %d", frameLen, streamLen)
		}

		actualFrameLen, err := io.CopyN(io.Discard, streamBuf, int64(viLen)+int64(frameLen))
		streamLen += actualFrameLen
		if err != nil {
			if err != io.EOF {
				log.Fatalf("unexpected error at offset %d: %s", streamLen-actualFrameLen, err)
			}
			return streamLen, fmt.Errorf("aborting car stream parse: truncated frame at offset %d: expected %d bytes but read %d: %s", streamLen-actualFrameLen, frameLen, actualFrameLen, err)
		}
	}
	return streamLen, nil
}

func fastCommp(reader io.Reader) (writer.DataCIDSize, error) {
	cp := new(commp.Calc)
	streamBuf := bufio.NewReaderSize(
		io.TeeReader(reader, cp),
		BufSize,
	)

	var streamLen int64
	var carHdr *CarHeader

	if maybeHeaderLen, err := streamBuf.Peek(10); err == nil {
		if hdrLen, viLen := binary.Uvarint(maybeHeaderLen); viLen > 0 && hdrLen > 0 {
			actualViLen, err := io.CopyN(io.Discard, streamBuf, int64(viLen))
			streamLen += actualViLen
			if err == nil {
				hdrBuf := make([]byte, hdrLen)
				actualHdrLen, err := io.ReadFull(streamBuf, hdrBuf)
				streamLen += int64(actualHdrLen)
				if err == nil {
					carHdr = new(CarHeader)
					if cbor.DecodeInto(hdrBuf, carHdr) != nil {
						carHdr = nil
					} else if carHdr.Version == 1 {
						streamLen, err = process(streamBuf, streamLen)
						if err != nil {
							log.Fatal(err)
							return writer.DataCIDSize{}, err
						}
					}
				}
			}
		}
	}

	n, err := io.Copy(io.Discard, streamBuf)
	streamLen += n
	if err != nil && err != io.EOF {
		log.Fatalf("unexpected error at offset %d: %s", streamLen, err)
	}

	rawCommP, paddedSize, err := cp.Digest()
	if err != nil {
		log.Fatal(err)
	}

	commCid, err := commcid.DataCommitmentV1ToCID(rawCommP)
	if err != nil {
		log.Fatal(err)
	}

	return writer.DataCIDSize{
		PayloadSize: streamLen,
		PieceSize:   abi.PaddedPieceSize(paddedSize),
		PieceCID:    commCid,
	}, nil
}
