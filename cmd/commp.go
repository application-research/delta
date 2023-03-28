package cmd

import (
	"bufio"
	"bytes"
	"context"
	"delta/api"
	c "delta/config"
	"delta/core"
	"delta/utils"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/application-research/filclient"
	"github.com/application-research/whypfs-core"
	"github.com/filecoin-project/go-commp-utils/writer"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	_ "github.com/ipfs/go-merkledag"
	"github.com/urfave/cli/v2"
)

type PieceCommitmentResult struct {
	FileName        string                     `json:"file_name,omitempty"`
	Size            int64                      `json:"size,omitempty"`
	Cid             string                     `json:"cid,omitempty"`
	PieceCommitment api.PieceCommitmentRequest `json:"piece_commitment,omitempty"`
}

type FilclientCommp struct {
	PayloadCid        cid.Cid
	PaddedPieceSize   uint64
	UnpaddedPieceSize abi.UnpaddedPieceSize
}

// CommpCmd A CLI command that generates a piece commitment for a given file.
// `CommpCmd()` returns a slice of `*cli.Command`s
func CommpCmd(cfg *c.DeltaConfig) []*cli.Command {
	// add a command to run API node
	var commpCommands []*cli.Command
	var commpService = new(core.CommpService)
	commpFileCmd := &cli.Command{
		Name:        "commp",
		Usage:       "Run the piece commitment computation and generation for a given file.",
		Description: " `commp` is a CLI command that generates a piece commitment for a given file or a directory of files.",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file",
				Usage: "specify the file",
			},
			&cli.StringFlag{
				Name:  "dir",
				Usage: "specify the directory to read and create a piece commitment for all the files in the directory",
			},
			&cli.StringFlag{
				Name:  "mode",
				Usage: "specify the mode of the piece commitment generation (default: fast. options: filboost, stream, fast)",
				Value: "fast",
			},
			&cli.BoolFlag{
				Name:  "include-payload-cid",
				Usage: "specify whether to include the payload cid in the piece commitment or not (default: false)",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {

			file := c.String("file")
			dir := c.String("dir")
			includePayloadCID := c.Bool("include-payload-cid")

			// TODO: make it work for now and clean up after.
			if file != "" {
				var commpRequest PieceCommitmentResult
				openFile, err := os.Open(file)
				reader := bufio.NewReader(openFile)

				if err != nil {
					fmt.Println(err)
					return err
				}
				var whypfsNode *whypfs.Node
				if includePayloadCID {
					params := whypfs.NewNodeParams{
						Ctx:       context.Background(),
						Datastore: whypfs.NewInMemoryDatastore(),
					}
					whypfsNode, err = whypfs.NewNode(params)
					if err != nil {
						fmt.Println(err)
						return err
					}
					fileNode, err := whypfsNode.AddPinFile(context.Background(), reader, nil)
					if err != nil {
						fmt.Println(err)
						return err
					}
					commpRequest.Cid = fileNode.Cid().String()
					commpRequest.FileName = openFile.Name()
				}

				var carV2PieceInfo *abi.PieceInfo
				var dataCidPieceInfo writer.DataCIDSize
				var filclientCommp FilclientCommp

				if c.String("mode") == "stream" {
					fileToStream, err := os.Open(file)
					if err != nil {
						fmt.Println(err)
						return err
					}
					fileToStreamReader := bufio.NewReader(fileToStream)
					carV2PieceInfo, err = commpService.GenerateCommPCarV2(fileToStreamReader)
					if err != nil {
						fmt.Println(err)
						return err
					}
				} else if c.String("mode") == "filboost" {
					cidToCompute, err := cid.Decode(commpRequest.Cid)
					if err != nil {
						fmt.Println(err)
						return err
					}
					payloadCid, paddedPieceSize, unpaddedPieceSize, err := filclient.GeneratePieceCommitment(context.Background(), cidToCompute, whypfsNode.Blockstore)
					if err != nil {
						fmt.Println(err)
						return err
					}
					filclientCommp = FilclientCommp{
						PayloadCid:        payloadCid,
						PaddedPieceSize:   paddedPieceSize,
						UnpaddedPieceSize: unpaddedPieceSize,
					}
				} else {
					fileToStream, err := os.Open(file)
					if err != nil {
						fmt.Println(err)
						return err
					}
					dataCidPieceInfo, err = commpService.GenerateCommp(fileToStream)
					if err != nil {
						fmt.Println(err)
						return err
					}
				}

				if err != nil {
					fmt.Println(err)
					return err
				}

				// return json to console.
				commpRequest.PieceCommitment.Piece = func() string {
					if c.String("mode") == "stream" {
						return carV2PieceInfo.PieceCID.String()
					} else if c.String("mode") == "filboost" {
						return filclientCommp.PayloadCid.String()
					}
					return dataCidPieceInfo.PieceCID.String()
				}()
				commpRequest.PieceCommitment.PaddedPieceSize = func() uint64 {
					if c.String("mode") == "stream" {
						return uint64(carV2PieceInfo.Size)
					} else if c.String("mode") == "filboost" {
						return uint64(filclientCommp.PaddedPieceSize)
					}
					return uint64(dataCidPieceInfo.PieceSize)
				}()

				commpRequest.PieceCommitment.UnPaddedPieceSize = func() uint64 {
					if c.String("mode") == "stream" {
						return uint64(carV2PieceInfo.Size.Unpadded())
					} else if c.String("mode") == "filboost" {
						return uint64(filclientCommp.UnpaddedPieceSize)
					}
					return uint64(dataCidPieceInfo.PieceSize.Unpadded())
				}()

				if err != nil {
					fmt.Println(err)
					return err
				}
				commpRequest.Size = func() int64 {
					if c.String("mode") == "stream" {
						commpRequest.Size = int64(reader.Size())
					} else if c.String("mode") == "filboost" {
						commpRequest.Size = int64(filclientCommp.UnpaddedPieceSize)
					}
					return dataCidPieceInfo.PayloadSize
				}()

				var buffer bytes.Buffer
				err = utils.PrettyEncode(commpRequest, &buffer)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(buffer.String())
			}
			if dir != "" {

				var dealRequests []PieceCommitmentResult

				var whypfsNode *whypfs.Node
				var err error
				if includePayloadCID {
					params := whypfs.NewNodeParams{
						Ctx:       context.Background(),
						Datastore: whypfs.NewInMemoryDatastore(),
					}
					whypfsNode, err = whypfs.NewNode(params)
					if err != nil {
						fmt.Println(err)
						return err
					}
				}

				err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}

					// Check if the entry is a file
					if !d.Type().IsRegular() {
						return nil
					}

					// Open the file
					fileOpen, err := os.Open(path)

					if err != nil {
						return err
					}

					var carV2PieceInfo *abi.PieceInfo
					var dataCidPieceInfo writer.DataCIDSize
					var filclientCommp FilclientCommp
					var requestInDir PieceCommitmentResult

					if c.String("mode") == "stream" {
						fileToStreamReader := bufio.NewReader(fileOpen)
						carV2PieceInfo, err = commpService.GenerateCommPCarV2(fileToStreamReader)

						if err != nil {
							fmt.Println(err)
							return err
						}
						if includePayloadCID {

							fileNode, err := whypfsNode.AddPinFile(context.Background(), fileToStreamReader, nil)
							if err != nil {
								fmt.Println(err)
								return err
							}
							requestInDir.Cid = fileNode.Cid().String()
							requestInDir.FileName = fileOpen.Name()
						}
					} else {
						dataCidPieceInfo, err = commpService.GenerateCommp(fileOpen)
						if err != nil {
							fmt.Println(err)
							return err
						}
						if includePayloadCID {
							fileNode, err := whypfsNode.AddPinFile(context.Background(), fileOpen, nil)
							if err != nil {
								fmt.Println(err)
								return err
							}
							requestInDir.Cid = fileNode.Cid().String()
							requestInDir.FileName = fileOpen.Name()
						}
					}

					if err != nil {
						fmt.Println(err)
						return err
					}

					// return json to console.
					requestInDir.PieceCommitment.Piece = func() string {
						if c.String("mode") == "stream" {
							return carV2PieceInfo.PieceCID.String()
						} else if c.String("mode") == "filboost" {
							return filclientCommp.PayloadCid.String()
						}
						return dataCidPieceInfo.PieceCID.String()
					}()
					requestInDir.PieceCommitment.PaddedPieceSize = func() uint64 {
						if c.String("mode") == "stream" {
							return uint64(carV2PieceInfo.Size)
						} else if c.String("mode") == "filboost" {
							return uint64(filclientCommp.PaddedPieceSize)
						}
						return uint64(dataCidPieceInfo.PieceSize)
					}()

					requestInDir.PieceCommitment.UnPaddedPieceSize = func() uint64 {
						if c.String("mode") == "stream" {
							return uint64(carV2PieceInfo.Size.Unpadded())
						} else if c.String("mode") == "filboost" {
							return uint64(filclientCommp.UnpaddedPieceSize)
						}
						return uint64(dataCidPieceInfo.PieceSize.Unpadded())
					}()

					if err != nil {
						fmt.Println(err)
						return err
					}
					requestInDir.Size = func() int64 {
						if c.String("mode") == "stream" {
							requestInDir.Size = int64(carV2PieceInfo.Size)
						}
						return dataCidPieceInfo.PayloadSize
					}()

					dealRequests = append(dealRequests, requestInDir)

					return nil
				})
				var buffer bytes.Buffer
				err = utils.PrettyEncode(dealRequests, &buffer)
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println(buffer.String())
				if err != nil {
					fmt.Println(err)
				}
			}
			return nil
		},
	}
	commpCommands = append(commpCommands, commpFileCmd)

	return commpCommands

}

func getPieceInfoRequest(file string, node whypfs.Node) (PieceCommitmentResult, error) {
	var request PieceCommitmentResult
	var commpService = new(core.CommpService)

	fileOpen, err := os.Open(file)
	if err != nil {
		return request, err
	}

	fileToStreamReader := bufio.NewReader(fileOpen)
	carV2PieceInfo, err := commpService.GenerateCommPCarV2(fileToStreamReader)
	if err != nil {
		fmt.Println(err)
		return request, err
	}

	// return json to console.
	request.PieceCommitment.Piece = func() string {
		return carV2PieceInfo.PieceCID.String()
	}()
	request.PieceCommitment.PaddedPieceSize = func() uint64 {
		return uint64(carV2PieceInfo.Size)
	}()

	request.PieceCommitment.UnPaddedPieceSize = func() uint64 {
		return uint64(carV2PieceInfo.Size.Unpadded())
	}()

	if err != nil {
		fmt.Println(err)
		return request, err
	}
	request.Size = func() int64 {
		return int64(carV2PieceInfo.Size)
	}()

	return request, nil
}
