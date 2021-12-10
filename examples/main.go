package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Replace(strings.ToUpper("leona-apiserver"), "-", "_", -1))

	fmt.Println(strings.NewReplacer(".", "_", "-", "_").Replace("leona-apiserver-mysql-addr"))
}
