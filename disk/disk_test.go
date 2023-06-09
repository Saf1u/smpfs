package disk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDisk(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr error
		diskSize    int
		blockSize   int
	}{

		{name: "blocksize exceeding disk space", expectedErr: ErrBlockSizeExceedsDriveSize, diskSize: 10, blockSize: 20},
		{name: "successfully initalized disk", expectedErr: nil, diskSize: 10, blockSize: 10},
		{name: "successfully initalized disk", expectedErr: nil, diskSize: 20, blockSize: 10},
	}
	for _, testcase := range tests {
		_, err := NewDisk(testcase.diskSize, testcase.blockSize)
		assert.Equal(t, testcase.expectedErr, err)
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr error
		data        []byte
		diskSetup   func() Disk
		blocksUsed  int
	}{

		{name: "data size exceeding  available disk space",
			expectedErr: ErrInsufficentMemoryError,
			diskSetup: func() Disk {
				disk, _ := NewDisk(100, 10)
				return disk
			},
			data: make([]byte, 101),
		},
		{name: "data duccesfully written",
			expectedErr: nil,
			diskSetup: func() Disk {
				disk, _ := NewDisk(100, 10)
				return disk
			},
			data:       make([]byte, 100),
			blocksUsed: 10,
		},
		{name: "data duccesfully written",
			expectedErr: nil,
			diskSetup: func() Disk {
				disk, _ := NewDisk(100, 100)
				return disk
			},
			data:       make([]byte, 100),
			blocksUsed: 1,
		},
	}
	for _, testcase := range tests {
		disk := testcase.diskSetup()
		rec, err := disk.Write(testcase.data)
		assert.Equal(t, testcase.expectedErr, err)
		if err == nil {
			assert.Equal(t, len(rec.blocks), testcase.blocksUsed)
		}
	}
}

func TestRead(t *testing.T) {
	tests := []struct {
		name         string
		expectedErr  error
		diskSetup    func() (Disk, *BlockRecord)
		expectedData []byte
	}{

		{name: "data duccesfully written",
			expectedErr: nil,
			diskSetup: func() (Disk, *BlockRecord) {
				disk, _ := NewDisk(1000, 5)
				disk.Write([]byte("hello I'm a string scattered in the memory region but is not read"))
				manifest, _ := disk.Write([]byte("hello I'm a string scattered in the memory region"))
				return disk, manifest
			},
			expectedData: []byte("hello I'm a string scattered in the memory region"),
		},
	}
	for _, testcase := range tests {
		disk, manifest := testcase.diskSetup()
		rec, err := disk.Read(manifest)
		assert.Equal(t, testcase.expectedErr, err)
		if err == nil {
			assert.Equal(t, rec, testcase.expectedData)
		}
	}
}

func TestGetUnfilledBlock(t *testing.T) {
	record := NewBlockRecord()
	record.addBlock(block{startIndex: 0, endIndex: 100, used: 95, size: 100})
	extraBlock := record.getUnfilledBlock()
	assert.Equal(t, block{startIndex: 0, endIndex: 100, used: 95, size: 100}, *extraBlock)

	record = NewBlockRecord()
	record.addBlock(block{startIndex: 0, endIndex: 100, used: 100, size: 100})
	extraBlock = record.getUnfilledBlock()
	var block *block
	//nil block
	assert.Equal(t, extraBlock, block)
}
