package filesystem

import (
	"testing"

	"github.com/Saf1u/smpfs/disk"
	"github.com/stretchr/testify/assert"
)

func TestParseDir(t *testing.T) {
	tests := []struct {
		name              string
		expectedErr       error
		path              string
		expectedStructure []string
	}{

		{name: "bad path",
			expectedErr:       ErrMalformedPathStructure,
			path:              "abcdedfg",
			expectedStructure: []string{},
		},
		{name: "succesful",
			expectedErr:       nil,
			path:              "/abcdedfg",
			expectedStructure: []string{"abcdedfg"},
		},
		{name: "succesful",
			expectedErr:       nil,
			path:              "/abcd/efg/hij",
			expectedStructure: []string{"abcd", "efg", "hij"},
		},
		{name: "succesful",
			expectedErr:       nil,
			path:              "/abcd/efg/hij/",
			expectedStructure: []string{"abcd", "efg", "hij"},
		},
	}
	for _, testcase := range tests {
		pathStructure, err := parseDirStruture(testcase.path)
		assert.Equal(t, testcase.expectedErr, err)
		if err == nil {
			assert.Equal(t, testcase.expectedStructure, pathStructure)
		}
	}
}

func TestFindParentDir(t *testing.T) {
	tests := []struct {
		name             string
		expectedErr      error
		pathStructure    []string
		expectedBasePath string
		setupFs          func() *fileSystem
	}{

		{
			name: "shouldLocateParentPath",
			//eg:/usr/home/desktop desktop parentDir is home
			pathStructure:    []string{"usr", "home", "desktop"},
			expectedBasePath: "home",
			expectedErr:      nil,
			setupFs: func() *fileSystem {
				desktop := &directory{
					dirName: "desktop",
				}
				home := &directory{
					dirName: "home",
					contents: map[string]item{
						"desktop": desktop,
					},
				}
				usr := &directory{
					dirName: "usr",
					contents: map[string]item{
						"home": home,
					},
				}
				root := &directory{
					dirName: "root",
					contents: map[string]item{
						"usr": usr,
					},
				}
				fs := &fileSystem{
					root: root,
				}
				return fs
			},
		},
		{
			name: "shouldLocateParentPath",
			//eg:/usr/home/desktop desktop parentDir is home
			pathStructure:    []string{"usr"},
			expectedBasePath: "root",
			expectedErr:      nil,
			setupFs: func() *fileSystem {

				usr := &directory{
					dirName: "usr",
				}
				root := &directory{
					dirName: "root",
					contents: map[string]item{
						"usr": usr,
					},
				}
				fs := &fileSystem{
					root: root,
				}
				return fs
			},
		},
		{
			name: "shouldLocateParentPath",
			//eg:/usr/home/desktop desktop parentDir is home
			pathStructure:    []string{"usr", "sus", "desktop"},
			expectedBasePath: "",
			expectedErr:      ErrPathDoesNotExists,
			setupFs: func() *fileSystem {
				desktop := &directory{
					dirName: "desktop",
				}
				home := &directory{
					dirName: "home",
					contents: map[string]item{
						"desktop": desktop,
					},
				}
				usr := &directory{
					dirName: "usr",
					contents: map[string]item{
						"home": home,
					},
				}
				root := &directory{
					dirName: "root",
					contents: map[string]item{
						"usr": usr,
					},
				}
				fs := &fileSystem{
					root: root,
				}
				return fs
			},
		},
	}
	for _, testcase := range tests {
		fs := testcase.setupFs()
		item, err := fs.root.(*directory).findParentDir(testcase.pathStructure)
		assert.Equal(t, testcase.expectedErr, err)
		if err == nil {
			assert.Equal(t, testcase.expectedBasePath, item.name())
		}
	}
}

func TestCreateDir(t *testing.T) {
	tests := []struct {
		name         string
		expectedErr  error
		pathName     string
		expectedName string
		setupFs      func() *fileSystem
	}{

		{
			name: "succesful create",
			//eg:/usr/home/desktop desktop parentDir is home
			pathName:     "/usr/home/desktop/newdir",
			expectedName: "newdir",
			expectedErr:  nil,
			setupFs: func() *fileSystem {
				desktop := &directory{
					dirName:  "desktop",
					contents: map[string]item{},
				}
				home := &directory{
					dirName: "home",
					contents: map[string]item{
						"desktop": desktop,
					},
				}
				usr := &directory{
					dirName: "usr",
					contents: map[string]item{
						"home": home,
					},
				}
				root := &directory{
					dirName: "root",
					contents: map[string]item{
						"usr": usr,
					},
				}
				fs := &fileSystem{
					root: root,
				}
				return fs
			},
		},
		{
			name: "creating a directory that exists",
			//eg:/usr/home/desktop desktop parentDir is home
			pathName:     "/usr",
			expectedName: "",
			expectedErr:  ErrDirrAlreadyExist,
			setupFs: func() *fileSystem {

				usr := &directory{
					dirName:  "usr",
					contents: map[string]item{},
				}
				root := &directory{
					dirName: "root",
					contents: map[string]item{
						"usr": usr,
					},
				}
				fs := &fileSystem{
					root: root,
				}
				return fs
			},
		},
	}
	for _, testcase := range tests {
		fs := testcase.setupFs()
		err := fs.CreateDir(testcase.pathName)
		assert.Equal(t, testcase.expectedErr, err)
	}
}

