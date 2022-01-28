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

package gstash

import (
	"github.com/dairongpeng/leona/internal/gstash/config"
	"github.com/dairongpeng/leona/internal/gstash/options"
	genericapiserver "github.com/dairongpeng/leona/internal/pkg/server"
	"github.com/dairongpeng/leona/pkg/app"
	"github.com/dairongpeng/leona/pkg/log"
)

const commandDesc = `LEONA GSTASH is a pluggable analytics purger to move Analytics generated by your leona nodes to any back-end.

Find more leona-gstash information at:
    https://github.com/dairongpeng/leona/blob/master/docs/guide/en-US/cmd/leona-gstash.md`

// NewApp creates a App object with default parameters.
func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("LEONA analytics server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Flush()

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		stopCh := genericapiserver.SetupSignalHandler()

		return Run(cfg, stopCh)
	}
}
