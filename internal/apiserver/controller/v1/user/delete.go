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
	"github.com/dairongpeng/leona/internal/apiserver/analytics"
	"github.com/dairongpeng/leona/pkg/core"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/dairongpeng/leona/pkg/log"
)

// Delete delete an user by the user identifier.
// Only administrator can call this function.
func (u *UserController) Delete(c *gin.Context) {
	log.L(c).Info("delete user function called.")

	if err := u.srv.Users().Delete(c, c.Param("name"), metav1.DeleteOptions{Unscoped: true}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	// 收集数据
	record := analytics.AnalyticsRecord{
		TimeStamp: time.Now().Unix(),
		Username:  "-",
		Effect:    "delete",
	}
	record.SetExpiry(0)
	_ = analytics.GetAnalytics().RecordHit(&record)

	core.WriteResponse(c, nil, nil)
}
