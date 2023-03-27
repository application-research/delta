package core

import (
	"bytes"
	"context"
	"github.com/application-research/whypfs-core"
	"github.com/ipfs/go-cid"
	carv2 "github.com/ipld/go-car/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommpService_GenerateCommPFile(t *testing.T) {
	// create a mock DeltaNode instance
	deltaNode := &DeltaNode{}

	// create a mock blockstore instance
	params := whypfs.NewNodeParams{
		Ctx:       context.Background(),
		Datastore: whypfs.NewInMemoryDatastore(),
		Repo:      ".test",
	}
	whypfsPeer, err := whypfs.NewNode(params)
	blockstore := whypfsPeer.Blockstore

	// create a CommpService instance using the mock DeltaNode instance
	commpService := CommpService{
		DeltaNode: deltaNode,
	}

	// create a mock payload CID
	payloadCid, _ := cid.Decode("bafy2bzacedx6vywq6so7e6m43g6op7dhcqkntwk2lrzvt6yljvf6xlzshxh5w")

	// call the GenerateCommPFile method
	pieceCid, payloadSize, unPaddedPieceSize, err := commpService.GenerateCommPFile(context.Background(), payloadCid, blockstore)

	// assert that there is no error
	assert.NoError(t, err)

	// assert that the piece CID, payload size, and unpadded piece size are not empty
	assert.NotEmpty(t, pieceCid)
	assert.NotEmpty(t, payloadSize)
	assert.NotEmpty(t, unPaddedPieceSize)
}

func TestCommpService_GetSize(t *testing.T) {
	// create a mock payload reader
	reader := bytes.NewReader([]byte("hello world"))

	// create a CommpService instance
	commpService := CommpService{}

	// call the GetSize method
	size := commpService.GetSize(reader)

	// assert that the size is correct
	assert.Equal(t, 11, size)
}

func TestCommpService_GetCarSize(t *testing.T) {
	// create a mock CARv2 file reader
	reader := bytes.NewReader([]byte("file content"))

	// create a CARv2 reader from the mock CARv2 file reader
	carReader, _ := carv2.NewReader(reader)

	// create a CommpService instance
	commpService := CommpService{}

	// call the GetCarSize method
	size, err := commpService.GetCarSize(reader, carReader)

	// assert that there is no error
	assert.NoError(t, err)

	// assert that the size is correct
	assert.Equal(t, int64(12), size)
}

func TestCommpService_GenerateParallelCommp(t *testing.T) {
	// create a mock payload reader
	reader := bytes.NewReader([]byte("hello world"))

	// create a CommpService instance
	commpService := CommpService{}

	// call the GenerateParallelCommp method
	dataCIDSize, err := commpService.GenerateParallelCommp(reader)

	// assert that there is no error
	assert.NoError(t, err)

	// assert that the Data CID and Size are not empty
	assert.NotEmpty(t, dataCIDSize.PieceCID)
	assert.NotEmpty(t, dataCIDSize.PieceSize)
}
