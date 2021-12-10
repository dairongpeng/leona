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

package reflect

import (
	"reflect"
	"testing"
)

func TestGetObjFieldsMap(t *testing.T) {
	type Obj struct {
		A int
		B int
		C int
	}

	org := &Obj{
		A: 1,
		B: 2,
		C: 3,
	}

	m := GetObjFieldsMap(org, []string{})
	if !reflect.DeepEqual(m, map[string]interface{}{
		"A": 1,
		"B": 2,
		"C": 3,
	}) {
		t.Fatalf("not equal")
	}

	m = GetObjFieldsMap(org, []string{"A"})
	if !reflect.DeepEqual(m, map[string]interface{}{
		"A": 1,
	}) {
		t.Fatalf("not equal")
	}
}

func TestCopyObj(t *testing.T) {
	type Obj struct {
		A int
		B int
		C int
	}

	org := &Obj{
		A: 1,
		B: 2,
		C: 3,
	}

	des := &Obj{
		A: 4,
		B: 5,
		C: 6,
	}

	changed, err := CopyObj(org, des, []string{"A"})
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !changed {
		t.Fatalf("expect changed")
	}

	if des.A != org.A {
		t.Fatalf("A not copy")
	}

	if des.B != 5 || des.C != 6 {
		t.Fatalf("B and C changed")
	}

	des.A = org.A
	changed, err = CopyObj(org, des, []string{"A"})
	if err != nil {
		t.Fatalf(err.Error())
	}

	if changed {
		t.Fatalf("expect not changed")
	}
}
