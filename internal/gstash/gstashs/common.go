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

import "github.com/dairongpeng/leona/internal/gstash/analytics"

// CommonGstashConfig defines common options used by all persistent store, like elasticsearch, kafka, mongo and etc.
type CommonGstashConfig struct {
	filters               analytics.AnalyticsFilters
	timeout               int
	OmitDetailedRecording bool
}

// SetFilters set attributes `filters` for CommonGstashConfig.
func (p *CommonGstashConfig) SetFilters(filters analytics.AnalyticsFilters) {
	p.filters = filters
}

// GetFilters get attributes `filters` for CommonGstashConfig.
func (p *CommonGstashConfig) GetFilters() analytics.AnalyticsFilters {
	return p.filters
}

// SetTimeout set attributes `timeout` for CommonGstashConfig.
func (p *CommonGstashConfig) SetTimeout(timeout int) {
	p.timeout = timeout
}

// GetTimeout get attributes `timeout` for CommonGstashConfig.
func (p *CommonGstashConfig) GetTimeout() int {
	return p.timeout
}

// SetOmitDetailedRecording set attributes `OmitDetailedRecording` for CommonGstashConfig.
func (p *CommonGstashConfig) SetOmitDetailedRecording(omitDetailedRecording bool) {
	p.OmitDetailedRecording = omitDetailedRecording
}

// GetOmitDetailedRecording get attributes `OmitDetailedRecording` for CommonGstashConfig.
func (p *CommonGstashConfig) GetOmitDetailedRecording() bool {
	return p.OmitDetailedRecording
}
