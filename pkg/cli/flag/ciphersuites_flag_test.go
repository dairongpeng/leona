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

package flag

import (
	"crypto/tls"
	"fmt"
	"go/importer"
	"reflect"
	"strings"
	"testing"
)

func TestStrToUInt16(t *testing.T) {
	tests := []struct {
		flag          []string
		expected      []uint16
		expectedError bool
	}{
		{
			// Happy case
			flag: []string{
				"TLS_RSA_WITH_RC4_128_SHA",
				"TLS_RSA_WITH_AES_128_CBC_SHA",
				"TLS_ECDHE_RSA_WITH_RC4_128_SHA",
				"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
			},
			expected: []uint16{
				tls.TLS_RSA_WITH_RC4_128_SHA,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			},
			expectedError: false,
		},
		{
			// One flag only
			flag:          []string{"TLS_RSA_WITH_RC4_128_SHA"},
			expected:      []uint16{tls.TLS_RSA_WITH_RC4_128_SHA},
			expectedError: false,
		},
		{
			// Empty flag
			flag:          []string{},
			expected:      nil,
			expectedError: false,
		},
		{
			// Duplicated flag
			flag: []string{
				"TLS_RSA_WITH_RC4_128_SHA",
				"TLS_RSA_WITH_AES_128_CBC_SHA",
				"TLS_ECDHE_RSA_WITH_RC4_128_SHA",
				"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
				"TLS_RSA_WITH_RC4_128_SHA",
			},
			expected: []uint16{
				tls.TLS_RSA_WITH_RC4_128_SHA,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_RC4_128_SHA,
			},
			expectedError: false,
		},
		{
			// Invalid flag
			flag:          []string{"foo"},
			expected:      nil,
			expectedError: true,
		},
	}

	for i, test := range tests {
		uIntFlags, err := TLSCipherSuites(test.flag)
		if reflect.DeepEqual(uIntFlags, test.expected) == false {
			t.Errorf("%d: expected %+v, got %+v", i, test.expected, uIntFlags)
		}
		if test.expectedError && err == nil {
			t.Errorf("%d: expecting error, got %+v", i, err)
		}
	}
}

func TestConstantMaps(t *testing.T) {
	pkg, err := importer.Default().Import("crypto/tls")
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	discoveredVersions := map[string]bool{}
	discoveredCiphers := map[string]bool{}
	for _, declName := range pkg.Scope().Names() {
		if strings.HasPrefix(declName, "VersionTLS") {
			discoveredVersions[declName] = true
		}
		if strings.HasPrefix(declName, "TLS_") && !strings.HasPrefix(declName, "TLS_FALLBACK_") {
			discoveredCiphers[declName] = true
		}
	}

	acceptedCiphers := allCiphers()
	for k := range discoveredCiphers {
		if _, ok := acceptedCiphers[k]; !ok {
			t.Errorf("discovered cipher tls.%s not in ciphers map", k)
		}
	}
	for k := range acceptedCiphers {
		if _, ok := discoveredCiphers[k]; !ok {
			t.Errorf("ciphers map has %s not in tls package", k)
		}
	}
	for k := range discoveredVersions {
		if _, ok := versions[k]; !ok {
			t.Errorf("discovered version tls.%s not in version map", k)
		}
	}
	for k := range versions {
		if _, ok := discoveredVersions[k]; !ok {
			t.Errorf("versions map has %s not in tls package", k)
		}
	}
}
