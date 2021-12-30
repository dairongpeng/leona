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

package user

import (
	v1 "github.com/dairongpeng/leona/api/apiserver/v1"
	"github.com/dairongpeng/leona/pkg/auth"
	"github.com/dairongpeng/leona/pkg/core"
	"github.com/dairongpeng/leona/pkg/errors"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/dairongpeng/leona/internal/pkg/code"
	"github.com/dairongpeng/leona/pkg/log"
)

// Create add new user to the storage.
func (u *UserController) Create(c *gin.Context) {
	log.L(c).Info("user create function called.")

	var r v1.User

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	if errs := r.Validate(); len(errs) != 0 {
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), nil)

		return
	}

	r.Password, _ = auth.Encrypt(r.Password)
	r.Status = 1
	r.LoginedAt = time.Now()

	// Insert the user to the storage.
	if err := u.srv.Users().Create(c, &r, metav1.CreateOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, r)
}
