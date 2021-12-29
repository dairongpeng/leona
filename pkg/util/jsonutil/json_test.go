package jsonutil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type demo struct {
	S   string
	Arr []string
	I   int
	M   map[string]string
}

func TestToJson(t *testing.T) {
	arr := make([]string, 0)
	arr = append(arr, "t1")
	arr = append(arr, "t2")
	m := make(map[string]string, 0)
	m["d1"] = "aa"
	m["d2"] = "bb"
	d := demo{
		S:   "Demo",
		Arr: arr,
		I:   99,
		M:   m,
	}

	// trans to simpleJson
	json := ToJson(d)
	pretty, err := json.EncodePretty()
	assert.NoError(t, err)
	fmt.Println(string(pretty))
}
