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

// nolint: lll
package logger_test

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"encoding/json"
	"github.com/jinzhu/now"
	"gorm.io/gorm/logger"
)

type JSON json.RawMessage

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

type ExampleStruct struct {
	Name string
	Val  string
}

func (s ExampleStruct) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// nolint: unparam
func format(v []byte, escaper string) string {
	return escaper + strings.ReplaceAll(string(v), escaper, "\\"+escaper) + escaper
}

func TestExplainSQL(t *testing.T) {
	type role string
	type password []byte
	var (
		tt     = now.MustParse("2020-02-23 11:10:10")
		myrole = role("admin")
		pwd    = password([]byte("pass"))
		jsVal  = []byte(`{"Name":"test","Val":"test"}`)
		js     = JSON(jsVal)
		esVal  = []byte(`{"Name":"test","Val":"test"}`)
		es     = ExampleStruct{Name: "test", Val: "test"}
	)

	results := []struct {
		SQL           string
		NumericRegexp *regexp.Regexp
		Vars          []interface{}
		Result        string
	}{
		{
			SQL:           "create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			NumericRegexp: nil,
			Vars: []interface{}{
				"jinzhu",
				1,
				999.99,
				true,
				[]byte("12345"),
				tt,
				&tt,
				nil,
				"w@g.\"com",
				myrole,
				pwd,
			},
			Result: `create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass) values ("jinzhu", 1, 999.990000, true, "12345", "2020-02-23 11:10:10", "2020-02-23 11:10:10", NULL, "w@g.\"com", "admin", "pass")`,
		},
		{
			SQL:           "create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			NumericRegexp: nil,
			Vars: []interface{}{
				"jinzhu?",
				1,
				999.99,
				true,
				[]byte("12345"),
				tt,
				&tt,
				nil,
				"w@g.\"com",
				myrole,
				pwd,
			},
			Result: `create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass) values ("jinzhu?", 1, 999.990000, true, "12345", "2020-02-23 11:10:10", "2020-02-23 11:10:10", NULL, "w@g.\"com", "admin", "pass")`,
		},
		{
			SQL:           "create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass) values (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11)",
			NumericRegexp: regexp.MustCompile(`@p(\d+)`),
			Vars: []interface{}{
				"jinzhu",
				1,
				999.99,
				true,
				[]byte("12345"),
				tt,
				&tt,
				nil,
				"w@g.com",
				myrole,
				pwd,
			},
			Result: `create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass) values ("jinzhu", 1, 999.990000, true, "12345", "2020-02-23 11:10:10", "2020-02-23 11:10:10", NULL, "w@g.com", "admin", "pass")`,
		},
		{
			SQL:           "create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass) values ($3, $4, $1, $2, $7, $8, $5, $6, $9, $10, $11)",
			NumericRegexp: regexp.MustCompile(`\$(\d+)`),
			Vars: []interface{}{
				999.99,
				true,
				"jinzhu",
				1,
				&tt,
				nil,
				[]byte("12345"),
				tt,
				"w@g.com",
				myrole,
				pwd,
			},
			Result: `create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass) values ("jinzhu", 1, 999.990000, true, "12345", "2020-02-23 11:10:10", "2020-02-23 11:10:10", NULL, "w@g.com", "admin", "pass")`,
		},
		{
			SQL:           "create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass) values (@p1, @p11, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10)",
			NumericRegexp: regexp.MustCompile(`@p(\d+)`),
			Vars: []interface{}{
				"jinzhu",
				999.99,
				true,
				[]byte("12345"),
				tt,
				&tt,
				nil,
				"w@g.com",
				myrole,
				pwd,
				1,
			},
			Result: `create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass) values ("jinzhu", 1, 999.990000, true, "12345", "2020-02-23 11:10:10", "2020-02-23 11:10:10", NULL, "w@g.com", "admin", "pass")`,
		},
		{
			SQL:           "create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass, json_struct, example_struct) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			NumericRegexp: nil,
			Vars: []interface{}{
				"jinzhu",
				1,
				999.99,
				true,
				[]byte("12345"),
				tt,
				&tt,
				nil,
				"w@g.\"com",
				myrole,
				pwd,
				js,
				es,
			},
			Result: fmt.Sprintf(
				`create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass, json_struct, example_struct) values ("jinzhu", 1, 999.990000, true, "12345", "2020-02-23 11:10:10", "2020-02-23 11:10:10", NULL, "w@g.\"com", "admin", "pass", %v, %v)`,
				format(jsVal, `"`),
				format(esVal, `"`),
			),
		},
		{
			SQL:           "create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass, json_struct, example_struct) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			NumericRegexp: nil,
			Vars: []interface{}{
				"jinzhu",
				1,
				999.99,
				true,
				[]byte("12345"),
				tt,
				&tt,
				nil,
				"w@g.\"com",
				myrole,
				pwd,
				&js,
				&es,
			},
			Result: fmt.Sprintf(
				`create table users (name, age, height, actived, bytes, create_at, update_at, deleted_at, email, role, pass, json_struct, example_struct) values ("jinzhu", 1, 999.990000, true, "12345", "2020-02-23 11:10:10", "2020-02-23 11:10:10", NULL, "w@g.\"com", "admin", "pass", %v, %v)`,
				format(jsVal, `"`),
				format(esVal, `"`),
			),
		},
	}

	for idx, r := range results {
		if result := logger.ExplainSQL(r.SQL, r.NumericRegexp, `"`, r.Vars...); result != r.Result {
			t.Errorf("Explain SQL #%v expects %v, but got %v", idx, r.Result, result)
		}
	}
}
