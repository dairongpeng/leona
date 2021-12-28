package fileutil

import (
	"fmt"
	"testing"
)

func TestListDir(t *testing.T) {
	dir := ListDir("/Users/dairongpeng")
	fmt.Println(dir)
}

func TestGetHomeDirectory(t *testing.T) {
	directory := GetHomeDirectory()
	fmt.Println(directory)
}

func TestGetParent(t *testing.T) {
	parent := GetParent("/Users/dairongpeng")
	fmt.Println(parent)
}
