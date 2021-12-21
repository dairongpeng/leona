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
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"math/rand"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	addr     = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	lockName = flag.String("name", "my-test-lock", "lock name")
)

// 类似zookeeper 的分布式锁原理，节点宕机对应session 销毁，持有的锁会被释放
func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	// etcd地址
	endpoints := strings.Split(*addr, ",")
	// 生成一个etcd client
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	useLock(cli) // 测试锁
}

// useLock 使用lock方式实现分布式锁
func useLock(cli *clientv3.Client) {
	// 为锁生成session
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	//得到一个分布式锁
	locker := concurrency.NewLocker(s1, *lockName)

	// 请求锁
	log.Println("acquiring lock")
	locker.Lock()
	log.Println("acquired lock")

	// 等待一段时间
	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)
	locker.Unlock() // 释放锁

	log.Println("released lock")
}

// useMutex 使用Mutex实现分布式锁，需要ctx
func useMutex(cli *clientv3.Client) {
	// 为锁生成session
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	m1 := concurrency.NewMutex(s1, *lockName)

	//在请求锁之前查询key
	log.Printf("before acquiring. key: %s", m1.Key())
	// 请求锁
	log.Println("acquiring lock")
	if err := m1.Lock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	log.Printf("acquired lock. key: %s", m1.Key())

	//等待一段时间
	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)

	// 释放锁
	if err := m1.Unlock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	log.Println("released lock")
}
