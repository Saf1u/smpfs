package filesystem

import (
	"testing"

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
