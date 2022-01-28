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

import "testing"

func TestShouldFilter(t *testing.T) {
	record := AnalyticsRecord{
		Username: "colin",
	}

	// test skip_usernames
	filter := AnalyticsFilters{
		SkippedUsernames: []string{"colin"},
	}
	shouldFilter := filter.ShouldFilter(record)
	if shouldFilter == false {
		t.Fatal("filter should be filtering the record")
	}

	// test different usernames
	filter = AnalyticsFilters{
		Usernames: []string{"james"},
	}
	shouldFilter = filter.ShouldFilter(record)
	if shouldFilter == false {
		t.Fatal("filter should be filtering the record")
	}

	// test no filter
	filter = AnalyticsFilters{}
	shouldFilter = filter.ShouldFilter(record)
	if shouldFilter == true {
		t.Fatal("filter should not be filtering the record")
	}
}

func TestHasFilter(t *testing.T) {
	filter := AnalyticsFilters{}

	hasFilter := filter.HasFilter()
	if hasFilter == true {
		t.Fatal("Has filter should be false.")
	}

	filter = AnalyticsFilters{
		Usernames: []string{"colin"},
	}
	hasFilter = filter.HasFilter()
	if hasFilter == false {
		t.Fatal("HasFilter should be true.")
	}
}
