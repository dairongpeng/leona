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

package options

import (
	"encoding/json"
	"github.com/dairongpeng/leona/internal/gstash/analytics"
	genericoptions "github.com/dairongpeng/leona/internal/pkg/options"
	cliflag "github.com/dairongpeng/leona/pkg/cli/flag"
	"github.com/dairongpeng/leona/pkg/log"
)

// GStashConfig defines options for gstash back-end.
type GStashConfig struct {
	Type                  string                     `json:"type"                    mapstructure:"type"`
	Filters               analytics.AnalyticsFilters `json:"filters"                 mapstructure:"filters"`
	Timeout               int                        `json:"timeout"                 mapstructure:"timeout"`
	OmitDetailedRecording bool                       `json:"omit-detailed-recording" mapstructure:"omit-detailed-recording"`
	Meta                  map[string]interface{}     `json:"meta"                    mapstructure:"meta"`
}

// Options runs a gstashserver.
type Options struct {
	PurgeDelay            int                          `json:"purge-delay"             mapstructure:"purge-delay"`
	Gstashs               map[string]GStashConfig      `json:"gstashs"                   mapstructure:"gstashs"`
	HealthCheckPath       string                       `json:"health-check-path"       mapstructure:"health-check-path"`
	HealthCheckAddress    string                       `json:"health-check-address"    mapstructure:"health-check-address"`
	OmitDetailedRecording bool                         `json:"omit-detailed-recording" mapstructure:"omit-detailed-recording"`
	RedisOptions          *genericoptions.RedisOptions `json:"redis"                   mapstructure:"redis"`
	Log                   *log.Options                 `json:"log"                     mapstructure:"log"`
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	s := Options{
		PurgeDelay: 10,
		Gstashs: map[string]GStashConfig{
			"csv": {
				Type: "csv",
				Meta: map[string]interface{}{
					"csv_dir": "./analytics-data",
				},
			},
		},
		HealthCheckPath:    "healthz",
		HealthCheckAddress: "0.0.0.0:7070",
		RedisOptions:       genericoptions.NewRedisOptions(),
		Log:                log.NewOptions(),
	}

	return &s
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.Log.AddFlags(fss.FlagSet("logs"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	fs.IntVar(&o.PurgeDelay, "purge-delay", o.PurgeDelay, ""+
		"This setting the purge delay (in seconds) when purge the data from Redis to MongoDB or other data stores.")
	fs.StringVar(&o.HealthCheckPath, "health-check-path", o.HealthCheckPath, ""+
		"Specifies liveness health check request path.")
	fs.StringVar(&o.HealthCheckAddress, "health-check-address", o.HealthCheckAddress, ""+
		"Specifies liveness health check bind address.")
	fs.BoolVar(&o.OmitDetailedRecording, "omit-detailed-recording", o.OmitDetailedRecording, ""+
		"Setting this to true will avoid writing policy fields for each authorization request in gstashs.")

	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
