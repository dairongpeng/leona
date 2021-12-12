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
	"github.com/dairongpeng/leona/pkg/core"
	"github.com/dairongpeng/leona/pkg/errors"
	"github.com/gin-gonic/gin"

	"github.com/dairongpeng/leona/internal/apiserver/controller/v1/user"
	"github.com/dairongpeng/leona/internal/apiserver/store/mysql"
	"github.com/dairongpeng/leona/internal/pkg/code"
	"github.com/dairongpeng/leona/internal/pkg/middleware"
	"github.com/dairongpeng/leona/internal/pkg/middleware/auth"

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

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) *gin.Engine {
	// Middlewares.
	jwtStrategy, _ := newJWTAuth().(auth.JWTStrategy)
	g.POST("/login", jwtStrategy.LoginHandler)
	g.POST("/logout", jwtStrategy.LogoutHandler)
	// Refresh time can be longer than token timeout
	g.POST("/refresh", jwtStrategy.RefreshHandler)

	auto := newAutoAuth()
	g.NoRoute(auto.AuthFunc(), func(c *gin.Context) {
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})

	// v1 handlers, requiring authentication
	storeIns, _ := mysql.GetMySQLFactoryOr(nil)
	v1 := g.Group("/v1")
	{
		// user RESTful resource
		userv1 := v1.Group("/users")
		{
			userController := user.NewUserController(storeIns)

			userv1.POST("", userController.Create)
			userv1.Use(auto.AuthFunc(), middleware.Validation())
			// v1.PUT("/find_password", userController.FindPassword)
			userv1.DELETE("", userController.DeleteCollection) // admin api
			userv1.DELETE(":name", userController.Delete)      // admin api
			userv1.PUT(":name/change-password", userController.ChangePassword)
			userv1.PUT(":name", userController.Update)
			userv1.GET("", userController.List)
			userv1.GET(":name", userController.Get) // admin api
		}

		v1.Use(auto.AuthFunc())

		// policy RESTful resource
		//policyv1 := v1.Group("/policies", middleware.Publish())
		//{
		//	policyController := policy.NewPolicyController(storeIns)
		//
		//	policyv1.POST("", policyController.Create)
		//	policyv1.DELETE("", policyController.DeleteCollection)
		//	policyv1.DELETE(":name", policyController.Delete)
		//	policyv1.PUT(":name", policyController.Update)
		//	policyv1.GET("", policyController.List)
		//	policyv1.GET(":name", policyController.Get)
		//}

		// secret RESTful resource
		//secretv1 := v1.Group("/secrets", middleware.Publish())
		//{
		//	secretController := secret.NewSecretController(storeIns)
		//
		//	secretv1.POST("", secretController.Create)
		//	secretv1.DELETE(":name", secretController.Delete)
		//	secretv1.PUT(":name", secretController.Update)
		//	secretv1.GET("", secretController.List)
		//	secretv1.GET(":name", secretController.Get)
		//}
	}

	return g
}
