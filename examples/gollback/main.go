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
	"errors"
	"fmt"
	"github.com/vardius/gollback"
	"time"
)

func main() {
	AllTest()
	RetryTest()
}

func AllTest() {
	rs, errs := gollback.All( // 调用All方法
		context.Background(),
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(3 * time.Second)
			return 1, nil // 第一个任务没有错误，返回1
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("failed") // 第二个任务返回一个错误
		},
		func(ctx context.Context) (interface{}, error) {
			return 3, nil // 第三个任务没有错误，返回3
		},
	)

	fmt.Println(rs)   // 输出子任务的结果
	fmt.Println(errs) // 输出子任务的错误信息
}

// RetryTest Retry 不是执行一组子任务，而是执行一个子任务。如果子任务执行失败，它会尝试一定的次数，如果一直不成功 ，就会返回失败错误 ，如果执行成功，它会立即返回。
//  如果 retires 等于 0，它会永远尝试，直到成功。
func RetryTest() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试5次，或者超时返回
	res, err := gollback.Retry(ctx, 5, func(ctx context.Context) (interface{}, error) {
		return nil, errors.New("failed")
	})

	fmt.Println(res) // 输出结果
	fmt.Println(err) // 输出错误信息
}
