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

// Code generated by "codegen -type=int /home/colin/workspace/golang/src/github.com/dairongpeng/leona/internal/pkg/code"; DO NOT EDIT.

package code

// init register error codes defines in this source code to `github.com/dairongpeng/leona/pkg/errors`
func init() {
	register(ErrUserNotFound, 404, "User not found")
	register(ErrUserAlreadyExist, 400, "User already exist")
	register(ErrReachMaxCount, 400, "Secret reach the max count")
	register(ErrSecretNotFound, 404, "Secret not found")
	register(ErrPolicyNotFound, 404, "Policy not found")
	register(ErrSuccess, 200, "OK")
	register(ErrUnknown, 500, "Internal server error")
	register(ErrBind, 400, "Error occurred while binding the request body to the struct")
	register(ErrValidation, 400, "Validation failed")
	register(ErrTokenInvalid, 401, "Token invalid")
	register(ErrPageNotFound, 404, "Page not found")
	register(ErrDatabase, 500, "Database error")
	register(ErrEncrypt, 401, "Error occurred while encrypting the user password")
	register(ErrSignatureInvalid, 401, "Signature is invalid")
	register(ErrExpired, 401, "Token expired")
	register(ErrInvalidAuthHeader, 401, "Invalid authorization header")
	register(ErrMissingHeader, 401, "The `Authorization` header was empty")
	register(ErrPasswordIncorrect, 401, "Password was incorrect")
	register(ErrPermissionDenied, 403, "Permission denied")
	register(ErrEncodingFailed, 500, "Encoding failed due to an error with the data")
	register(ErrDecodingFailed, 500, "Decoding failed due to an error with the data")
	register(ErrInvalidJSON, 500, "Data is not valid JSON")
	register(ErrEncodingJSON, 500, "JSON data could not be encoded")
	register(ErrDecodingJSON, 500, "JSON data could not be decoded")
	register(ErrInvalidYaml, 500, "Data is not valid Yaml")
	register(ErrEncodingYaml, 500, "Yaml data could not be encoded")
	register(ErrDecodingYaml, 500, "Yaml data could not be decoded")
}
