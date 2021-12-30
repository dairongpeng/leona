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

package fake

import (
	"fmt"
	"sync"

	v1 "github.com/dairongpeng/leona/api/apiserver/v1"
	"github.com/dairongpeng/leona/internal/apiserver/store"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"
)

// ResourceCount defines the number of fake resources.
const ResourceCount = 1000

type datastore struct {
	sync.RWMutex
	users []*v1.User
}

func (ds *datastore) Users() store.UserStore {
	return newUsers(ds)
}

func (ds *datastore) Close() error {
	return nil
}

var (
	fakeFactory store.Factory
	once        sync.Once
)

// GetFakeFactoryOr create fake store.
func GetFakeFactoryOr() (store.Factory, error) {
	once.Do(func() {
		fakeFactory = &datastore{
			users: FakeUsers(ResourceCount),
		}
	})

	if fakeFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v", fakeFactory)
	}

	return fakeFactory, nil
}

// FakeUsers returns fake user data.
func FakeUsers(count int) []*v1.User {
	// init some user records
	users := make([]*v1.User, 0)
	for i := 1; i <= count; i++ {
		users = append(users, &v1.User{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("user%d", i),
				ID:   uint64(i),
			},
			Nickname: fmt.Sprintf("user%d", i),
			Password: fmt.Sprintf("User%d@2020", i),
			Email:    fmt.Sprintf("user%d@qq.com", i),
		})
	}

	return users
}
