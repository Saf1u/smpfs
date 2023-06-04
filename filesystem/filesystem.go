package filesystem

import (
	"strings"

	"github.com/Saf1u/smpfs/disk"
)

type fileSystem struct {
	root directory
	disk *disk.Disk
}

type item interface {
	isFile() bool
	name() string
	read() []byte
	write([]byte)
}

type directory struct {
	contents []item
}

type file struct {
	info disk.BlockRecord
}

func (f *fileSystem) CreateDir(path string, data []byte) {

}

func parseDirStruture(path string) []string {
	return strings.Split(path, "/")

}
