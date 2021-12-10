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

package errors

/*
WARNING - changing the line numbers in this file will break the
examples.
*/

import (
	"fmt"
)

const (
	// Error codes below 1000 are reserved future use by the
	// "github.com/bdlm/errors" package.
	ConfigurationNotValid int = iota + 1000
	ErrInvalidJSON
	ErrEOF
	ErrLoadConfigFailed
)

func init() {
	Register(defaultCoder{ConfigurationNotValid, 500, "ConfigurationNotValid error", ""})
	Register(defaultCoder{ErrInvalidJSON, 500, "Data is not valid JSON", ""})
	Register(defaultCoder{ErrEOF, 500, "End of input", ""})
	Register(defaultCoder{ErrLoadConfigFailed, 500, "Load configuration file failed", ""})
}

func loadConfig() error {
	err := decodeConfig()
	return WrapC(err, ConfigurationNotValid, "service configuration could not be loaded")
}

func decodeConfig() error {
	err := readConfig()
	return WrapC(err, ErrInvalidJSON, "could not decode configuration data")
}

func readConfig() error {
	err := fmt.Errorf("read: end of input")
	return WrapC(err, ErrEOF, "could not read configuration file")
}
