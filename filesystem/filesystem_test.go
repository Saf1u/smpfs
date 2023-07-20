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

func TestRead(t *testing.T) {
	tests := []struct {
		name         string
		expectedErr  error
		dataToExpect string
		pathName     string
		setupFs      func(t *testing.T) *fileSystem
	}{

		{
			name: "reading succesfully from a file that exists",
			//eg:/usr/home/desktop desktop parentDir is home
			pathName:     "/home/usr.txt",
			dataToExpect: "random text here in the block",
			expectedErr:  nil,
			setupFs: func(t *testing.T) *fileSystem {
				usr := NewFile("usr.txt").(item)
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
				fl, err := fs.OpenFile("/home/usr.txt")
				if err != nil {
					t.Fatal(err)
				}
				err = fs.WriteFile(fl, []byte("random text here in the block"))
				if err != nil {
					t.Fatal(err)
				}
				return fs
			},
		},
	}
	for _, testcase := range tests {
		fs := testcase.setupFs(t)
		file, _ := fs.OpenFile(testcase.pathName)
		data, err := fs.ReadFile(file)
		assert.Equal(t, testcase.expectedErr, err)
		assert.Equal(t, []byte(testcase.dataToExpect), data)

	}
}

func TestListsDir(t *testing.T) {
	tests := []struct {
		name              string
		expectedItemNames []string
		expectedError     error
		pathName          string
		setupFs           func(t *testing.T) *fileSystem
	}{

		{
			name: "listing files from an exisiting dir",
			//eg:/usr/home/desktop desktop parentDir is home
			pathName:          "/home",
			expectedItemNames: []string{"textfile1.txt", "textfile2.txt"},
			setupFs: func(t *testing.T) *fileSystem {
				textfilea := NewFile("textfile1.txt").(item)
				textfileb := NewFile("textfile2.txt").(item)
				home := &directory{
					dirName: "home",
					contents: map[string]item{
						"textfile1.txt": textfilea,
						"textfile2.txt": textfileb,
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
			name: "listing folders",
			//eg:/usr/home/desktop desktop parentDir is home
			pathName:          "/home",
			expectedItemNames: []string{"images", "videos"},
			setupFs: func(t *testing.T) *fileSystem {
				imageDir := &directory{
					dirName: "images",
				}
				videoDir := &directory{
					dirName: "videos",
				}
				home := &directory{
					dirName: "home",
					contents: map[string]item{
						"images": imageDir,
						"videos": videoDir,
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
			name: "empty dir",
			//eg:/usr/home/desktop desktop parentDir is home
			pathName:          "/home",
			expectedItemNames: []string{},
			setupFs: func(t *testing.T) *fileSystem {

				home := &directory{
					dirName: "home",
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
			name: "non existing path",
			//eg:/usr/home/desktop desktop parentDir is home
			pathName:          "/home/files",
			expectedItemNames: []string{},
			expectedError:     ErrPathDoesNotExists,
			setupFs: func(t *testing.T) *fileSystem {

				home := &directory{
					dirName: "home",
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
	}
	for _, testcase := range tests {
		fs := testcase.setupFs(t)
		items, err := fs.ListDir(testcase.pathName)
		if err == nil {
			assert.ElementsMatch(t, testcase.expectedItemNames, items)
		}
		assert.Equal(t, testcase.expectedError, err)

	}
}

func TestDeleteFile(t *testing.T) {
	tests := []struct {
		name                        string
		fileName                    string
		expectedError               error
		fileSize                    int
		diskSize                    int
		diskSizeBeforeDelete        int
		expectedDiskSizeAfterDelete int
		setupFs                     func(t *testing.T, fileSize int, diskSize int) *fileSystem
	}{

		{
			name:     "deleting a file that exists",
			fileName: "/home/textfile1.txt",
			setupFs: func(t *testing.T, fileSize int, diskSize int) *fileSystem {
				textfilea := NewFile("textfile1.txt")
				home := &directory{
					dirName: "home",
					contents: map[string]item{
						"textfile1.txt": textfilea,
					},
				}
				root := &directory{
					dirName: "root",
					contents: map[string]item{
						"home": home,
					},
				}
				disk, _ := disk.NewDisk(diskSize, 10)
				record, err := disk.Write(make([]byte, 20))
				if err != nil {
					t.Fatal(err)
				}
				textfilea.setManifest(record)
				fs := &fileSystem{
					root: root,
					disk: disk,
				}
				return fs
			},
			fileSize:                    20,
			diskSize:                    100,
			expectedDiskSizeAfterDelete: 100,
			diskSizeBeforeDelete:        80,
		},

		{
			name:     "deleting a file that does exists",
			fileName: "/home/textfiledoesnotexist.txt",
			setupFs: func(t *testing.T, fileSize int, diskSize int) *fileSystem {
				textfilea := NewFile("textfileexist.txt")
				home := &directory{
					dirName: "home",
					contents: map[string]item{
						"textfile1.txt": textfilea,
					},
				}
				root := &directory{
					dirName: "root",
					contents: map[string]item{
						"home": home,
					},
				}
				disk, _ := disk.NewDisk(diskSize, 10)
				record, err := disk.Write(make([]byte, 20))
				if err != nil {
					t.Fatal(err)
				}
				textfilea.setManifest(record)
				fs := &fileSystem{
					root: root,
					disk: disk,
				}
				return fs
			},
			fileSize:                    20,
			diskSize:                    100,
			expectedDiskSizeAfterDelete: 80,
			diskSizeBeforeDelete:        80,
			expectedError: ErrFileDoesNotExist,
		},
	}
	for _, testcase := range tests {
		fs := testcase.setupFs(t, testcase.fileSize, testcase.diskSize)
		sizeBefore := fs.GetAvailableMemory()
		err := fs.DeleteFile(testcase.fileName)
		if err == nil {
			assert.Equal(t, testcase.diskSizeBeforeDelete, sizeBefore)
			assert.Equal(t, testcase.expectedDiskSizeAfterDelete, fs.GetAvailableMemory())
			_, err := fs.OpenFile(testcase.fileName)
			assert.Equal(t, err, ErrFileDoesNotExist)
		}
		assert.Equal(t, testcase.expectedError, err)

	}
}
