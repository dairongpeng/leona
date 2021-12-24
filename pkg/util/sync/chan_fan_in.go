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

import "reflect"

// FanInReflect 多个chan 扇入到一个目标chan
func FanInReflect(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		// 构造SelectCase slice
		var cases []reflect.SelectCase
		for _, c := range chans {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		// 循环，从cases中选择一个可用的
		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok { // 此channel已经close
				cases = append(cases[:i], cases[i+1:]...)
				continue
			}
			out <- v.Interface()
		}
	}()
	return out
}

// FanInRec 二分递归加速多个chan往最终的一个chan扇入
func FanInRec(chans ...<-chan interface{}) <-chan interface{} {
	switch len(chans) {
	case 0:
		c := make(chan interface{})
		close(c)
		return c
	case 1:
		return chans[0]
	case 2:
		return mergeTwo(chans[0], chans[1])
	default:
		m := len(chans) / 2
		return mergeTwo(
			FanInRec(chans[:m]...),
			FanInRec(chans[m:]...))
	}
}

func mergeTwo(a, b <-chan interface{}) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		for a != nil || b != nil { //只要还有可读的chan
			select {
			case v, ok := <-a:
				if !ok { // a 已关闭，设置为nil
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok { // b 已关闭，设置为nil
					b = nil
					continue
				}
				c <- v
			}
		}
	}()
	return c
}
