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
		diskSetup   func() *Disk
		blocksUsed  int
	}{

		{name: "data size exceeding  available disk space",
			expectedErr: ErrInsufficentMemoryError,
			diskSetup: func() *Disk {
				disk, _ := NewDisk(100, 10)
				return disk
			},
			data: make([]byte, 101),
		},
		{name: "data duccesfully written",
			expectedErr: nil,
			diskSetup: func() *Disk {
				disk, _ := NewDisk(100, 10)
				return disk
			},
			data:       make([]byte, 100),
			blocksUsed: 10,
		},
		{name: "data duccesfully written",
			expectedErr: nil,
			diskSetup: func() *Disk {
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
