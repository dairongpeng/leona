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

package errors

import (
	"encoding/json"
	"regexp"
	"testing"
)

func TestFrameMarshalText(t *testing.T) {
	var tests = []struct {
		Frame
		want string
	}{{
		initpc,
		`^github.com/dairongpeng/leona/pkg/errors\.init(\.ializers)? .+/github\.com/dairongpeng/leona/pkg/errors/stack_test.go:\d+$`,
	}, {
		0,
		`^unknown$`,
	}}
	for i, tt := range tests {
		got, err := tt.Frame.MarshalText()
		if err != nil {
			t.Fatal(err)
		}
		if !regexp.MustCompile(tt.want).Match(got) {
			t.Errorf("test %d: MarshalJSON:\n got %q\n want %q", i+1, string(got), tt.want)
		}
	}
}

func TestFrameMarshalJSON(t *testing.T) {
	var tests = []struct {
		Frame
		want string
	}{{
		initpc,
		`^"github\.com/dairongpeng/leona/pkg/errors\.init(\.ializers)? .+/github\.com/dairongpeng/leona/pkg/errors/stack_test.go:\d+"$`,
	}, {
		0,
		`^"unknown"$`,
	}}
	for i, tt := range tests {
		got, err := json.Marshal(tt.Frame)
		if err != nil {
			t.Fatal(err)
		}
		if !regexp.MustCompile(tt.want).Match(got) {
			t.Errorf("test %d: MarshalJSON:\n got %q\n want %q", i+1, string(got), tt.want)
		}
	}
}
