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
	"time"
)

// Mutex 使用chan实现互斥锁
type Mutex struct {
	ch chan struct{}
}

// NewMutex 使用锁需要初始化
func NewMutex() *Mutex {
	mu := &Mutex{make(chan struct{}, 1)} // 容量为1的chan
	mu.ch <- struct{}{}                  // 写入一个interface当做是锁
	return mu
}

// Lock 请求锁，直到获取到
func (m *Mutex) Lock() {
	<-m.ch
}

// Unlock 解锁
func (m *Mutex) Unlock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")
	}
}

// TryLock 尝试获取锁
func (m *Mutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
	}
	return false
}

// LockTimeout 加入一个超时的设置
func (m *Mutex) LockTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case <-m.ch:
		timer.Stop()
		return true
	case <-timer.C:
	}
	return false
}

// IsLocked 锁是否已被持有
func (m *Mutex) IsLocked() bool {
	return len(m.ch) == 0
}
