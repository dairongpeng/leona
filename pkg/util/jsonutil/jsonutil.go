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

package jsonutil

import (
	"fmt"
	"strings"

	"github.com/dairongpeng/leona/pkg/json"
)

type JSONRawMessage []byte

func (m JSONRawMessage) Find(key string) JSONRawMessage {
	var objmap map[string]json.RawMessage
	err := json.Unmarshal(m, &objmap)
	if err != nil {
		fmt.Printf("Resolve JSON Key failed, find key =%s, err=%s",
			key, err)
		return nil
	}
	return JSONRawMessage(objmap[key])
}

func (m JSONRawMessage) ToList() []JSONRawMessage {
	var lists []json.RawMessage
	err := json.Unmarshal(m, &lists)
	if err != nil {
		fmt.Printf("Resolve JSON List failed, err=%s",
			err)
		return nil
	}
	var res []JSONRawMessage
	for _, v := range lists {
		res = append(res, JSONRawMessage(v))
	}
	return res
}

func (m JSONRawMessage) ToString() string {
	res := strings.ReplaceAll(string(m[:]), "\"", "")
	return res
}
