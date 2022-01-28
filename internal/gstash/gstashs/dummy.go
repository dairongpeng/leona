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

	"github.com/dairongpeng/leona/pkg/log"
)

// DummyGstash  defines a dummy gstash with dummy specific options and common options.
type DummyGstash struct {
	CommonGstashConfig
}

// New create a dummy gstash instance.
func (p *DummyGstash) New() Gstash {
	newGstash := DummyGstash{}

	return &newGstash
}

// GetName returns the dummy gstash name.
func (p *DummyGstash) GetName() string {
	return "Dummy Gstash"
}

// Init initialize the dummy gstash instance.
func (p *DummyGstash) Init(conf interface{}) error {
	log.Debug("Dummy Initialized")

	return nil
}

// WriteData write analyzed data to dummy persistent back-end storage.
func (p *DummyGstash) WriteData(ctx context.Context, data []interface{}) error {
	log.Infof("Writing %d records", len(data))

	return nil
}
