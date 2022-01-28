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
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/dairongpeng/leona/pkg/errors"
	"github.com/mitchellh/mapstructure"

	"github.com/dairongpeng/leona/internal/gstash/analytics"
	"github.com/dairongpeng/leona/pkg/log"
)

// CSVGstash defines a csv gstash with csv specific options and common options.
type CSVGstash struct {
	csvConf *CSVConf
	CommonGstashConfig
}

// CSVConf defines csv specific options.
type CSVConf struct {
	// Specify the directory used to store automatically generated csv file which contains analyzed data.
	CSVDir string `mapstructure:"csv_dir"`
}

// New create a csv gstash instance.
func (c *CSVGstash) New() Gstash {
	newGstash := CSVGstash{}

	return &newGstash
}

// GetName returns the csv gstash name.
func (c *CSVGstash) GetName() string {
	return "CSV Gstash"
}

// Init initialize the csv gstash instance.
func (c *CSVGstash) Init(conf interface{}) error {
	c.csvConf = &CSVConf{}
	err := mapstructure.Decode(conf, &c.csvConf)
	if err != nil {
		log.Fatalf("Failed to decode configuration: %s", err.Error())
	}

	ferr := os.MkdirAll(c.csvConf.CSVDir, 0o777)
	if ferr != nil {
		log.Error(ferr.Error())
	}

	log.Debug("CSV Initialized")

	return nil
}

// WriteData write analyzed data to csv persistent back-end storage.
func (c *CSVGstash) WriteData(ctx context.Context, data []interface{}) error {
	curtime := time.Now()
	fname := fmt.Sprintf("%d-%s-%d-%d.csv", curtime.Year(), curtime.Month().String(), curtime.Day(), curtime.Hour())
	fname = path.Join(c.csvConf.CSVDir, fname)

	var outfile *os.File
	var appendHeader bool

	if _, err := os.Stat(fname); os.IsNotExist(err) {
		var createErr error
		outfile, createErr = os.Create(fname)
		if createErr != nil {
			log.Errorf("Failed to create new CSV file: %s", createErr.Error())
		}
		appendHeader = true
	} else {
		var appendErr error
		outfile, appendErr = os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, 0o600)
		if appendErr != nil {
			log.Errorf("Failed to open CSV file: %s", appendErr.Error())
		}
	}

	defer outfile.Close()
	writer := csv.NewWriter(outfile)

	if appendHeader {
		startRecord := analytics.AnalyticsRecord{}
		headers := startRecord.GetFieldNames()

		err := writer.Write(headers)
		if err != nil {
			log.Errorf("Failed to write file headers: %s", err.Error())

			return errors.Wrap(err, "failed to write file headers")
		}
	}

	for _, v := range data {
		decoded, _ := v.(analytics.AnalyticsRecord)

		toWrite := decoded.GetLineValues()
		err := writer.Write(toWrite)
		if err != nil {
			log.Error("File write failed!")
			log.Error(err.Error())
		}
	}

	writer.Flush()

	return nil
}
