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
	"github.com/ory/ladon"
	"gorm.io/gorm"

	"github.com/dairongpeng/leona/pkg/json"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"
	"github.com/dairongpeng/leona/pkg/util/idutil"
)

// AuthzPolicy defines leona policy type.
type AuthzPolicy struct {
	ladon.DefaultPolicy
}

// Policy represents a policy restful resource, include a ladon policy.
// It is also used as gorm model.
type Policy struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// The user of the policy.
	Username string `json:"username" gorm:"column:username" validate:"omitempty"`

	// AuthzPolicy policy, will not be stored in db.
	Policy AuthzPolicy `json:"policy,omitempty" gorm:"-" validate:"omitempty"`
	// Policy ladon.DefaultPolicy `json:"policy,omitempty" gorm:"-" validate:"omitempty"`

	// The ladon policy content, just a string format of ladon.DefaultPolicy. DO NOT modify directly.
	PolicyShadow string `json:"-" gorm:"column:policyShadow" validate:"omitempty"`
}

// PolicyList is the whole list of all policies which have been stored in stroage.
type PolicyList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	metav1.ListMeta `json:",inline"`

	// List of policies.
	Items []*Policy `json:"items"`
}

// TableName maps to mysql table name.
func (p *Policy) TableName() string {
	return "policy"
}

// String returns the string format of Policy.
func (ap AuthzPolicy) String() string {
	data, _ := json.Marshal(ap)

	return string(data)
}

// BeforeCreate run before create database record.
func (p *Policy) BeforeCreate(tx *gorm.DB) (err error) {
	p.PolicyShadow = p.Policy.String()
	p.ExtendShadow = p.Extend.String()

	return
}

// AfterCreate run after create database record.
func (p *Policy) AfterCreate(tx *gorm.DB) (err error) {
	p.InstanceID = idutil.GetInstanceID(p.ID, "policy-")

	return tx.Save(p).Error
}

// BeforeUpdate run before update database record.
func (p *Policy) BeforeUpdate(tx *gorm.DB) (err error) {
	p.PolicyShadow = p.Policy.String()
	p.ExtendShadow = p.Extend.String()

	return
}

// AfterFind run after find to unmarshal a policy string into ladon.DefaultPolicy struct.
func (p *Policy) AfterFind(tx *gorm.DB) (err error) {
	if err := json.Unmarshal([]byte(p.PolicyShadow), &p.Policy); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(p.ExtendShadow), &p.Extend); err != nil {
		return err
	}

	return
}
