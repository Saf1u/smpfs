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
