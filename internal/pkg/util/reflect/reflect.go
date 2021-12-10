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
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func ToGormDBMap(obj interface{}, fields []string) (map[string]interface{}, error) {
	reflectType := reflect.ValueOf(obj).Type()
	reflectValue := reflect.ValueOf(obj)
	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
		reflectValue = reflect.ValueOf(obj).Elem()
	}

	ret := make(map[string]interface{})
	for _, f := range fields {
		fs, exist := reflectType.FieldByName(f)
		if !exist {
			return nil, fmt.Errorf("unknow field " + f)
		}

		tagMap := parseTagSetting(fs.Tag)
		gormfiled, exist := tagMap["COLUMN"]
		if !exist {
			return nil, fmt.Errorf("undef gorm field " + f)
		}

		ret[gormfiled] = reflectValue.FieldByName(f)
	}
	return ret, nil
}

func parseTagSetting(tags reflect.StructTag) map[string]string {
	setting := map[string]string{}
	for _, str := range []string{tags.Get("sql"), tags.Get("gorm")} {
		if str == "" {
			continue
		}
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				setting[k] = strings.Join(v[1:], ":")
			} else {
				setting[k] = k
			}
		}
	}
	return setting
}

func GetObjFieldsMap(obj interface{}, fields []string) map[string]interface{} {
	ret := make(map[string]interface{})

	modelReflect := reflect.ValueOf(obj)
	if modelReflect.Kind() == reflect.Ptr {
		modelReflect = modelReflect.Elem()
	}

	modelRefType := modelReflect.Type()
	fieldsCount := modelReflect.NumField()
	var fieldData interface{}
	for i := 0; i < fieldsCount; i++ {
		field := modelReflect.Field(i)
		if len(fields) != 0 && !findString(fields, modelRefType.Field(i).Name) {
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			fallthrough
		case reflect.Ptr:
			fieldData = GetObjFieldsMap(field.Interface(), []string{})
		default:
			fieldData = field.Interface()
		}

		ret[modelRefType.Field(i).Name] = fieldData
	}

	return ret
}

func CopyObj(from interface{}, to interface{}, fields []string) (changed bool, err error) {
	fromMap := GetObjFieldsMap(from, fields)
	toMap := GetObjFieldsMap(to, fields)
	if reflect.DeepEqual(fromMap, toMap) {
		return false, nil
	}

	t := reflect.ValueOf(to).Elem()
	for k, v := range fromMap {
		val := t.FieldByName(k)
		val.Set(reflect.ValueOf(v))
	}
	return true, nil
}

// CopyObjViaYaml marshal "from" to yaml data, then unMarshal data to "to".
func CopyObjViaYaml(to interface{}, from interface{}) error {
	if from == nil || to == nil {
		return nil
	}

	data, err := yaml.Marshal(from)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, to)
}

// findString return true if target in slice, return false if not.
func findString(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}
