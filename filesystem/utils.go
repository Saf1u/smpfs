package filesystem

import (
	"regexp"
	"strings"
)

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
