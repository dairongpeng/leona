// Copyright 2021 dairongpeng <dairongpeng@foxmail.com>. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"github.com/dairongpeng/leona/pkg/log"
)

var (
	h bool

	level  int
	format string
)

func main() {
	flag.BoolVar(&h, "h", false, "Print this help.")
	flag.IntVar(&level, "l", 0, "Log level.")
	flag.StringVar(&format, "f", "console", "log output format.")

	flag.Parse()

	if h {
		flag.Usage()

		return
	}

	// logger配置
	opts := &log.Options{
		Level:            "debug",
		Format:           "console",
		EnableColor:      true,
		DisableCaller:    true,
		OutputPaths:      []string{"test.log", "stdout"},
		ErrorOutputPaths: []string{"error.log"},
	}
	// 初始化全局logger
	log.Init(opts)
	defer log.Flush()

	// WithValues使用
	lv := log.WithValues("X-Request-ID", "7a7b9f24-4cae-4b2a-9464-69088b45b904")

	// Context使用
	lv.Infof("Start to call pirntString function")
	ctx := lv.WithContext(context.Background())
	pirntString(ctx, "World")
}

func pirntString(ctx context.Context, str string) {
	lc := log.FromContext(ctx)
	lc.Infof("Hello %s", str)
}
