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

package runtime

import (
	"fmt"

	"github.com/dairongpeng/leona/pkg/json"
)

// NegotiateError is returned when a ClientNegotiator is unable to locate
// a serializer for the requested operation.
type NegotiateError struct {
	ContentType string
	Stream      bool
}

func (e NegotiateError) Error() string {
	if e.Stream {
		return fmt.Sprintf("no stream serializers registered for %s", e.ContentType)
	}
	return fmt.Sprintf("no serializers registered for %s", e.ContentType)
}

type apimachineryClientNegotiator struct{}

var _ ClientNegotiator = &apimachineryClientNegotiator{}

func (n *apimachineryClientNegotiator) Encoder() (Encoder, error) {
	return &apimachineryClientNegotiatorSerializer{}, nil
}

func (n *apimachineryClientNegotiator) Decoder() (Decoder, error) {
	return &apimachineryClientNegotiatorSerializer{}, nil
}

type apimachineryClientNegotiatorSerializer struct{}

var _ Decoder = &apimachineryClientNegotiatorSerializer{}

func (s *apimachineryClientNegotiatorSerializer) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (s *apimachineryClientNegotiatorSerializer) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// NewSimpleClientNegotiator will negotiate for a single serializer. This should only be used
// for testing or when the caller is taking responsibility for setting the GVK on encoded objects.
func NewSimpleClientNegotiator() ClientNegotiator {
	return &apimachineryClientNegotiator{}
}
