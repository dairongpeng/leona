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

// Package gormutil is a util to convert offset and limit to default values.
package gormutil

// DefaultLimit define the default number of records to be retrieved.
const DefaultLimit = 1000

// LimitAndOffset contains offset and limit fields.
type LimitAndOffset struct {
	Offset int
	Limit  int
}

// Unpointer fill LimitAndOffset with default values if offset/limit is nil
// or it will be filled with the passed value.
func Unpointer(offset *int64, limit *int64) *LimitAndOffset {
	var o, l int = 0, DefaultLimit

	if offset != nil {
		o = int(*offset)
	}

	if limit != nil {
		l = int(*limit)
	}

	return &LimitAndOffset{
		Offset: o,
		Limit:  l,
	}
}
