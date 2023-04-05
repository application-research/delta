package core

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/filecoin-project/go-commp-utils/writer"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	carv2 "github.com/ipld/go-car/v2"
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
func extractCarV1_777(file io.ReadSeekCloser, offset, length int) (*bytes.Reader, error) {
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

func extractCarV1(file io.ReadSeekCloser, offset, length int) (io.ReadSeekCloser, error) {
	// Slice out the portion of the file

	_, err := file.Seek(int64(offset), 0)
	if err != nil {
		return nil, err
	}

	/*
		slice := make([]byte, length)
		_, err = file.Read(slice)
		if err != nil {
			return nil, err
		}
	*/

	// Create a new io.Reader from the slice
	//sliceReader := bytes.NewReader(slice)

	return file, nil
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
func process(streamBuf *bufio.Reader, streamLen int64, maxSize int64) (strLen int64, err error) {
	for {
		nextBlockBuffer, err := streamBuf.Peek(10)
		if err == io.EOF {
			break
		}

		if maxSize != 0 && streamLen+10 >= maxSize {
			fmt.Println("max size reached, diff: ", maxSize, streamLen, maxSize-streamLen)
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

// GenerateCommP calculates commp locally
func GenerateCommP(filepath string) (*abi.PieceInfo, error) {
	rd, err := carv2.OpenReader(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to get CARv2 reader: %w", err)
	}

	defer func() {
		if err := rd.Close(); err != nil {
			fmt.Printf("failed to close CARv2 reader for %s: %w", filepath, err)
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

	// get the size of the CAR file
	size, err := getCarSize(filepath, rd)
	if err != nil {
		return nil, err
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

func getCarSize(filepath string, rd *carv2.Reader) (int64, error) {
	var size int64
	switch rd.Version {
	case 2:
		size = int64(rd.Header.DataSize)
	case 1:
		st, err := os.Stat(filepath)
		if err != nil {
			return 0, fmt.Errorf("failed to get CARv1 file size: %w", err)
		}
		size = st.Size()
	}
	return size, nil
}

// fastCommp calculates the commp of a CARv1 or CARv2 file
func fastCommp(filename string) (writer.DataCIDSize, error) {
	result, err := GenerateCommP(filename)

	if err != nil {
		fmt.Println("Error generating CommP:", err)
		return writer.DataCIDSize{}, err
	}

	fmt.Println("CommP:", result.PieceCID)
	return writer.DataCIDSize{
		PayloadSize: 0,
		PieceSize:   abi.PaddedPieceSize(result.Size),
		PieceCID:    result.PieceCID,
	}, nil
}
