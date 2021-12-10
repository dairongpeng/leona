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

	v1 "github.com/dairongpeng/leona/api/apiserver/v1"
	"github.com/dairongpeng/leona/pkg/errors"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"

	"github.com/dairongpeng/leona/internal/apiserver/store"
	"github.com/dairongpeng/leona/internal/pkg/code"
)

// PolicySrv defines functions used to handle policy request.
type PolicySrv interface {
	Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOptions) error
	Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOptions) error
	Delete(ctx context.Context, username string, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, username string, names []string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, username string, name string, opts metav1.GetOptions) (*v1.Policy, error)
	List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.PolicyList, error)
}

type policyService struct {
	store store.Factory
}

var _ PolicySrv = (*policyService)(nil)

func newPolicies(srv *service) *policyService {
	return &policyService{store: srv.store}
}

func (s *policyService) Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOptions) error {
	if err := s.store.Policies().Create(ctx, policy, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *policyService) Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOptions) error {
	// Save changed fields.
	if err := s.store.Policies().Update(ctx, policy, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *policyService) Delete(ctx context.Context, username, name string, opts metav1.DeleteOptions) error {
	if err := s.store.Policies().Delete(ctx, username, name, opts); err != nil {
		return err
	}

	return nil
}

func (s *policyService) DeleteCollection(
	ctx context.Context,
	username string,
	names []string,
	opts metav1.DeleteOptions,
) error {
	if err := s.store.Policies().DeleteCollection(ctx, username, names, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *policyService) Get(ctx context.Context, username, name string, opts metav1.GetOptions) (*v1.Policy, error) {
	policy, err := s.store.Policies().Get(ctx, username, name, opts)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (s *policyService) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.PolicyList, error) {
	policies, err := s.store.Policies().List(ctx, username, opts)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return policies, nil
}
