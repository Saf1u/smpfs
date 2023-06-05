package main

import (
	"fmt"
	"regexp"
)

func main() {

	// disk, err := disk.NewDisk(800*1024, 5)

	// if err != nil {
	// 	panic(err)
	// }
	// disk.SaveDisk()
	match, err := regexp.MatchString(`^(/[^/ ]*)+/?$`, "/abc/")
	if err != nil {
		panic(err)
	}
	fmt.Println(match)
}
