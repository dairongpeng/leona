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

// Package validator defines leona custom binding validators used by gin.
package validator

import (
	"github.com/dairongpeng/leona/pkg/validation"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// validateUsername checks if a given username is illegal.
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if errs := validation.IsQualifiedName(username); len(errs) > 0 {
		return false
	}

	return true
}

// validatePassword checks if a given password is illegal.
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if err := validation.IsValidPassword(password); err != nil {
		return false
	}

	return true
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("username", validateUsername)
		_ = v.RegisterValidation("password", validatePassword)
	}
}
