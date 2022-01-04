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

package sync

import (
	"fmt"
	"time"
)

type Token struct{}

// NewWorker 基于令牌信号传递goroutine的执行顺序
func NewWorker(id int, ch chan Token, nextCh chan Token) {
	for {
		token := <-ch       // 取得令牌
		fmt.Println(id + 1) // id从1开始
		time.Sleep(time.Second)
		nextCh <- token
	}
}
