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

package etcd

import (
	"context"
	"fmt"

	v1 "github.com/dairongpeng/leona/api/apiserver/v1"
	"github.com/dairongpeng/leona/pkg/errors"
	"github.com/dairongpeng/leona/pkg/json"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"
	"github.com/dairongpeng/leona/pkg/util/jsonutil"
)

type users struct {
	ds *datastore
}

func newUsers(ds *datastore) *users {
	return &users{ds: ds}
}

var keyUser = "/users/%v"

func (u *users) getKey(name string) string {
	return fmt.Sprintf(keyUser, name)
}

// Create creates a new user account.
func (u *users) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	return u.ds.Put(ctx, u.getKey(user.Name), jsonutil.ToString(user))
}

// Update updates an user account information.
func (u *users) Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error {
	return u.ds.Put(ctx, u.getKey(user.Name), jsonutil.ToString(user))
}

// Delete deletes the user by the user identifier.
func (u *users) Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error {

	if _, err := u.ds.Delete(ctx, u.getKey(username)); err != nil {
		return err
	}

	return nil
}

// Get return an user by the user identifier.
func (u *users) Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error) {
	resp, err := u.ds.Get(ctx, u.getKey(username))
	if err != nil {
		return nil, err
	}

	var user v1.User
	if err := json.Unmarshal(resp, &user); err != nil {
		return nil, errors.Wrap(err, "unmarshal to User struct failed")
	}

	return &user, nil
}

// List return all users.
func (u *users) List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	kvs, err := u.ds.List(ctx, u.getKey(""))
	if err != nil {
		return nil, err
	}

	ret := &v1.UserList{
		ListMeta: metav1.ListMeta{
			TotalCount: int64(len(kvs)),
		},
	}

	for _, v := range kvs {
		var user v1.User
		if err := json.Unmarshal(v.Value, &user); err != nil {
			return nil, errors.Wrap(err, "unmarshal to User struct failed")
		}

		ret.Items = append(ret.Items, &user)
	}

	return ret, nil
}
