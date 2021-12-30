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

// Package cache defines a cache service which can return all secrets and policies.
package cache

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/dairongpeng/leona/internal/apiserver/store"
)

func TestGetCacheInsOr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)

	type args struct {
		store store.Factory
	}
	tests := []struct {
		name    string
		args    args
		want    *Cache
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				store: mockFactory,
			},
			want: &Cache{
				store: mockFactory,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCacheInsOr(tt.args.store)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCacheInsOr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCacheInsOr() = %v, want %v", got, tt.want)
			}
		})
	}
}
