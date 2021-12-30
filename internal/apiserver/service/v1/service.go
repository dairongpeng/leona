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

//go:generate mockgen -self_package=github.com/dairongpeng/leona/internal/parseserver/service/v1 -destination mock_service.go -package v1 github.com/dairongpeng/leona/internal/parseserver/service/v1 Service,UserSrv,SecretSrv,PolicySrv

import "github.com/dairongpeng/leona/internal/apiserver/store"

// Service defines functions used to return resource interface.
type Service interface {
	Users() UserSrv
}

type service struct {
	store store.Factory
}

// NewService returns Service interface.
func NewService(store store.Factory) Service {
	return &service{
		store: store,
	}
}

func (s *service) Users() UserSrv {
	return newUsers(s)
}
