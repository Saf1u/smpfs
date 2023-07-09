package main

import (
	"fmt"

	"github.com/Saf1u/smpfs/disk"
	"github.com/Saf1u/smpfs/filesystem"
)

func main() {
	disk, err := disk.NewDisk(100, 10)
	if err != nil {
		panic(err)
	}
	fs := filesystem.NewFileSystem(disk)
	err = fs.CreateDir("/home/usr/path")
	if err != nil {
		panic(err)
	}
	err = fs.CreateFile("/home/usr/path/file1.txt")
	if err != nil {
		panic(err)
	}
	err = fs.CreateFile("/home/usr/path/file2.txt")
	if err != nil {
		panic(err)
	}
	files, err := fs.ListDir("/home/usr/path")
	if err != nil {
		panic(err)
	}
	fmt.Println(files)
	fl, err := fs.OpenFile("/home/usr/path/file1.txt")
	if err != nil {
		panic(err)
	}
	err = fs.WriteFile(fl, []byte("random data that will pressist hopefully"))
	if err != nil {
		panic(err)
	}
	data, err := fs.ReadFile(fl)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	fl_b, err := fs.OpenFile("/home/usr/path/file2.txt")
	if err != nil {
		panic(err)
	}
	err = fs.WriteFile(fl_b, []byte("random data that will pressist hopefully part2"))
	if err != nil {
		panic(err)
	}
	data, err = fs.ReadFile(fl_b)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
