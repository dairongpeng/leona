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

// Package cache defines a cache service which can return all secrets and policies.
package cache

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/dairongpeng/leona/api/proto/apiserver/v1"
	"github.com/dairongpeng/leona/pkg/errors"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"

	"github.com/dairongpeng/leona/internal/apiserver/store"
	"github.com/dairongpeng/leona/internal/pkg/code"
	"github.com/dairongpeng/leona/pkg/log"
)

// Cache defines a cache service used to list all secrets and policies.
type Cache struct {
	store store.Factory
}

var (
	cacheServer *Cache
	once        sync.Once
)

// GetCacheInsOr return cache server instance with given factory.
func GetCacheInsOr(store store.Factory) (*Cache, error) {
	if store != nil {
		once.Do(func() {
			cacheServer = &Cache{store}
		})
	}

	if cacheServer == nil {
		return nil, fmt.Errorf("got nil cache server")
	}

	return cacheServer, nil
}

// ListSecrets returns all secrets.
func (c *Cache) ListSecrets(ctx context.Context, r *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	log.L(ctx).Info("list secrets function called.")
	opts := metav1.ListOptions{
		Offset: r.Offset,
		Limit:  r.Limit,
	}

	secrets, err := c.store.Secrets().List(ctx, "", opts)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	items := make([]*pb.SecretInfo, 0)
	for _, secret := range secrets.Items {
		items = append(items, &pb.SecretInfo{
			SecretId:    secret.SecretID,
			Username:    secret.Username,
			SecretKey:   secret.SecretKey,
			Expires:     secret.Expires,
			Description: secret.Description,
			CreatedAt:   secret.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   secret.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &pb.ListSecretsResponse{
		TotalCount: secrets.TotalCount,
		Items:      items,
	}, nil
}

// ListPolicies returns all policies.
func (c *Cache) ListPolicies(ctx context.Context, r *pb.ListPoliciesRequest) (*pb.ListPoliciesResponse, error) {
	log.L(ctx).Info("list policies function called.")
	opts := metav1.ListOptions{
		Offset: r.Offset,
		Limit:  r.Limit,
	}

	policies, err := c.store.Policies().List(ctx, "", opts)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	items := make([]*pb.PolicyInfo, 0)
	for _, pol := range policies.Items {
		items = append(items, &pb.PolicyInfo{
			Name:         pol.Name,
			Username:     pol.Username,
			PolicyShadow: pol.PolicyShadow,
			CreatedAt:    pol.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &pb.ListPoliciesResponse{
		TotalCount: policies.TotalCount,
		Items:      items,
	}, nil
}
