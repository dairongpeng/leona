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
	"fmt"
	"golang.org/x/sync/semaphore"
	"log"
	"runtime"
	"time"
)

//type Weighted struct {
//	size    int64         // 最大资源数
//	cur     int64         // 当前已被使用的资源
//	mu      sync.Mutex    // 互斥锁，对字段的保护
//	waiters list.List     // 等待队列
//}

var (
	maxWorkers = runtime.GOMAXPROCS(0) // worker数量 CPU的核数
	// 使用信号量遵循的原则就是请求多少资源，就释放多少资源
	// 你一定要注意，必须使用正确的方法传递整数，不要“耍小聪明”，而且，请求的资源数一定不要超过最大资源数。
	sema = semaphore.NewWeighted(int64(maxWorkers)) //信号量
	task = make([]int, maxWorkers*4)                // 任务数，是worker的四倍
)

// main 函数在信号量的处理过程中类似于调度员dispatcher
func main() {
	ctx := context.Background()

	// 多个任务可以通过信号量来进行缓冲
	for i := range task {
		// Acquire 相当于 P 操作，你可以一次获取多个资源，如果没有足够多的资源，调用者就会被阻塞。
		// 它的第一个参数是 Context，这就意味着，你可以通过 Context 增加超时或者 cancel 的机制。
		// 如果是正常获取了资源，就返回 nil；否则，就返回 ctx.Err()，信号量不改变。
		if err := sema.Acquire(ctx, 1); err != nil {
			break
		}

		// TryAcquire
		// sema.TryAcquire(2)

		// 启动worker goroutine
		go func(i int) {
			// Release 相当于 V 操作，可以将 n 个资源释放，返还给信号量。
			defer sema.Release(1)
			time.Sleep(100 * time.Millisecond) // 模拟一个耗时操作
			task[i] = i + 1
		}(i)
	}

	// 请求所有的worker,这样能确保前面的worker都执行完
	// 如果在实际应用中，你想等所有的 Worker 都执行完，就可以获取最大计数值的信号量。
	if err := sema.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.Printf("获取所有的worker失败: %v", err)
	}

	fmt.Println(task)
}
