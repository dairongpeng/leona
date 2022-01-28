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

package analytics

// AnalyticsFilters defines the analytics options.
type AnalyticsFilters struct {
	Usernames        []string `json:"usernames"`
	SkippedUsernames []string `json:"skip_usernames"`
}

// ShouldFilter determine whether a record should to be filtered out.
func (filters AnalyticsFilters) ShouldFilter(record AnalyticsRecord) bool {
	switch {
	case len(filters.SkippedUsernames) > 0 && stringInSlice(record.Username, filters.SkippedUsernames):
		return true
	case len(filters.Usernames) > 0 && !stringInSlice(record.Username, filters.Usernames):
		return true
	}

	return false
}

// HasFilter determine whether a record has a filter.
func (filters AnalyticsFilters) HasFilter() bool {
	if len(filters.SkippedUsernames) == 0 && len(filters.Usernames) == 0 {
		return false
	}

	return true
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}
