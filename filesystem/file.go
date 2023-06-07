package filesystem

import "github.com/Saf1u/smpfs/disk"

type file struct {
	fileName string
	info     disk.BlockRecord
}

func (fl *file) isFile() bool {
	return true
}
func (fl *file) name() string {
	return fl.fileName
}
