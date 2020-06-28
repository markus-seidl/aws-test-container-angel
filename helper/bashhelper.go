package helper

import (
	"regexp"
)

var (
	spaceRegexp = regexp.MustCompile("[\\s]+")
)

func SplitIntoArray(command string, args ...string) (output []string) {
	return spaceRegexp.Split(command, -1)
}
