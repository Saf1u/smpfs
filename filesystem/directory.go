package filesystem

type directory struct {
	dirName  string
	contents map[string]item
}

func (dir *directory) isFile() bool {
	return false
}
func (dir *directory) name() string {
	return dir.dirName
}

func (dir *directory) createFile(levels []string) error {
	baseDir, err := dir.findParentDir(levels)
	if err != nil {
		return err
	}

	fileName := levels[len(levels)-1]
	if fsItem, exist := baseDir.(*directory).contents[fileName]; exist && fsItem.isFile() {
		return ErrFileAlreadyExist
	} else {
		newFile := NewFile(fileName).(item)
		baseDir.(*directory).contents[fileName] = newFile
		return nil
	}

}
func (dir *directory) openFile(levels []string) (File, error) {
	baseDir, err := dir.findParentDir(levels)
	if err != nil {
		return nil, err
	}

	fileName := levels[len(levels)-1]
	if fsItem, exist := baseDir.(*directory).contents[fileName]; exist && fsItem.isFile() {
		return fsItem.(File), nil
	} else {
		return nil, ErrFileDoesNotExist
	}

}

func (dir *directory) createDir(levels []string) error {
	baseDir, err := dir.findParentDir(levels)
	if err != nil {
		return err
	}

	folderName := levels[len(levels)-1]
	if fsItem, exist := baseDir.(*directory).contents[folderName]; exist && !fsItem.isFile() {
		return ErrDirrAlreadyExist
	} else {
		newDir := &directory{folderName, map[string]item{}}
		baseDir.(*directory).contents[folderName] = newDir
		return nil
	}

}
func (dir *directory) listDir(levels []string) ([]string, error) {
	baseDir, err := dir.findParentDir(levels)
	if err != nil {
		return nil, err
	}
	concDir := baseDir.(*directory)
	items := make([]string, 0)
	for names:= range concDir.contents {
		items = append(items, names)
	}
	return items, nil
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
