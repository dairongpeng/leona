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
	"fmt"
	"log/syslog"

	"github.com/mitchellh/mapstructure"

	"github.com/dairongpeng/leona/internal/gstash/analytics"
	"github.com/dairongpeng/leona/pkg/log"
)

// SyslogGstash defines a syslog gstash with syslog specific options and common options.
type SyslogGstash struct {
	syslogConf *SyslogConf
	writer     *syslog.Writer
	filters    analytics.AnalyticsFilters
	timeout    int
	CommonGstashConfig
}

var logPrefix = "syslog-gstash"

// SyslogConf defines syslog specific options.
type SyslogConf struct {
	Transport   string `mapstructure:"transport"`
	NetworkAddr string `mapstructure:"network_addr"`
	LogLevel    int    `mapstructure:"log_level"`
	Tag         string `mapstructure:"tag"`
}

// New create a syslog gstash instance.
func (s *SyslogGstash) New() Gstash {
	newStash := SyslogGstash{}

	return &newStash
}

// GetName returns the syslog gstash name.
func (s *SyslogGstash) GetName() string {
	return "Syslog Gstash"
}

// Init initialize the syslog gstash instance.
func (s *SyslogGstash) Init(config interface{}) error {
	// Read configuration file
	s.syslogConf = &SyslogConf{}
	err := mapstructure.Decode(config, &s.syslogConf)
	if err != nil {
		log.Fatalf("Failed to decode configuration: %s", err.Error())
	}

	// Init the configs
	initConfigs(s)

	// Init the Syslog writer
	initWriter(s)

	log.Debug("Syslog Gstash active")

	return nil
}

func initWriter(s *SyslogGstash) {
	tag := logPrefix
	if s.syslogConf.Tag != "" {
		tag = s.syslogConf.Tag
	}
	syslogWriter, err := syslog.Dial(
		s.syslogConf.Transport,
		s.syslogConf.NetworkAddr,
		syslog.Priority(s.syslogConf.LogLevel),
		tag)
	if err != nil {
		log.Fatalf("failed to connect to Syslog Daemon: %s", err.Error())
	}

	s.writer = syslogWriter
}

// Set default values if they are not explicitly given and perform validation.
func initConfigs(gstash *SyslogGstash) {
	if gstash.syslogConf.Transport == "" {
		gstash.syslogConf.Transport = "udp"
		log.Info("No Transport given, using 'udp'")
	}

	if gstash.syslogConf.Transport != "udp" &&
		gstash.syslogConf.Transport != "tcp" &&
		gstash.syslogConf.Transport != "tls" {
		log.Fatal("Chosen invalid Transport type.  Please use a supported Transport type for Syslog")
	}

	if gstash.syslogConf.NetworkAddr == "" {
		gstash.syslogConf.NetworkAddr = "localhost:5140"
		log.Info("No host given, using 'localhost:5140'")
	}

	if gstash.syslogConf.LogLevel == 0 {
		log.Warn("Using Log Level 0 (KERNEL) for Syslog gstash")
	}
}

// WriteData write analyzed data to syslog persistent back-end storage.
func (s *SyslogGstash) WriteData(ctx context.Context, data []interface{}) error {
	// Data is all the analytics being written
	for _, v := range data {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Decode the raw analytics into Form
			decoded, _ := v.(analytics.AnalyticsRecord)
			message := Message{
				"timestamp":  decoded.TimeStamp,
				"username":   decoded.Username,
				"effect":     decoded.Effect,
				"conclusion": decoded.Conclusion,
				"request":    decoded.Request,
				"policies":   decoded.Policies,
				"deciders":   decoded.Deciders,
				"expireAt":   decoded.ExpireAt,
			}

			// Print to Syslog
			_, _ = fmt.Fprintf(s.writer, "%s", message)
		}
	}

	return nil
}

// SetTimeout set attributes `timeout` for SyslogGstash.
func (s *SyslogGstash) SetTimeout(timeout int) {
	s.timeout = timeout
}

// GetTimeout get attributes `timeout` for SyslogGstash.
func (s *SyslogGstash) GetTimeout() int {
	return s.timeout
}

// SetFilters set attributes `filters` for SyslogGstash.
func (s *SyslogGstash) SetFilters(filters analytics.AnalyticsFilters) {
	s.filters = filters
}

// GetFilters get attributes `filters` for SyslogGstash.
func (s *SyslogGstash) GetFilters() analytics.AnalyticsFilters {
	return s.filters
}
