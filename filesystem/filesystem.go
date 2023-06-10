package filesystem

import (
	"errors"
	"time"

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
	ErrFileDoesNotExist       = errors.New("the file does not exists")
	ErrFileCouldNotBeWritten  = errors.New("not enough emmoey to write to files")
	ErrUnkonwnError           = errors.New("???")
)

type item interface {
	isFile() bool
	name() string
}

// CreateDir creates a directory in the nested tree structure,it does not create all parent paths of the final path
func (f *fileSystem) CreateDir(path string) error {
	structure, err := parseDirStruture(path)
	if err != nil {
		return err
	}
	return f.root.(*directory).createDir(structure)

}

// WriteFile truncates the file and writes the data to the file
func (f *fileSystem) WriteFile(fileHandle File, data []byte) error {

	//Truncate File
	if fileHandle.getManifest() != nil {
		f.disk.Delete(fileHandle.getManifest())
	}
	fileManifest, err := f.disk.Write(data)
	if err != nil {
		if errors.Is(err, disk.ErrInsufficentMemoryError) {
			return ErrFileCouldNotBeWritten
		} else {
			return ErrUnkonwnError
		}
	}
	fileHandle.setManifest(fileManifest)
	fileHandle.updateAccessTs(time.Now())
	return nil
}

// ReadFile reads the data stored in the file
func (f *fileSystem) ReadFile(fileHandle File) ([]byte, error) {

	data, err := f.disk.Read(fileHandle.getManifest())
	if err != nil {
		return nil, err
	}
	fileHandle.updateAccessTs(time.Now())
	return data, nil
}

// CreateFile creates a file in the nested tree structure,it does not create all parent paths of the final path
func (f *fileSystem) CreateFile(path string) error {
	structure, err := parseDirStruture(path)
	if err != nil {
		return err
	}
	return f.root.(*directory).createFile(structure)

}

// OpenFile Searches for a file in the directory structure, and returns a file pointer to enable reads and writes
func (f *fileSystem) OpenFile(path string) (File, error) {
	structure, err := parseDirStruture(path)
	if err != nil {
		return nil, err
	}
	return f.root.(*directory).openFile(structure)

}
