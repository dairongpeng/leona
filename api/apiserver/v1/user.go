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
	"time"

	"github.com/dairongpeng/leona/pkg/auth"
	"github.com/dairongpeng/leona/pkg/json"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"
	"github.com/dairongpeng/leona/pkg/util/idutil"
	"gorm.io/gorm"
)

// User represents a user restful resource. It is also used as gorm model.
type User struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status int `json:"status" gorm:"column:status" validate:"omitempty"`

	// Required: true
	Nickname string `json:"nickname" gorm:"column:nickname" validate:"required,min=1,max=30"`

	// Required: true
	Password string `json:"password,omitempty" gorm:"column:password" validate:"required"`

	// Required: true
	Email string `json:"email" gorm:"column:email" validate:"required,email,min=1,max=100"`

	Phone string `json:"phone" gorm:"column:phone" validate:"omitempty"`

	IsAdmin int `json:"isAdmin,omitempty" gorm:"column:isAdmin" validate:"omitempty"`

	TotalPolicy int64 `json:"totalPolicy" gorm:"-" validate:"omitempty"`

	LoginedAt time.Time `json:"loginedAt,omitempty" gorm:"column:loginedAt"`
}

// UserList is the whole list of all users which have been stored in stroage.
type UserList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*User `json:"items"`
}

// TableName maps to mysql table name.
func (u *User) TableName() string {
	return "user"
}

// Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
func (u *User) Compare(pwd string) (err error) {
	err = auth.Compare(u.Password, pwd)

	return
}

// BeforeCreate run before create database record.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	return
}

// AfterCreate run after create database record.
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	u.InstanceID = idutil.GetInstanceID(u.ID, "user-")

	// Enable user
	u.Status = 1

	// NOTICE: tx.Save will trigger u.BeforeUpdate
	return tx.Save(u).Error
}

// BeforeUpdate run before update database record.
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.ExtendShadow = u.Extend.String()

	return err
}

// AfterFind run after find to unmarshal a extend shadown string into metav1.Extend struct.
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	if err := json.Unmarshal([]byte(u.ExtendShadow), &u.Extend); err != nil {
		return err
	}

	return nil
}
