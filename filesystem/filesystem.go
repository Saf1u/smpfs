package filesystem

import (
	"errors"

	"github.com/Saf1u/smpfs/disk"
)

type fileSystem struct {
	root item
	disk *disk.Disk
}

var (
	ErrDirrAlreadyExist       = errors.New("the directory already exists")
	ErrPathDoesNotExists      = errors.New("path to dir does not exists")
	ErrMalformedPathStructure = errors.New("the provided path is invalid")
	ErrFileAlreadyExist       = errors.New("the file already exists")
)

type item interface {
	isFile() bool
	name() string
}

func (f *fileSystem) CreateDir(path string) error {
	structure, err := parseDirStruture(path)
	if err != nil {
		return err
	}
	return f.root.(*directory).createDir(structure)

}

func (f *fileSystem) CreateFile(path string) error {
	structure, err := parseDirStruture(path)
	if err != nil {
		return err
	}
	return f.root.(*directory).createFile(structure)

}
