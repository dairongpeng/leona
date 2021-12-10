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
	"time"

	v1 "github.com/dairongpeng/leona/api/apiserver/v1"
	"github.com/dairongpeng/leona/pkg/errors"
	"github.com/dairongpeng/leona/pkg/fields"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"
	"github.com/dairongpeng/leona/pkg/util/stringutil"

	"github.com/dairongpeng/leona/internal/pkg/code"
	"github.com/dairongpeng/leona/internal/pkg/util/gormutil"
	reflectutil "github.com/dairongpeng/leona/internal/pkg/util/reflect"
)

type policies struct {
	ds *datastore
}

func newPolicies(ds *datastore) *policies {
	return &policies{ds}
}

// Create creates a new ladon policy.
func (p *policies) Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOptions) error {
	p.ds.Lock()
	defer p.ds.Unlock()

	for _, pol := range p.ds.policies {
		if pol.Username == policy.Username && pol.Name == policy.Name {
			return errors.New("record already exist")
		}
	}

	if len(p.ds.policies) > 0 {
		policy.ID = p.ds.policies[len(p.ds.policies)-1].ID + 1
	}
	p.ds.policies = append(p.ds.policies, policy)

	return nil
}

// Update updates policy by the policy identifier.
func (p *policies) Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOptions) error {
	p.ds.Lock()
	defer p.ds.Unlock()

	for _, pol := range p.ds.policies {
		if pol.Username == policy.Username && pol.Name == policy.Name {
			if _, err := reflectutil.CopyObj(policy, pol, nil); err != nil {
				return errors.Wrap(err, "copy policy failed")
			}
		}
	}

	return nil
}

// Delete deletes the policy by the policy identifier.
func (p *policies) Delete(ctx context.Context, username, name string, opts metav1.DeleteOptions) error {
	p.ds.Lock()
	defer p.ds.Unlock()

	policies := p.ds.policies
	p.ds.policies = make([]*v1.Policy, 0)
	for _, pol := range policies {
		if pol.Username == username && pol.Name == name {
			continue
		}

		p.ds.policies = append(p.ds.policies, pol)
	}

	return nil
}

// DeleteCollection batch deletes policies by policies ids.
func (p *policies) DeleteCollection(
	ctx context.Context,
	username string,
	names []string,
	opts metav1.DeleteOptions,
) error {
	p.ds.Lock()
	defer p.ds.Unlock()

	policies := p.ds.policies
	p.ds.policies = make([]*v1.Policy, 0)
	for _, pol := range policies {
		if pol.Username == username && stringutil.StringIn(pol.Name, names) {
			continue
		}

		p.ds.policies = append(p.ds.policies, pol)
	}

	return nil
}

func (p *policies) DeleteByUser(ctx context.Context, username string, opts metav1.DeleteOptions) error {
	p.ds.Lock()
	defer p.ds.Unlock()

	policies := p.ds.policies
	p.ds.policies = make([]*v1.Policy, 0)
	for _, pol := range policies {
		if pol.Username == username {
			continue
		}

		p.ds.policies = append(p.ds.policies, pol)
	}

	return nil
}

// DeleteCollectionByUser batch deletes policies usernames.
func (p *policies) DeleteCollectionByUser(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	p.ds.Lock()
	defer p.ds.Unlock()

	policies := p.ds.policies
	p.ds.policies = make([]*v1.Policy, 0)
	for _, pol := range policies {
		if stringutil.StringIn(pol.Username, usernames) {
			continue
		}

		p.ds.policies = append(p.ds.policies, pol)
	}

	return nil
}

// Get return policy by the policy identifier.
func (p *policies) Get(ctx context.Context, username, name string, opts metav1.GetOptions) (*v1.Policy, error) {
	p.ds.RLock()
	defer p.ds.RUnlock()

	for _, pol := range p.ds.policies {
		if pol.Username == username && pol.Name == name {
			return pol, nil
		}
	}

	return nil, errors.WithCode(code.ErrPolicyNotFound, "record not found")
}

// List return all policies.
func (p *policies) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.PolicyList, error) {
	p.ds.RLock()
	defer p.ds.RUnlock()

	ol := gormutil.Unpointer(opts.Offset, opts.Limit)
	selector, _ := fields.ParseSelector(opts.FieldSelector)
	name, _ := selector.RequiresExactMatch("name")

	policies := make([]*v1.Policy, 0)
	i := 0
	for _, pol := range p.ds.policies {
		if i == ol.Limit {
			break
		}

		if pol.Username != username {
			continue
		}

		if !strings.Contains(pol.Name, name) {
			continue
		}

		policies = append(policies, pol)
		i++
	}

	// Simulate database query latency, sleep 2 millisecond
	time.Sleep(2 * time.Millisecond)

	return &v1.PolicyList{
		ListMeta: metav1.ListMeta{
			TotalCount: int64(len(p.ds.policies)),
		},
		Items: policies,
	}, nil
}
