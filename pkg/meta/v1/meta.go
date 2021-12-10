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

package v1

import (
	"time"

	"github.com/dairongpeng/leona/pkg/scheme"
)

type ObjectMetaAccessor interface {
	GetObjectMeta() Object
}

// Object lets you work with object metadata from any of the versioned or
// internal API objects. Attempting to set or retrieve a field on an object that does
// not support that field (Name, UID, Namespace on lists) will be a no-op and return
// a default value.
type Object interface {
	GetID() uint64
	SetID(id uint64)
	GetName() string
	SetName(name string)
	GetCreatedAt() time.Time
	SetCreatedAt(createdAt time.Time)
	GetUpdatedAt() time.Time
	SetUpdatedAt(updatedAt time.Time)
}

// ListInterface lets you work with list metadata from any of the versioned or
// internal API objects. Attempting to set or retrieve a field on an object that does
// not support that field will be a no-op and return a default value.
type ListInterface interface {
	GetTotalCount() int64
	SetTotalCount(count int64)
}

// Type exposes the type and APIVersion of versioned or internal API objects.
type Type interface {
	GetAPIVersion() string
	SetAPIVersion(version string)
	GetKind() string
	SetKind(kind string)
}

var _ ListInterface = &ListMeta{}

func (meta *ListMeta) GetTotalCount() int64      { return meta.TotalCount }
func (meta *ListMeta) SetTotalCount(count int64) { meta.TotalCount = count }

var _ Type = &TypeMeta{}

func (obj *TypeMeta) GetObjectKind() scheme.ObjectKind { return obj }

// SetGroupVersionKind satisfies the ObjectKind interface for all objects that embed TypeMeta.
func (obj *TypeMeta) SetGroupVersionKind(gvk scheme.GroupVersionKind) {
	obj.APIVersion, obj.Kind = gvk.ToAPIVersionAndKind()
}

// GroupVersionKind satisfies the ObjectKind interface for all objects that embed TypeMeta.
func (obj *TypeMeta) GroupVersionKind() scheme.GroupVersionKind {
	return scheme.FromAPIVersionAndKind(obj.APIVersion, obj.Kind)
}

func (meta *TypeMeta) GetAPIVersion() string        { return meta.APIVersion }
func (meta *TypeMeta) SetAPIVersion(version string) { meta.APIVersion = version }
func (meta *TypeMeta) GetKind() string              { return meta.Kind }
func (meta *TypeMeta) SetKind(kind string)          { meta.Kind = kind }

func (obj *ListMeta) GetListMeta() ListInterface { return obj }

func (obj *ObjectMeta) GetObjectMeta() Object { return obj }

var _ Object = &ObjectMeta{}

func (meta *ObjectMeta) GetID() uint64                    { return meta.ID }
func (meta *ObjectMeta) SetID(id uint64)                  { meta.ID = id }
func (meta *ObjectMeta) GetName() string                  { return meta.Name }
func (meta *ObjectMeta) SetName(name string)              { meta.Name = name }
func (meta *ObjectMeta) GetCreatedAt() time.Time          { return meta.CreatedAt }
func (meta *ObjectMeta) SetCreatedAt(createdAt time.Time) { meta.CreatedAt = createdAt }
func (meta *ObjectMeta) GetUpdatedAt() time.Time          { return meta.UpdatedAt }
func (meta *ObjectMeta) SetUpdatedAt(updatedAt time.Time) { meta.UpdatedAt = updatedAt }
