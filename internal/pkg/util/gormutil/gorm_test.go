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

// Package gormutil is a util to convert offset and limit to default values.
package gormutil

import (
	"reflect"
	"testing"

	"github.com/AlekSi/pointer"
)

func TestUnpointer(t *testing.T) {
	type args struct {
		offset *int64
		limit  *int64
	}
	tests := []struct {
		name string
		args args
		want *LimitAndOffset
	}{
		{
			name: "both offset and limit are not zero",
			args: args{
				offset: pointer.ToInt64(0),
				limit:  pointer.ToInt64(10),
			},
			want: &LimitAndOffset{
				Offset: 0,
				Limit:  10,
			},
		},
		{
			name: "both offset and limit are zero",
			want: &LimitAndOffset{
				Offset: 0,
				Limit:  1000,
			},
		},
		{
			name: "offset not zero and limit zero",
			args: args{
				offset: pointer.ToInt64(2),
			},
			want: &LimitAndOffset{
				Offset: 2,
				Limit:  1000,
			},
		},
		{
			name: "offset zero and limit not zero",
			args: args{
				limit: pointer.ToInt64(10),
			},
			want: &LimitAndOffset{
				Offset: 0,
				Limit:  10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unpointer(tt.args.offset, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unpointer() = %v, want %v", got, tt.want)
			}
		})
	}
}
