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

package gstashs

import (
	"context"
	"errors"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/ory/ladon"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/dairongpeng/leona/internal/gstash/analytics"
	"github.com/dairongpeng/leona/pkg/log"
)

// PrometheusGstash defines a prometheus gstash with prometheus specific options and common options.
type PrometheusGstash struct {
	conf *PrometheusConf
	// Per service
	TotalStatusMetrics *prometheus.CounterVec

	CommonGstashConfig
}

// PrometheusConf defines prometheus specific options.
type PrometheusConf struct {
	Addr string `mapstructure:"listen_address"`
	Path string `mapstructure:"path"`
}

// New create a prometheus gstash instance.
func (p *PrometheusGstash) New() Gstash {
	newGStash := PrometheusGstash{}
	newGStash.TotalStatusMetrics = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "iam_user_authorization_status_total",
			Help: "authorization effect per user",
		},
		[]string{"code", "username"},
	)

	prometheus.MustRegister(newGStash.TotalStatusMetrics)

	return &newGStash
}

// GetName returns the prometheus gstash name.
func (p *PrometheusGstash) GetName() string {
	return "Prometheus GStash"
}

// Init initialize the prometheus gstash instance.
func (p *PrometheusGstash) Init(conf interface{}) error {
	p.conf = &PrometheusConf{}
	err := mapstructure.Decode(conf, &p.conf)
	if err != nil {
		log.Fatalf("Failed to decode configuration: %s", err.Error())
	}

	if p.conf.Path == "" {
		p.conf.Path = "/metrics"
	}

	if p.conf.Addr == "" {
		return errors.New("prometheus listen_addr not set")
	}

	log.Infof("Starting prometheus listener on: %s", p.conf.Addr)

	http.Handle(p.conf.Path, promhttp.Handler())

	go func() {
		log.Fatal(http.ListenAndServe(p.conf.Addr, nil).Error())
	}()

	return nil
}

// WriteData write analyzed data to prometheus persistent back-end storage.
func (p *PrometheusGstash) WriteData(ctx context.Context, data []interface{}) error {
	log.Debugf("Writing %d records", len(data))

	for _, item := range data {
		record, _ := item.(analytics.AnalyticsRecord)
		code := "0"
		if record.Effect != ladon.AllowAccess {
			code = "1"
		}

		p.TotalStatusMetrics.WithLabelValues(code, record.Username).Inc()
	}

	return nil
}
