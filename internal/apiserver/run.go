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

package apiserver

import "github.com/dairongpeng/leona/internal/apiserver/config"

// Run runs the specified APIServer. This should never exit.
func Run(cfg *config.Config) error {
	// the configuration used to create the HTTP/GRPC server according to the application configuration
	server, err := createAPIServer(cfg)
	if err != nil {
		return err
	}

	// PrepareRun for initializing the HTTP/GRPC server before it starts
	// Afterwards, the Run method is called to start the GRPC and HTTP server
	return server.PrepareRun().Run()
}
