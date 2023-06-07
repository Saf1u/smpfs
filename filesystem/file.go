package filesystem

import "github.com/Saf1u/smpfs/disk"

type file struct {
	fileName string
	info     disk.BlockRecord
}

type File interface {
	Read() ([]byte, error)
	Close() error
	Append([]byte) error
}

func (fl *file) isFile() bool {
	return true
}
func (fl *file) name() string {
	return fl.fileName
}

func (fl *file) Read() ([]byte, error) {
	return nil, nil
}
func (fl *file) Close() error {
	return nil
}
func (fl *file) Append([]byte) error {
	return nil
}
