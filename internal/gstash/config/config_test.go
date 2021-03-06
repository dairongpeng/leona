// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
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

package config

import (
	"reflect"
	"testing"

	"github.com/dairongpeng/leona/internal/gstash/options"
)

func TestCreateConfigFromOptions(t *testing.T) {
	opts := options.NewOptions()
	type args struct {
		opts *options.Options
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				opts: opts,
			},
			want:    &Config{opts},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateConfigFromOptions(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateConfigFromOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateConfigFromOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
