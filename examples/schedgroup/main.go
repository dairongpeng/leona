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
	"github.com/mdlayher/schedgroup"
	"log"
	"time"
)

func main() {
	sg := schedgroup.New(context.Background())

	// 设置子任务分别在100、200、300ms之后执行
	for i := 0; i < 3; i++ {
		n := i + 1
		sg.Delay(time.Duration(n)*1000*time.Millisecond, func() {
			log.Println(n) //输出任务编号
		})
	}

	// 等待所有的子任务都完成
	if err := sg.Wait(); err != nil {
		log.Fatalf("failed to wait: %v", err)
	}
}
