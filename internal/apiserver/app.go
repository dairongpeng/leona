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

import (
	"github.com/dairongpeng/leona/internal/apiserver/config"
	"github.com/dairongpeng/leona/internal/apiserver/options"
	"github.com/dairongpeng/leona/pkg/app"
	"github.com/dairongpeng/leona/pkg/log"
)

const commandDesc = `The LEONA API server responsible for leona's resource
scheduling and provides a simple rest api.

Find more leona-apiserver information at:
    https://github.com/dairongpeng/leona/cmd/leona-apiserver.md`

// NewApp creates a App object with default parameters.
func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	// Build app by generic opts
	application := app.NewApp("LEONA API Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

// run creates a hook func be used to run app
func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Flush()

		// Build application configuration
		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)
	}
}
