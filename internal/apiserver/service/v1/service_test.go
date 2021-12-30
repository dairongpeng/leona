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

package v1

import (
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"

	"github.com/dairongpeng/leona/internal/apiserver/store"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)

	type args struct {
		store store.Factory
	}
	tests := []struct {
		name string
		args args
		want *service
	}{
		{
			name: "default",
			args: args{
				store: mockFactory,
			},
			want: &service{
				store: mockFactory,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Users(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)
	s := &service{
		store: mockFactory,
	}

	tests := []struct {
		name string
		want UserSrv
	}{
		{
			name: "default",
			want: newUsers(s),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.Users(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Users() = %v, want %v", got, tt.want)
			}
		})
	}
}
