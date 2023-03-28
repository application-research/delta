package core

import (
	"bufio"
	"bytes"
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
	Version uint64
	Roots   []cid.Cid
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

// Pragma is the first 11 bytes of a CARv2 file.
var Pragma = []byte{0x0a, 0xa1, 0x67, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x02}

const (
	// BufSize is the size of the buffer used to read the CAR file.
	BufSize = (4 << 20) / 128 * 127
	// PragmaSize is the size of the CARv2 pragma in bytes.
	PragmaSize = 11
	// HeaderSize is the fixed size of CARv2 header in number of bytes.
	HeaderSize = 40
	// CharacteristicsSize is the fixed size of Characteristics bitfield within CARv2 header in number of bytes.
	CharacteristicsSize = 16
)

func init() {
	cbor.RegisterCborType(CarHeader{})
}

// checkCarV2 checks if the given file is a CARv2 file and returns the header if it is.
func checkCarV2(reader io.ReadSeekCloser) (bool, *CarV2Header) {
	defer reader.Seek(0, 0)
	// Read the first 11 bytes of the file into a byte slice
	pragmaHeader := make([]byte, PragmaSize)
	_, err := reader.Read(pragmaHeader)
	if err != nil {
		fmt.Println("Error reading file:", err)
		panic(err)
	}

	carV2Header := &CarV2Header{}

	// Compare the first 11 bytes of the file to the expected header
	if bytes.Equal(pragmaHeader, Pragma) {
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

// readCarHeader reads the CARv1 header from a CARv2 file
func readCarHeader(streamBuf *bufio.Reader, streamLen int64) (carHeader *CarHeader, strLen int64, err error) {
	// Read the first 10 bytes of the file into a byte slice
	headerLengthBytes, err := streamBuf.Peek(10)
	if err != nil {
		return nil, 0, err
	}
	// Read the header length
	headerLength, headerBytesRead := binary.Uvarint(headerLengthBytes)
	if headerLength == 0 || headerBytesRead < 0 {
		return nil, 0, fmt.Errorf("invalid header length")
	}
	// Read the header
	realHeaderLength, err := io.CopyN(io.Discard, streamBuf, int64(headerBytesRead))
	if err != nil {
		return nil, 0, err
	}
	streamLen += realHeaderLength
	headerBuffer := make([]byte, headerLength)
	actualHdrLen, err := io.ReadFull(streamBuf, headerBuffer)
	if err != nil {
		return nil, 0, err
	}
	streamLen += int64(actualHdrLen)

	// Decode the header
	carHeader = new(CarHeader)
	err = cbor.DecodeInto(headerBuffer, carHeader)
	if err != nil {
		return nil, 0, err
	}
	return carHeader, streamLen, nil
}

// process the carfile blocks
func process(streamBuf *bufio.Reader, streamLen int64) (strLen int64, err error) {
	for {
		nextBlockBuffer, err := streamBuf.Peek(10)
		if err == io.EOF {
			break
		}
		if err != nil && err != bufio.ErrBufferFull {
			return streamLen, err
		}
		if len(nextBlockBuffer) == 0 {
			return streamLen, err
		}

		blockLength, viLen := binary.Uvarint(nextBlockBuffer)
		if viLen <= 0 {
			return streamLen, err
		}
		if blockLength > 2<<20 {
			// anything over ~2MiB got to be a mistake
			return streamLen, fmt.Errorf("large block length too large: %d bytes at offset %d", blockLength, streamLen)
		}
		actualBlockLength, err := io.CopyN(io.Discard, streamBuf, int64(viLen)+int64(blockLength))
		streamLen += actualBlockLength
		if err != nil {
			if err != io.EOF {
				log.Fatalf("unexpected error at offset %d: %s", streamLen-actualBlockLength, err)
			}
			return streamLen, fmt.Errorf("aborting car stream parse: truncated block at offset %d: expected %d bytes but read %d: %s", streamLen-actualBlockLength, blockLength, actualBlockLength, err)
		}
	}
	return streamLen, nil
}

// fastCommp calculates the commp of a CARv1 or CARv2 file
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

	// The length of the stream
	var streamLen int64

	// Read the header
	carHeader, streamLen, err := readCarHeader(streamBuf, streamLen)
	if err != nil {
		return writer.DataCIDSize{}, err
	}

	if carHeader.Version == 1 || carHeader.Version == 2 {
		streamLen, err = process(streamBuf, streamLen)
		if err != nil {
			log.Fatal(err)
			return writer.DataCIDSize{}, err
		}
	} else {
		return writer.DataCIDSize{}, fmt.Errorf("invalid car version: %d", carHeader.Version)
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
