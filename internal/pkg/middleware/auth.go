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
)

// AuthStrategy defines the set of methods used to do resource authentication.
type AuthStrategy interface {
	AuthFunc() gin.HandlerFunc
}

// AuthOperator used to switch between different authentication strategy.
type AuthOperator struct {
	strategy AuthStrategy
}

// SetStrategy used to set to another authentication strategy.
func (operator *AuthOperator) SetStrategy(strategy AuthStrategy) {
	operator.strategy = strategy
}

// AuthFunc execute resource authentication.
func (operator *AuthOperator) AuthFunc() gin.HandlerFunc {
	return operator.strategy.AuthFunc()
}
