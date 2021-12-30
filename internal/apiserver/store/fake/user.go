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
	"context"
	"strings"

	v1 "github.com/dairongpeng/leona/api/apiserver/v1"
	"github.com/dairongpeng/leona/pkg/errors"
	"github.com/dairongpeng/leona/pkg/fields"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"
	"github.com/dairongpeng/leona/pkg/util/stringutil"

	"github.com/dairongpeng/leona/internal/pkg/code"
	"github.com/dairongpeng/leona/internal/pkg/util/gormutil"
	reflectutil "github.com/dairongpeng/leona/internal/pkg/util/reflect"
)

type users struct {
	ds *datastore
}

func newUsers(ds *datastore) *users {
	return &users{ds}
}

// Create creates a new user account.
func (u *users) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	u.ds.Lock()
	defer u.ds.Unlock()

	for _, u := range u.ds.users {
		if u.Name == user.Name {
			return errors.WithCode(code.ErrUserAlreadyExist, "record already exist")
		}
	}

	if len(u.ds.users) > 0 {
		user.ID = u.ds.users[len(u.ds.users)-1].ID + 1
	}
	u.ds.users = append(u.ds.users, user)

	return nil
}

// Update updates an user account information.
func (u *users) Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error {
	u.ds.Lock()
	defer u.ds.Unlock()

	for _, u := range u.ds.users {
		if u.Name == user.Name {
			if _, err := reflectutil.CopyObj(user, u, nil); err != nil {
				return errors.Wrap(err, "copy user failed")
			}
		}
	}

	return nil
}

// Delete deletes the user by the user identifier.
func (u *users) Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error {
	u.ds.Lock()
	defer u.ds.Unlock()

	users := u.ds.users
	u.ds.users = make([]*v1.User, 0)
	for _, user := range users {
		if user.Name == username {
			continue
		}

		u.ds.users = append(u.ds.users, user)
	}

	return nil
}

// DeleteCollection batch deletes the users.
func (u *users) DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	u.ds.Lock()
	defer u.ds.Unlock()

	users := u.ds.users
	u.ds.users = make([]*v1.User, 0)
	for _, user := range users {
		if stringutil.StringIn(user.Name, usernames) {
			continue
		}

		u.ds.users = append(u.ds.users, user)
	}

	return nil
}

// Get return an user by the user identifier.
func (u *users) Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error) {
	u.ds.RLock()
	defer u.ds.RUnlock()

	for _, u := range u.ds.users {
		if u.Name == username {
			return u, nil
		}
	}

	return nil, errors.WithCode(code.ErrUserNotFound, "record not found")
}

// List return all users.
func (u *users) List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	u.ds.RLock()
	defer u.ds.RUnlock()

	ol := gormutil.Unpointer(opts.Offset, opts.Limit)
	selector, _ := fields.ParseSelector(opts.FieldSelector)
	username, _ := selector.RequiresExactMatch("name")

	users := make([]*v1.User, 0)
	i := 0
	for _, user := range u.ds.users {
		if i == ol.Limit {
			break
		}
		if !strings.Contains(user.Name, username) {
			continue
		}
		users = append(users, user)
		i++
	}

	return &v1.UserList{
		ListMeta: metav1.ListMeta{
			TotalCount: int64(len(u.ds.users)),
		},
		Items: users,
	}, nil
}
