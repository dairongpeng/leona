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

import "sync"

// 基于chan实现信号量Semaphore 数据结构，并且还实现了Locker接口
type semaphore struct {
	sync.Locker
	ch chan struct{}
}

// NewSemaphore 创建一个新的信号量
func NewSemaphore(capacity int) sync.Locker {
	if capacity <= 0 {
		capacity = 1 // 容量为1就变成了一个互斥锁
	}
	return &semaphore{ch: make(chan struct{}, capacity)}
}

// Lock 请求一个资源
func (s *semaphore) Lock() {
	s.ch <- struct{}{}
}

// Unlock 释放资源
func (s *semaphore) Unlock() {
	<-s.ch
}
