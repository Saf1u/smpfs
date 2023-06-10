package filesystem

import (
	"time"

	"github.com/Saf1u/smpfs/disk"
)

type file struct {
	fileName     string
	info         *disk.BlockRecord
	createdAt    time.Time
	lastModified time.Time
}

type File interface {
	Close() error
}

func (fl *file) isFile() bool {
	return true
}
func (fl *file) name() string {
	return fl.fileName
}

func (fl *file) Close() error {
	return nil
}
