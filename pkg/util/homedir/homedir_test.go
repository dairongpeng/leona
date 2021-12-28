package homedir

import (
	"fmt"
	"testing"
)

func TestHomeDir(t *testing.T) {
	dir := HomeDir()
	fmt.Println(dir)
}
