package cmd

import (
	"bufio"
	"context"
	c "delta/config"
	"delta/utils"
	"encoding/json"
	"fmt"
	"github.com/filecoin-project/go-fil-commp-hashhash"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type CommpResult struct {
	commp     string
	pieceSize uint64
}

type Result struct {
	Ipld      *utils.FsNode
	DataCid   string
	PieceCid  string
	PieceSize uint64
	CidMap    map[string]utils.CidMapValue
}

type Input []utils.Finfo

type CarHeader struct {
	Roots   []cid.Cid
	Version uint64
}

func init() {
	cbor.RegisterCborType(CarHeader{})
}

const BufSize = (4 << 20) / 128 * 127

func CarCmd(cfg *c.DeltaConfig) []*cli.Command {
	ctx := context.TODO()
	// add a command to run API node
	var carCommands []*cli.Command
	carCmd := &cli.Command{
		Name:  "car",
		Usage: "Create a car file from a list of files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "source",
				Usage: "Source of the input (file path, dir path or json string)",
			},
			&cli.StringFlag{
				Name:  "split-size",
				Usage: "Split size in bytes",
				Value: "",
			},
			&cli.StringFlag{
				Name:  "output-dir",
				Usage: "Output directory for the car file. If not specified, the car file will be created in the same directory as the input file",
				Value: ".",
			},
		},
		Action: func(c *cli.Context) error {
			sourceInput := c.String("source")
			splitSizeInput := c.String("split-size")
			outDir := c.String("out-dir")
			fmt.Println(splitSizeInput)

			if splitSizeInput != "" {
				splitSizeA, err := strconv.Atoi(splitSizeInput)
				if err != nil {
					return err
				}
				splitSize := int64(splitSizeA)
				var input Input
				var outputs []Result
				err = filepath.Walk(sourceInput, func(sourcePath string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if info.IsDir() {
						return nil
					}

					// Calculate the number of chunks based on the split size.
					chunks := (info.Size() + splitSize - 1) / splitSize

					// Open the source file for reading.
					sourceFile, err := os.Open(sourcePath)
					if err != nil {
						return err
					}
					//defer sourceFile.Close()

					// Create a new directory for the file chunks.
					fileName := info.Name()
					chunkDir := filepath.Join(outDir, fileName)
					err = os.MkdirAll(chunkDir, 0755)
					if err != nil {
						return err
					}

					// Split the file into chunks.
					for i := int64(0); i < chunks; i++ {
						// Open a new file for writing the chunk.
						chunkFileName := fmt.Sprintf("%s_%04d", fileName, i)
						chunkFilePath := filepath.Join(chunkDir, chunkFileName)
						chunkFile, err := os.Create(chunkFilePath)
						if err != nil {
							return err
						}
						//defer chunkFile.Close()

						// Copy the chunk data from the source file to the chunk file.
						start := i * splitSize
						end := (i + 1) * splitSize
						if end > info.Size() {
							end = info.Size()
						}
						_, err = sourceFile.Seek(start, 0)
						if err != nil {
							return err
						}
						written, err := io.CopyN(chunkFile, sourceFile, end-start)
						if err != nil && err != io.EOF {
							return err
						}
						input = append(input, utils.Finfo{
							Path:  chunkFilePath,
							Size:  written,
							Start: 0,
							End:   written,
						})
						outFilename := uuid.New().String() + ".car"
						outPath := path.Join(outDir, outFilename)
						carF, err := os.Create(outPath)
						if err != nil {
							return err
						}
						cp := new(commp.Calc)
						writer := bufio.NewWriterSize(io.MultiWriter(carF, cp), BufSize)
						ipld, cid, cidMap, err := utils.GenerateCar(ctx, input, "", "", writer)
						if err != nil {
							return err
						}
						err = writer.Flush()
						if err != nil {
							return err
						}

						outputs = append(outputs, Result{
							Ipld:    ipld,
							DataCid: cid,
							CidMap:  cidMap,
						})
					}

					outputj, err := json.Marshal(outputs)
					if err != nil {
						return err
					}
					fmt.Println(string(outputj))
					return nil
				})
				if err != nil {
					panic(err)
				}
			} else {
				var input Input
				stat, err := os.Stat(sourceInput)
				if err != nil {
					return err
				}
				if stat.IsDir() {
					err := filepath.Walk(sourceInput, func(sourcePath string, info os.FileInfo, err error) error {

						//if splitSize == 0 {
						if err != nil {
							return err
						}
						if info.IsDir() {
							return nil
						}
						input = append(input, utils.Finfo{
							Path:  sourcePath,
							Size:  info.Size(),
							Start: 0,
							End:   info.Size(),
						})
						outFilename := uuid.New().String() + ".car"
						outPath := path.Join(outDir, outFilename)
						carF, err := os.Create(outPath)
						if err != nil {
							return err
						}
						cp := new(commp.Calc)
						writer := bufio.NewWriterSize(io.MultiWriter(carF, cp), BufSize)
						ipld, cid, cidMap, err := utils.GenerateCar(ctx, input, "", "", writer)
						if err != nil {
							return err
						}
						err = writer.Flush()
						if err != nil {
							return err
						}
						output, err := json.Marshal(Result{
							Ipld:    ipld,
							DataCid: cid,
							CidMap:  cidMap,
						})
						if err != nil {
							return err
						}
						fmt.Println(string(output))
						return nil
					})
					if err != nil {
						return err
					}
				} else {
					input = append(input, utils.Finfo{
						Path:  sourceInput,
						Size:  stat.Size(),
						Start: 0,
						End:   stat.Size(),
					})
					outFilename := uuid.New().String() + ".car"
					outPath := path.Join(outDir, outFilename)
					carF, err := os.Create(outPath)
					if err != nil {
						return err
					}
					cp := new(commp.Calc)
					writer := bufio.NewWriterSize(io.MultiWriter(carF, cp), BufSize)
					ipld, cid, cidMap, err := utils.GenerateCar(ctx, input, "", "", writer)
					if err != nil {
						return err
					}
					err = writer.Flush()
					if err != nil {
						return err
					}

					output, err := json.Marshal(Result{
						Ipld:    ipld,
						DataCid: cid,
						CidMap:  cidMap,
					})
					if err != nil {
						return err
					}
					fmt.Println(string(output))
				}
				return nil
			}
			return nil
		},
	}
	carCommands = append(carCommands, carCmd)

	return carCommands
}
