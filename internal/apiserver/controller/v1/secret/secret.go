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

package secret

import (
	srvv1 "github.com/dairongpeng/leona/internal/apiserver/service/v1"
	"github.com/dairongpeng/leona/internal/apiserver/store"
)

// SecretController create a secret handler used to handle request for secret resource.
type SecretController struct {
	srv srvv1.Service
}

// NewSecretController creates a secret handler.
func NewSecretController(store store.Factory) *SecretController {
	return &SecretController{
		srv: srvv1.NewService(store),
	}
}
