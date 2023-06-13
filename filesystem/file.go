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
	setCreationTs(time.Time)
	updateAccessTs(time.Time)
	getSize() int
	getManifest() *disk.BlockRecord
	setManifest(*disk.BlockRecord)
}

func NewFile(name string) File {
	return &file{fileName: name, info: disk.NewBlockRecord(), createdAt: time.Now(), lastModified: time.Now()}
}

func (fl *file) isFile() bool {
	return true
}
func (fl *file) name() string {
	return fl.fileName
}
func (fl *file) updateAccessTs(time time.Time) {
	fl.lastModified = time
}

func (fl *file) setCreationTs(time time.Time) {
	fl.createdAt = time
}

func (fl *file) getManifest() *disk.BlockRecord {
	return fl.info
}

func (fl *file) setManifest(manifest *disk.BlockRecord) {
	fl.info = manifest
}

func (fl *file) getSize() int {
	return 0
}
