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

type secrets struct {
	ds *datastore
}

func newSecrets(ds *datastore) *secrets {
	return &secrets{ds: ds}
}

var keySecret = "/secrets/%v/%v"

func (s *secrets) getKey(username string, secretID string) string {
	return fmt.Sprintf(keySecret, username, secretID)
}

// Create creates a new secret.
func (s *secrets) Create(ctx context.Context, secret *v1.Secret, opts metav1.CreateOptions) error {
	return s.ds.Put(ctx, s.getKey(secret.Username, secret.SecretID), jsonutil.ToString(secret))
}

// Update updates an secret information.
func (s *secrets) Update(ctx context.Context, secret *v1.Secret, opts metav1.UpdateOptions) error {
	return s.ds.Put(ctx, s.getKey(secret.Username, secret.SecretID), jsonutil.ToString(secret))
}

// Delete deletes the secret by the secret identifier.
func (s *secrets) Delete(ctx context.Context, username, secretID string, opts metav1.DeleteOptions) error {
	if _, err := s.ds.Delete(ctx, s.getKey(username, secretID)); err != nil {
		return err
	}

	return nil
}

// DeleteCollection batch deletes the secrets.
func (s *secrets) DeleteCollection(
	ctx context.Context,
	username string,
	secretIDs []string,
	opts metav1.DeleteOptions,
) error {
	return nil
}

// Get return an secret by the secret identifier.
func (s *secrets) Get(ctx context.Context, username, secretID string, opts metav1.GetOptions) (*v1.Secret, error) {
	resp, err := s.ds.Get(ctx, s.getKey(username, secretID))
	if err != nil {
		return nil, err
	}

	var secret v1.Secret
	if err := json.Unmarshal(resp, &secret); err != nil {
		return nil, errors.Wrap(err, "unmarshal to Secret struct failed")
	}

	return &secret, nil
}

// List return all secrets.
func (s *secrets) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error) {
	kvs, err := s.ds.List(ctx, s.getKey(username, ""))
	if err != nil {
		return nil, err
	}

	ret := &v1.SecretList{
		ListMeta: metav1.ListMeta{
			TotalCount: int64(len(kvs)),
		},
	}

	for _, v := range kvs {
		var secret v1.Secret
		if err := json.Unmarshal(v.Value, &secret); err != nil {
			return nil, errors.Wrap(err, "unmarshal to Secret struct failed")
		}

		ret.Items = append(ret.Items, &secret)
	}

	return ret, nil
}
