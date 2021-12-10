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
	"github.com/dairongpeng/leona/pkg/core"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"
	"github.com/gin-gonic/gin"

	"github.com/dairongpeng/leona/internal/pkg/middleware"
	"github.com/dairongpeng/leona/pkg/log"
)

// DeleteCollection delete secrets by secret names.
func (s *SecretController) DeleteCollection(c *gin.Context) {
	log.L(c).Info("batch delete policy function called.")

	if err := s.srv.Policies().DeleteCollection(
		c,
		c.GetString(middleware.UsernameKey),
		c.QueryArray("name"),
		metav1.DeleteOptions{},
	); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
