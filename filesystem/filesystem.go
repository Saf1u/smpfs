package filesystem

import (
	"errors"
	"regexp"
	"strings"

	"github.com/Saf1u/smpfs/disk"
)

type fileSystem struct {
	root item
	disk *disk.Disk
}
type itemTree struct {
	contents map[string]itemTree
}

var (
	ErrDirrAlreadyExistError  = errors.New("the directory already exists")
	ErrPathDoesNotExists      = errors.New("path to dir does not exists")
	ErrMalformedPathStructure = errors.New("the provided path is invalid")
)

type item interface {
	isFile() bool
	name() string
}

type directory struct {
	dirName  string
	contents map[string]item
}

type file struct {
	info disk.BlockRecord
}

func (f *fileSystem) CreateDir(path string) error {
	structure, err := parseDirStruture(path)
	if err != nil {
		return err
	}
	return f.root.(*directory).createDir(structure)

}
func (dir *directory) isFile() bool {
	return false
}
func (dir *directory) name() string {
	return dir.dirName
}

func (dir *directory) createDir(levels []string) error {
	baseDir, err := dir.findParentDir(levels)
	if err != nil {
		return err
	}

	folderName := levels[len(levels)-1]
	if _, exist := baseDir.(*directory).contents[folderName]; exist {
		return ErrDirrAlreadyExistError
	} else {
		newDir := &directory{folderName, map[string]item{}}
		baseDir.(*directory).contents[folderName] = newDir
		return nil
	}

}

func (dir *directory) findParentDir(levels []string) (item, error) {
	folderName := levels[0]
	if len(levels) == 1 {
		return dir, nil
	}
	if item, exist := dir.contents[folderName]; !exist || item.isFile() {
		return nil, ErrPathDoesNotExists
	}
	childDir := dir.contents[folderName].(*directory)
	return childDir.findParentDir(levels[1:])
}

func parseDirStruture(path string) ([]string, error) {
	if path[len(path)-1:] == "/" {
		path = path[:len(path)-1]
	}
	match, _ := regexp.MatchString(`^(/[^/ ]*)+/?$`, path)
	if !match {
		return nil, ErrMalformedPathStructure
	}

	return strings.Split(path, "/")[1:], nil
}
