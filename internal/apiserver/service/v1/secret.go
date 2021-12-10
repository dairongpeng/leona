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

// SecretSrv defines functions used to handle secret request.
type SecretSrv interface {
	Create(ctx context.Context, secret *v1.Secret, opts metav1.CreateOptions) error
	Update(ctx context.Context, secret *v1.Secret, opts metav1.UpdateOptions) error
	Delete(ctx context.Context, username, secretID string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, username string, secretIDs []string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, username, secretID string, opts metav1.GetOptions) (*v1.Secret, error)
	List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error)
}

type secretService struct {
	store store.Factory
}

var _ SecretSrv = (*secretService)(nil)

func newSecrets(srv *service) *secretService {
	return &secretService{store: srv.store}
}

func (s *secretService) Create(ctx context.Context, secret *v1.Secret, opts metav1.CreateOptions) error {
	if err := s.store.Secrets().Create(ctx, secret, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *secretService) Update(ctx context.Context, secret *v1.Secret, opts metav1.UpdateOptions) error {
	// Save changed fields.
	if err := s.store.Secrets().Update(ctx, secret, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *secretService) Delete(ctx context.Context, username, secretID string, opts metav1.DeleteOptions) error {
	if err := s.store.Secrets().Delete(ctx, username, secretID, opts); err != nil {
		return err
	}

	return nil
}

func (s *secretService) DeleteCollection(
	ctx context.Context,
	username string,
	secretIDs []string,
	opts metav1.DeleteOptions,
) error {
	if err := s.store.Secrets().DeleteCollection(ctx, username, secretIDs, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *secretService) Get(
	ctx context.Context,
	username, secretID string,
	opts metav1.GetOptions,
) (*v1.Secret, error) {
	secret, err := s.store.Secrets().Get(ctx, username, secretID, opts)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (s *secretService) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error) {
	secrets, err := s.store.Secrets().List(ctx, username, opts)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return secrets, nil
}
