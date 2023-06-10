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
	getManifest() *disk.BlockRecord
	setManifest(*disk.BlockRecord)
}

func (fl *file) isFile() bool {
	return true
}
func (fl *file) name() string {
	return fl.fileName
}
func (fl *file) updateAccessTs(time time.Time) {

}

func (fl *file) setCreationTs(time time.Time) {

}

func (fl *file) getManifest() *disk.BlockRecord {
	return nil
}

func (fl *file) setManifest(manifest *disk.BlockRecord) {

}
