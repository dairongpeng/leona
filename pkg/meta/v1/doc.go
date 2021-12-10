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

// Package v1 contains API types that are common to all versions.
//
// The package contains two categories of types:
// - external (serialized) types that lack their own version (e.g TypeMeta)
// - internal (never-serialized) types that are needed by several different
//   api groups, and so live here, to avoid duplication and/or import loops
//   (e.g. LabelSelector).
// In the future, we will probably move these categories of objects into
// separate packages.
package v1
