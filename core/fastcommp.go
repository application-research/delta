package core

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
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

// CarV2Header is the fixed-size header of a CARv2 file.
type CarV2Header struct {
	// Characteristics is a bitfield of characteristics that apply to the CARv2.
	Characteristics [16]byte
	// DataOffset is the byte offset from the beginning of the CARv2 to the beginning of the CARv1 data payload.
	DataOffset uint64
	// DataSize is the size of the CARv1 data payload in bytes.
	DataSize uint64
	// IndexOffset is the byte offset from the beginning of the CARv2 to the beginning of the CARv1 index payload.
	IndexOffset uint64
}

func init() {
	cbor.RegisterCborType(CarHeader{})
}

const (
	BufSize = (4 << 20) / 128 * 127
	// PragmaSize is the size of the CARv2 pragma in bytes.
	PragmaSize = 11
	// HeaderSize is the fixed size of CARv2 header in number of bytes.
	HeaderSize = 40
	// CharacteristicsSize is the fixed size of Characteristics bitfield within CARv2 header in number of bytes.
	CharacteristicsSize = 16
)

// checkCarV2 checks if the given file is a CARv2 file and returns the header if it is.
func checkCarV2(reader io.ReadSeekCloser) (bool, *CarV2Header) {
	defer reader.Seek(0, 0)
	// Read the first 11 bytes of the file into a byte slice
	pragmaHeader := make([]byte, 11)
	_, err := reader.Read(pragmaHeader)
	if err != nil {
		fmt.Println("Error reading file:", err)
		panic(err)
	}

	carV2Header := &CarV2Header{}

	// Convert the expected header to a byte slice
	expectedHeader, err := hex.DecodeString("0aa16776657273696f6e02")
	if err != nil {
		fmt.Println("Error decoding hex string:", err)
		panic(err)
	}

	// Compare the first 11 bytes of the file to the expected header
	if bytes.Equal(pragmaHeader, expectedHeader) {
		// Read the next 40 bytes of the file into a byte slice
		header := make([]byte, 40)
		_, err = reader.Read(header)
		if err != nil {
			fmt.Println("Error reading file:", err)
			panic(err)
		}

		// Read the characteristics
		copy(carV2Header.Characteristics[:], header[:16])

		// Read the data offset
		carV2Header.DataOffset = binary.LittleEndian.Uint64(header[16:24])

		// Read the data size
		carV2Header.DataSize = binary.LittleEndian.Uint64(header[24:32])

		// Read the index offset
		carV2Header.IndexOffset = binary.LittleEndian.Uint64(header[32:40])
		return true, carV2Header
	} else {
		return false, nil
	}
}

// extractCarV1 extracts the CARv1 data from a CARv2 file
func extractCarV1(file io.ReadSeekCloser, offset, length int) (*bytes.Reader, error) {
	// Slice out the portion of the file
	_, err := file.Seek(int64(offset), 0)
	if err != nil {
		return nil, err
	}
	slice := make([]byte, length)
	_, err = file.Read(slice)
	if err != nil {
		return nil, err
	}

	// Create a new io.Reader from the slice
	sliceReader := bytes.NewReader(slice)

	return sliceReader, nil
}

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

func fastCommp(reader io.ReadSeekCloser) (writer.DataCIDSize, error) {
	// Check if the file is a CARv2 file
	isVarV2, headerInfo := checkCarV2(reader)
	var streamBuf *bufio.Reader
	cp := new(commp.Calc)
	if isVarV2 {
		// Extract the CARv1 data from the CARv2 file
		sliced, err := extractCarV1(reader, int(headerInfo.DataOffset), int(headerInfo.DataSize))
		if err != nil {
			panic(err)
		}
		streamBuf = bufio.NewReaderSize(
			io.TeeReader(sliced, cp),
			BufSize,
		)
	} else {
		// Read the file as a CARv1 file
		streamBuf = bufio.NewReaderSize(
			io.TeeReader(reader, cp),
			BufSize,
		)
	}

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
