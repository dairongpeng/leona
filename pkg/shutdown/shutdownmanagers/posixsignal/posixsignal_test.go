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

package posixsignal

import (
	"syscall"
	"testing"
	"time"

	"github.com/dairongpeng/leona/pkg/shutdown"
)

type startShutdownFunc func(sm shutdown.ShutdownManager)

func (f startShutdownFunc) StartShutdown(sm shutdown.ShutdownManager) {
	f(sm)
}

func (f startShutdownFunc) ReportError(err error) {
}

func (f startShutdownFunc) AddShutdownCallback(shutdownCallback shutdown.ShutdownCallback) {
}

func waitSig(t *testing.T, c <-chan int) {
	select {
	case <-c:

	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for StartShutdown.")
	}
}

func TestStartShutdownCalledOnDefaultSignals(t *testing.T) {
	c := make(chan int, 100)

	psm := NewPosixSignalManager()
	_ = psm.Start(startShutdownFunc(func(sm shutdown.ShutdownManager) {
		c <- 1
	}))

	time.Sleep(time.Millisecond)

	_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	waitSig(t, c)

	_ = psm.Start(startShutdownFunc(func(sm shutdown.ShutdownManager) {
		c <- 1
	}))

	time.Sleep(time.Millisecond)

	_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	waitSig(t, c)
}

func TestStartShutdownCalledCustomSignal(t *testing.T) {
	c := make(chan int, 100)

	psm := NewPosixSignalManager(syscall.SIGHUP)
	_ = psm.Start(startShutdownFunc(func(sm shutdown.ShutdownManager) {
		c <- 1
	}))

	time.Sleep(time.Millisecond)

	_ = syscall.Kill(syscall.Getpid(), syscall.SIGHUP)

	waitSig(t, c)
}
