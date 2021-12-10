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
	"reflect"
	"testing"

	v1 "github.com/dairongpeng/leona/api/apiserver/v1"
	pb "github.com/dairongpeng/leona/api/proto/apiserver/v1"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"
	"github.com/golang/mock/gomock"

	"github.com/dairongpeng/leona/internal/apiserver/store"
	"github.com/dairongpeng/leona/internal/apiserver/store/fake"
)

func TestGetCacheInsOr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)

	type args struct {
		store store.Factory
	}
	tests := []struct {
		name    string
		args    args
		want    *Cache
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				store: mockFactory,
			},
			want: &Cache{
				store: mockFactory,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCacheInsOr(tt.args.store)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCacheInsOr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCacheInsOr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_ListSecrets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)
	mockSecretStore := store.NewMockSecretStore(ctrl)
	mockFactory.EXPECT().Secrets().Return(mockSecretStore)
	secrets := &v1.SecretList{
		ListMeta: metav1.ListMeta{
			TotalCount: 10,
		},
		Items: fake.FakeSecrets(3),
	}

	wantItems := make([]*pb.SecretInfo, 0)
	for _, secret := range secrets.Items {
		wantItems = append(wantItems, &pb.SecretInfo{
			SecretId:    secret.SecretID,
			Username:    secret.Username,
			SecretKey:   secret.SecretKey,
			Expires:     secret.Expires,
			Description: secret.Description,
			CreatedAt:   secret.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   secret.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	wantResponse := &pb.ListSecretsResponse{
		TotalCount: secrets.TotalCount,
		Items:      wantItems,
	}

	mockSecretStore.EXPECT().List(gomock.Any(), gomock.Eq(""), gomock.Any()).Return(secrets, nil)

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx context.Context
		r   *pb.ListSecretsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.ListSecretsResponse
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: mockFactory,
			},
			args: args{
				ctx: context.TODO(),
				r:   &pb.ListSecretsRequest{},
			},
			want:    wantResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				store: tt.fields.store,
			}
			got, err := c.ListSecrets(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cache.ListSecrets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.ListSecrets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_ListPolicies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)
	mockPolicyStore := store.NewMockPolicyStore(ctrl)
	mockFactory.EXPECT().Policies().Return(mockPolicyStore)
	policies := &v1.PolicyList{
		ListMeta: metav1.ListMeta{
			TotalCount: 10,
		},
		Items: fake.FakePolicies(3),
	}

	wantItems := make([]*pb.PolicyInfo, 0)
	for _, pol := range policies.Items {
		wantItems = append(wantItems, &pb.PolicyInfo{
			Name:         pol.Name,
			Username:     pol.Username,
			PolicyShadow: pol.PolicyShadow,
			CreatedAt:    pol.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	wantResponse := &pb.ListPoliciesResponse{
		TotalCount: policies.TotalCount,
		Items:      wantItems,
	}
	mockPolicyStore.EXPECT().List(gomock.Any(), gomock.Eq(""), gomock.Any()).Return(policies, nil)

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx context.Context
		r   *pb.ListPoliciesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.ListPoliciesResponse
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: mockFactory,
			},
			args: args{
				ctx: context.TODO(),
				r:   &pb.ListPoliciesRequest{},
			},
			want:    wantResponse,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				store: tt.fields.store,
			}
			got, err := c.ListPolicies(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cache.ListPolicies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.ListPolicies() = %v, want %v", got, tt.want)
			}
		})
	}
}
