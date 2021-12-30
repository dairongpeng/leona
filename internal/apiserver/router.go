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

package apiserver

import (
	"github.com/gin-gonic/gin"

	"github.com/dairongpeng/leona/internal/apiserver/controller/v1/user"
	"github.com/dairongpeng/leona/internal/apiserver/store/mysql"
	// custom gin validators.
	_ "github.com/dairongpeng/leona/pkg/validator"
)

// 初始化gin中间件和路由
func initRouter(g *gin.Engine) {
	// 初始化路由中间件
	installMiddleware(g)
	// 初始化api接口
	installController(g)
}

// installMiddleware 初始化路由中间件
func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) *gin.Engine {
	// v1 handlers, requiring authentication
	storeIns, _ := mysql.GetMySQLFactoryOr(nil)
	v1 := g.Group("/v1")
	{
		// user RESTful resource
		userv1 := v1.Group("/users")
		{
			userController := user.NewUserController(storeIns)
			userv1.POST("", userController.Create)
			userv1.DELETE(":name", userController.Delete) // admin api
			userv1.PUT(":name/change-password", userController.ChangePassword)
			userv1.PUT(":name", userController.Update)
			userv1.GET("", userController.List)
			userv1.GET(":name", userController.Get) // admin api
		}
	}

	return g
}
