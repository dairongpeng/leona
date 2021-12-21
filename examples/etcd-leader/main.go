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

// 导入所需的库
import (
	"bufio"
	"flag"
	"fmt"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"os"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 可以设置一些参数，比如节点ID
var (
	nodeID    = flag.Int("id", 0, "node ID")
	addr      = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	electName = flag.String("name", "my-test-elect", "election name")
)

func main() {
	flag.Parse()

	// 将etcd的地址解析成slice of string
	endpoints := strings.Split(*addr, ",")

	// 生成一个etcd的clien
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 创建session,如果程序宕机导致session断掉，etcd能检测到
	session, err := concurrency.NewSession(cli)
	defer session.Close()

	// 生成一个选举对象。下面主要使用它进行选举和查询等操作
	// 另一个方法ResumeElection可以使用既有的leader初始化Election
	e1 := concurrency.NewElection(session, *electName)

	// 从命令行读取命令
	consolescanner := bufio.NewScanner(os.Stdin)
	for consolescanner.Scan() {
		action := consolescanner.Text()
		switch action {
		case "elect": // 选举命令
			go elect(e1, *electName)
		case "proclaim": // 只更新leader的value
			proclaim(e1, *electName)
		case "resign": // 辞去leader,重新选举
			resign(e1, *electName)
		case "watch": // 监控leader的变动
			go watch(e1, *electName)
		case "query": // 查询当前的leader
			query(e1, *electName)
		case "rev":
			rev(e1, *electName)
		default:
			fmt.Println("unknown action")
		}
	}
}
