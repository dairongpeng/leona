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
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/dairongpeng/leona/api/apiserver/v1"
	metav1 "github.com/dairongpeng/leona/pkg/meta/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	srvv1 "github.com/dairongpeng/leona/internal/apiserver/service/v1"
)

func TestUserController_Update(t *testing.T) {
	user := &v1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: "admin",
			ID:   0,
		},
		Nickname: "admin",
		Password: "Admin@2020",
		Email:    "admin@foxmail.com",
		Phone:    "1812884xxxx",
	}

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := bytes.NewBufferString(`{"nickname":"admin2","email":"admin2@foxmail.com","phone":"1812885xxx"}`)
	c.Request, _ = http.NewRequest("PUT", "/v1/users/admin", body)
	c.Params = []gin.Param{{Key: "name", Value: "admin"}}
	c.Request.Header.Set("Content-Type", "application/json")

	// deep copy
	user2 := new(v1.User)
	*user2 = *user
	user2.Nickname = "admin2"
	user2.Email = "admin2@foxmail.com"
	user2.Phone = "1812885xxx"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := srvv1.NewMockService(ctrl)
	mockUserSrv := srvv1.NewMockUserSrv(ctrl)
	mockUserSrv.EXPECT().Get(gomock.Any(), gomock.Eq("admin"), gomock.Any()).Return(user, nil)
	mockUserSrv.EXPECT().Update(gomock.Any(), gomock.Eq(user2), gomock.Any()).Return(nil)
	mockService.EXPECT().Users().Return(mockUserSrv).Times(2)

	type fields struct {
		srv srvv1.Service
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "default",
			fields: fields{
				srv: mockService,
			},
			args: args{
				c: c,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserController{
				srv: tt.fields.srv,
			}
			u.Update(tt.args.c)
		})
	}
}