func TestCreateFile(t *testing.T) {
	tests := []struct {
		name         string
		expectedErr  error
		pathName     string
		expectedName string
		setupFs      func() *fileSystem
	}{

		{
			name: "succesful create",
			//eg:/usr/home/desktop desktop parentDir is home
			pathName:     "/usr/home/desktop/new.txt",
			expectedName: "new.txt",
			expectedErr:  nil,
			setupFs: func() *fileSystem {
				desktop := &directory{
					dirName:  "desktop",
					contents: map[string]item{},
				}
				home := &directory{
					dirName: "home",
					contents: map[string]item{
						"desktop": desktop,
					},
				}
				usr := &directory{
					dirName: "usr",
					contents: map[string]item{
						"home": home,
					},
				}
				root := &directory{
					dirName: "root",
					contents: map[string]item{
						"usr": usr,
					},
				}
				fs := &fileSystem{
					root: root,
				}
				return fs
			},
		},
		{
			name: "creating a file that exists",
			//eg:/usr/home/desktop desktop parentDir is home
			pathName:     "/home/usr.txt",
			expectedName: "",
			expectedErr:  ErrFileAlreadyExist,
			setupFs: func() *fileSystem {

				usr := &file{
					fileName: "usr.txt",
				}
				home := &directory{
					dirName: "home",
					contents: map[string]item{
						"usr.txt": usr,
					},
				}
				root := &directory{
					dirName: "root",
					contents: map[string]item{
						"home": home,
					},
				}
				fs := &fileSystem{
					root: root,
				}
				return fs
			},
		},
	}
	for _, testcase := range tests {
		fs := testcase.setupFs()
		err := fs.CreateFile(testcase.pathName)
		assert.Equal(t, testcase.expectedErr, err)
	}
}

func TestOpenFile(t *testing.T) {
	tests := []struct {
		name         string
		expectedErr  error
		pathName     string
		expectedName string
		setupFs      func() *fileSystem
	}{

		{
			name: "file cannot be found",
			//eg:/usr/home/desktop desktop parentDir is home
			pathName:     "/usr/home/desktop/new.txt",
			expectedName: "",
			expectedErr:  ErrFileDoesNotExist,
			setupFs: func() *fileSystem {
				desktop := &directory{
					dirName: "desktop",
					contents: map[string]item{
						"file.txt": &file{
							fileName: "file.txt",
						},
					},
				}
				home := &directory{
					dirName: "home",
					contents: map[string]item{
						"desktop": desktop,
					},
				}
				usr := &directory{
					dirName: "usr",
					contents: map[string]item{
						"home": home,
					},
				}
				root := &directory{
					dirName: "root",
					contents: map[string]item{
						"usr": usr,
					},
				}
				fs := &fileSystem{
					root: root,
				}
				return fs
			},
		},
		{
			name: "accesing a file that exists",
			//eg:/usr/home/desktop desktop parentDir is home
			pathName:     "/home/usr.txt",
			expectedName: "",
			expectedErr:  nil,
			setupFs: func() *fileSystem {

				usr := &file{
					fileName: "usr.txt",
				}
				home := &directory{
					dirName: "home",
					contents: map[string]item{
						"usr.txt": usr,
					},
				}
				root := &directory{
					dirName: "root",
					contents: map[string]item{
						"home": home,
					},
				}
				fs := &fileSystem{
					root: root,
				}
				return fs
			},
		},
	}
	for _, testcase := range tests {
		fs := testcase.setupFs()
		_, err := fs.OpenFile(testcase.pathName)
		assert.Equal(t, testcase.expectedErr, err)

	}
}

func TestWriteFile(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr error
		dataToWrite string
		pathName    string
		setupFs     func() *fileSystem
	}{

		{
			name: "writing succesfully to a file that exists",
			//eg:/usr/home/desktop desktop parentDir is home
			pathName:    "/home/usr.txt",
			dataToWrite: "random text",
			expectedErr: nil,
			setupFs: func() *fileSystem {
				usr := &file{
					fileName: "usr.txt",
					info:     &disk.BlockRecord{},
				}
				home := &directory{
					dirName: "home",
					contents: map[string]item{
						"usr.txt": usr,
					},
				}
				root := &directory{
					dirName: "root",
					contents: map[string]item{
						"home": home,
					},
				}
				disk, _ := disk.NewDisk(100, 10)
				fs := &fileSystem{
					root: root,
					disk: disk,
				}
				return fs
			},
		},

		{
			name: "writing unsuccesfully to a file due to disk issues",
			//eg:/usr/home/desktop desktop parentDir is home
			pathName:    "/home/usr.txt",
			dataToWrite: "random text",
			expectedErr: ErrFileCouldNotBeWritten,
			setupFs: func() *fileSystem {
				usr := &file{
					fileName: "usr.txt",
					info:     &disk.BlockRecord{},
				}
				home := &directory{
					dirName: "home",
					contents: map[string]item{
						"usr.txt": usr,
					},
				}
				root := &directory{
					dirName: "root",
					contents: map[string]item{
						"home": home,
					},
				}
				disk, _ := disk.NewDisk(3, 3)
				fs := &fileSystem{
					root: root,
					disk: disk,
				}
				return fs
			},
		},
	}
	for _, testcase := range tests {
		fs := testcase.setupFs()
		file, _ := fs.OpenFile(testcase.pathName)
		err := fs.WriteFile(file, []byte(testcase.dataToWrite))
		assert.Equal(t, testcase.expectedErr, err)

	}
}
