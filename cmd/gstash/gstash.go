package main

import (
	"github.com/dairongpeng/leona/internal/gstash"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	gstash.NewApp("leona-gstash").Run()
}
