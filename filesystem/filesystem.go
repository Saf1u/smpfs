package filesystem

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Saf1u/smpfs/disk"
)

type fileSystem struct {
	root item
	disk disk.Disk
}

type FileSystem interface {
	CreateDir(path string) error
	WriteFile(fileHandle File, data []byte) error
	ReadFile(fileHandle File) ([]byte, error)
	CreateFile(path string) error
	OpenFile(path string) (File, error)
	ListDir(path string) ([]string, error)
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

func NewFileSystem(disk disk.Disk) FileSystem {
	root := &directory{
		dirName:  "root",
		contents: map[string]item{},
	}
	return &fileSystem{root: root, disk: disk}
}

// CreateDir creates a directory in the nested tree structure.
// Will create the parent directories if they do not already exist.
func (f *fileSystem) CreateDir(path string) error {
	structure, err := parseDirStruture(path)
	if err != nil {
		return err
	}
	fmt.Println("CREATE DIR", structure, len(structure))
	if len(structure) > 1 {
		for i := 1; i < len(structure); i++ {
			if _, exist := f.root.(*directory).contents[structure[i-1]]; !exist {
				fmt.Println(structure[:i])
				err = f.root.(*directory).createDir(structure[:i])
				if err != nil {
					return fmt.Errorf("issue creating parent directory: /%s - %w", strings.Join(structure[:i], "/"), err)
				}
			}
		}
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

// ListDir lists filesystem dir contents
func (f *fileSystem) ListDir(path string) ([]string, error) {
	//add junk last path to levrage exisiting functionality that finds parent dir
	var structure []string
	var err error
	if path != "/" {
		path = fmt.Sprint(path, "/doesnotexist")

		structure, err = parseDirStruture(path)

		if err != nil {
			return nil, err
		}
	} else {
		structure = make([]string, 1)
		structure[0] = "junk"
	}

	return f.root.(*directory).listDir(structure)

}
