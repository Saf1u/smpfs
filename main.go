package main

import "github.com/Saf1u/smpfs/disk"

func main() {

	disk, err := disk.NewDisk(800*1024, 5)

	if err != nil {
		panic(err)
	}
	disk.SaveDisk()

}
