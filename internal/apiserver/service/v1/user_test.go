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
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/AlekSi/pointer"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"
	gomock "github.com/golang/mock/gomock"

	"github.com/dairongpeng/leona/internal/apiserver/store"
	"github.com/dairongpeng/leona/internal/apiserver/store/fake"
)

func TestMain(m *testing.M) {
	_, _ = fake.GetFakeFactoryOr()
	os.Exit(m.Run())
}

func BenchmarkListUser(b *testing.B) {
	opts := metav1.ListOptions{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(50),
	}
	storeIns, _ := fake.GetFakeFactoryOr()
	u := &userService{
		store: storeIns,
	}

	for i := 0; i < b.N; i++ {
		//_, _ = u.ListWithBadPerformance(context.TODO(), opts)
		_, _ = u.List(context.TODO(), opts)
	}
}

func Test_newUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)

	type args struct {
		srv *service
	}
	tests := []struct {
		name string
		args args
		want *userService
	}{
		{
			name: "default",
			args: args{
				srv: &service{store: mockFactory},
			},
			want: &userService{
				store: mockFactory,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUsers(tt.args.srv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}
