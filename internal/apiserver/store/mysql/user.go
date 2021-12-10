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

package mysql

import (
	"context"

	v1 "github.com/dairongpeng/leona/api/apiserver/v1"
	"github.com/dairongpeng/leona/pkg/errors"
	"github.com/dairongpeng/leona/pkg/fields"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"
	gorm "gorm.io/gorm"

	"github.com/dairongpeng/leona/internal/pkg/code"
	"github.com/dairongpeng/leona/internal/pkg/util/gormutil"
)

type users struct {
	db *gorm.DB
}

func newUsers(ds *datastore) *users {
	return &users{db: ds.db}
}

// Create creates a new user account.
func (u *users) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	return u.db.Create(&user).Error
}

// Update updates an user account information.
func (u *users) Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error {
	return u.db.Save(user).Error
}

// Delete deletes the user by the user identifier.
func (u *users) Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error {
	// delete related policy first
	pol := newPolicies(&datastore{u.db})
	if err := pol.DeleteByUser(ctx, username, opts); err != nil {
		return err
	}

	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}

	err := u.db.Where("name = ?", username).Delete(&v1.User{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

// DeleteCollection batch deletes the users.
func (u *users) DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	// delete related policy first
	pol := newPolicies(&datastore{u.db})
	if err := pol.DeleteCollectionByUser(ctx, usernames, opts); err != nil {
		return err
	}

	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}

	return u.db.Where("name in (?)", usernames).Delete(&v1.User{}).Error
}

// Get return an user by the user identifier.
func (u *users) Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error) {
	user := &v1.User{}
	err := u.db.Where("name = ? and status = 1", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return user, nil
}

// List return all users.
func (u *users) List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	ret := &v1.UserList{}
	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	selector, _ := fields.ParseSelector(opts.FieldSelector)
	username, _ := selector.RequiresExactMatch("name")
	d := u.db.Where("name like ? and status = 1", "%"+username+"%").
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, d.Error
}

// ListOptional show a more graceful query method.
func (u *users) ListOptional(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	ret := &v1.UserList{}
	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	where := v1.User{}
	whereNot := v1.User{
		IsAdmin: 0,
	}
	selector, _ := fields.ParseSelector(opts.FieldSelector)
	username, found := selector.RequiresExactMatch("name")
	if found {
		where.Name = username
	}

	d := u.db.Where(where).
		Not(whereNot).
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, d.Error
}
