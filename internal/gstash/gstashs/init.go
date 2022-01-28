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

var availableGstashs map[string]Gstash

// nolint: gochecknoinits
func init() {
	availableGstashs = make(map[string]Gstash)

	// Register all the storage handlers here
	availableGstashs["csv"] = &CSVGstash{}
	availableGstashs["mongo"] = &MongoGstash{}
	availableGstashs["dummy"] = &DummyGstash{}
	availableGstashs["elasticsearch"] = &ElasticsearchGstash{}
	availableGstashs["prometheus"] = &PrometheusGstash{}
	availableGstashs["kafka"] = &KafkaGstash{}
	availableGstashs["syslog"] = &SyslogGstash{}
}
