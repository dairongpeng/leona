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

type policies struct {
	ds *datastore
}

func newPolicies(ds *datastore) *policies {
	return &policies{ds: ds}
}

var keyPolicy = "/policies/%v/%v"

func (p *policies) getKey(username string, name string) string {
	return fmt.Sprintf(keyPolicy, username, name)
}

// Create creates a new policy.
func (p *policies) Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOptions) error {
	return p.ds.Put(ctx, p.getKey(policy.Username, policy.Name), jsonutil.ToString(policy))
}

// Update updates an policy information.
func (p *policies) Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOptions) error {
	return p.ds.Put(ctx, p.getKey(policy.Username, policy.Name), jsonutil.ToString(policy))
}

// Delete deletes the policy by the policy identifier.
func (p *policies) Delete(ctx context.Context, username, name string, opts metav1.DeleteOptions) error {
	if _, err := p.ds.Delete(ctx, p.getKey(username, name)); err != nil {
		return err
	}

	return nil
}

// DeleteByUser deletes policies by username.
func (p *policies) DeleteByUser(ctx context.Context, username string, opts metav1.DeleteOptions) error {
	if _, err := p.ds.Delete(ctx, p.getKey(username, "")); err != nil {
		return err
	}

	return nil
}

// DeleteCollection batch deletes the policies.
func (p *policies) DeleteCollection(
	ctx context.Context,
	username string,
	names []string,
	opts metav1.DeleteOptions,
) error {
	return nil
}

// DeleteCollectionByUser batch deletes policies usernames.
func (p *policies) DeleteCollectionByUser(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	return nil
}

// Get return an policy by the policy identifier.
func (p *policies) Get(ctx context.Context, username, name string, opts metav1.GetOptions) (*v1.Policy, error) {
	resp, err := p.ds.Get(ctx, p.getKey(username, name))
	if err != nil {
		return nil, err
	}

	var policy v1.Policy
	if err := json.Unmarshal(resp, &policy); err != nil {
		return nil, errors.Wrap(err, "unmarshal to Policy struct failed")
	}

	return &policy, nil
}

// List return all policies.
func (p *policies) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.PolicyList, error) {
	kvs, err := p.ds.List(ctx, p.getKey(username, ""))
	if err != nil {
		return nil, err
	}

	ret := &v1.PolicyList{
		ListMeta: metav1.ListMeta{
			TotalCount: int64(len(kvs)),
		},
	}

	for _, v := range kvs {
		var policy v1.Policy
		if err := json.Unmarshal(v.Value, &policy); err != nil {
			return nil, errors.Wrap(err, "unmarshal to Policy struct failed")
		}

		ret.Items = append(ret.Items, &policy)
	}

	return ret, nil
}
