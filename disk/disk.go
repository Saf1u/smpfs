package disk

import (
	"archive/zip"
	"bytes"
	"errors"
	"log"
	"math"
	"os"

	"time"

	"github.com/Saf1u/smpfs/pool"
)

// Package disk defines a in memory byte block allocator
type block struct {
	startIndex int
	endIndex   int
	used       int
	size       int
}

func (b *block) Reset() {
	b.used = b.size
}
func (b *block) SetUsed(size int) {
	b.used = size
}

type disk struct {
	buffer    []byte
	blockPool *pool.Pool
	blockSize int
}

type Disk interface {
	Write(fileBytes []byte) (*BlockRecord, error)
	Read(blockManifest *BlockRecord) ([]byte, error)
	Delete(blockManifest *BlockRecord)
	SaveDisk()
}

type BlockRecord struct {
	blocks []block
}

func NewBlockRecord() *BlockRecord {
	return &BlockRecord{make([]block, 0)}
}

func (blockRecord *BlockRecord) addBlock(b block) {
	blockRecord.blocks = append(blockRecord.blocks, b)
}

func (blockRecord *BlockRecord) getUnfilledBlock() *block {
	if len(blockRecord.blocks) == 0 {
		return nil
	}
	lastBlock := blockRecord.blocks[len(blockRecord.blocks)-1]
	if lastBlock.size == lastBlock.used {
		return nil
	}
	return &lastBlock
}

var (
	ErrBlockSizeExceedsDriveSize = errors.New("block size is greater than available disk")
	ErrInsufficentMemoryError    = errors.New("not enough memory present to store file")
)

func NewDisk(size int, blockSize int) (Disk, error) {
	if blockSize > size {
		return nil, ErrBlockSizeExceedsDriveSize
	}
	pool := pool.NewPool()
	numberOfBlocks := size / blockSize
	//floor div

	startIndex := 0
	endIndex := blockSize - 1
	for blockNum := 0; blockNum < numberOfBlocks; blockNum++ {
		pool.AddToPool(block{startIndex, endIndex, 0, blockSize})
		startIndex = endIndex + 1
		endIndex = (endIndex + blockSize)
	}
	return &disk{
		buffer:    make([]byte, size),
		blockPool: pool,
		blockSize: blockSize,
	}, nil
}

func (disk *disk) Write(fileBytes []byte) (*BlockRecord, error) {
	availableBlocks := disk.blockPool.AvaialbleResourceUnits()
	blocksNeeded := int(math.Ceil(float64(len(fileBytes)) / float64(disk.blockSize)))
	if blocksNeeded > availableBlocks {
		return nil, ErrInsufficentMemoryError
	}
	blockManifest := &BlockRecord{}

	//Wrap buffer for easy reads
	fileBuffer := bytes.NewBuffer(fileBytes)

	var dataBlock block

	for blocksNeeded != 0 {
		dataBlock = disk.blockPool.GetResource().(block)
		readSize := disk.writeDataToBlock(&dataBlock, fileBuffer, dataBlock.startIndex, dataBlock.endIndex+1)
		dataBlock.SetUsed(readSize)
		blockManifest.addBlock(dataBlock)
		blocksNeeded--
	}

	return blockManifest, nil
}
func (disk *disk) writeDataToBlock(dataBlock *block, fileBuffer *bytes.Buffer, startIndex, endIndex int) int {

	memoryBlock := disk.buffer[startIndex:endIndex]
	readSize, err := fileBuffer.Read(memoryBlock)
	if err != nil {
		panic(err)
	}
	return readSize

}

func (disk *disk) Read(blockManifest *BlockRecord) ([]byte, error) {
	outBuffer := make([]byte, 0)
	bufferWrapper := bytes.NewBuffer(outBuffer)
	for _, blocks := range blockManifest.blocks {
		_, err := bufferWrapper.Write(disk.buffer[blocks.startIndex : blocks.startIndex+blocks.used])
		if err != nil {
			panic(err)
		}

	}
	return bufferWrapper.Bytes(), nil
}

func (disk *disk) Append(blockManifest *BlockRecord, fileBytes []byte) error {
	//Wrap buffer for easy reads
	fileBuffer := bytes.NewBuffer(fileBytes)
	extraBlock := blockManifest.getUnfilledBlock()
	dataToReadSize := len(fileBytes)
	if extraBlock != nil {
		numBytesRead := disk.writeDataToBlock(extraBlock, fileBuffer, extraBlock.used+extraBlock.startIndex, extraBlock.endIndex+1)
		if numBytesRead == dataToReadSize {
			return nil
		}
		dataToReadSize = dataToReadSize - numBytesRead
	}
	appendedRecords, err := disk.Write(fileBytes[dataToReadSize:])
	if err != nil {
		return err
	}
	mergeBlockRecords(blockManifest, appendedRecords)
	return nil

}

func mergeBlockRecords(blockDest, blockSrc *BlockRecord) {
	for i := 0; i < len(blockSrc.blocks); i++ {
		blockDest.addBlock(blockDest.blocks[i])
	}
}

func (disk *disk) Delete(blockManifest *BlockRecord) {
	//no zeroing needed
	for _, block := range blockManifest.blocks {
		block.SetUsed(0)
		disk.blockPool.AddToPool(block)
	}
}

func (disk *disk) SaveDisk() {
	name := "DISKSNAPSHOT-" + time.Now().Format("Jan _2 15:04:05.000000000")
	zipFile, err := os.Create(name + ".zip")
	if err != nil {
		panic(err)
	}
	zipper := zip.NewWriter(zipFile)
	saveFile, err := zipper.Create(name)
	saveFile.Write(disk.buffer)
	if err != nil {
		panic(err)
	}
	err = zipper.Close()
	if err != nil {
		log.Fatal(err)
	}
}
