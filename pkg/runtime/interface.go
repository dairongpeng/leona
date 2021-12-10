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

// Package runtime defines some functions used to encode/decode object.
package runtime

// Encoder writes objects to a serialized form.
type Encoder interface {
	// Encode writes an object to a stream. Implementations may return errors if the versions are
	// incompatible, or if no conversion is defined.
	Encode(v interface{}) ([]byte, error)
}

// Decoder attempts to load an object from data.
type Decoder interface {
	Decode(data []byte, v interface{}) error
}

// ClientNegotiator handles turning an HTTP content type into the appropriate encoder.
// Use NewClientNegotiator or NewVersionedClientNegotiator to create this interface from
// a NegotiatedSerializer.
type ClientNegotiator interface {
	Encoder() (Encoder, error)
	Decoder() (Decoder, error)
}
