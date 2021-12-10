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

package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/dairongpeng/leona/pkg/log"
)

// UsernameKey defines the key in gin context which represents the owner of the secret.
const UsernameKey = "username"

// Context is a middleware that injects common prefix fields to gin.Context.
// Context 中间件，用来在 gin.Context 中设置 requestID和 username键
// 在打印日志时，将 gin.Context 类型的变量传递给 log.L() 函数，log.L() 函数会在日志输出中输出 requestID和 username域
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(log.KeyRequestID.String(), c.GetString(XRequestIDKey))
		c.Set(log.KeyUsername.String(), c.GetString(UsernameKey))
		c.Next()
	}
}
