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

package gstashs

import (
	"context"
	"errors"

	"github.com/dairongpeng/leona/internal/gstash/analytics"
)

// Gstash defines the interface for all analytics back-end.
type Gstash interface {
	GetName() string
	New() Gstash
	Init(interface{}) error
	WriteData(context.Context, []interface{}) error
	SetFilters(analytics.AnalyticsFilters)
	GetFilters() analytics.AnalyticsFilters
	SetTimeout(timeout int)
	GetTimeout() int
	SetOmitDetailedRecording(bool)
	GetOmitDetailedRecording() bool
}

// GetGstashByName returns the gstash instance by given name.
func GetGstashByName(name string) (Gstash, error) {
	if gstash, ok := availableGstashs[name]; ok && gstash != nil {
		return gstash, nil
	}

	return nil, errors.New(name + " Not found")
}
