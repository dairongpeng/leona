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
	"encoding/base64"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dairongpeng/leona/pkg/errors"
	"github.com/mitchellh/mapstructure"
	elastic "github.com/olivere/elastic/v7"

	"github.com/dairongpeng/leona/internal/gstash/analytics"
	"github.com/dairongpeng/leona/pkg/log"
)

// ElasticsearchGstash defines a elasticsearch gstash with elasticsearch specific options and common options.
type ElasticsearchGstash struct {
	operator ElasticsearchOperator
	esConf   *ElasticsearchConf
	CommonGstashConfig
}

// ElasticsearchConf defines elasticsearch specific options.
type ElasticsearchConf struct {
	BulkConfig       ElasticsearchBulkConfig `mapstructure:"bulk_config"`
	IndexName        string                  `mapstructure:"index_name"`
	ElasticsearchURL string                  `mapstructure:"elasticsearch_url"`
	DocumentType     string                  `mapstructure:"document_type"`
	AuthAPIKeyID     string                  `mapstructure:"auth_api_key_id"`
	AuthAPIKey       string                  `mapstructure:"auth_api_key"`
	Username         string                  `mapstructure:"auth_basic_username"`
	Password         string                  `mapstructure:"auth_basic_password"`
	EnableSniffing   bool                    `mapstructure:"use_sniffing"`
	RollingIndex     bool                    `mapstructure:"rolling_index"`
	DisableBulk      bool                    `mapstructure:"disable_bulk"`
}

// ElasticsearchBulkConfig defines elasticsearch bulk config.
type ElasticsearchBulkConfig struct {
	Workers       int `mapstructure:"workers"`
	FlushInterval int `mapstructure:"flush_interval"`
	BulkActions   int `mapstructure:"bulk_actions"`
	BulkSize      int `mapstructure:"bulk_size"`
}

// ElasticsearchOperator defines interface for all elasticsearch operator.
type ElasticsearchOperator interface {
	processData(ctx context.Context, data []interface{}, esConf *ElasticsearchConf) error
}

// Elasticsearch7Operator defines elasticsearch6 operator.
type Elasticsearch7Operator struct {
	esClient      *elastic.Client
	bulkProcessor *elastic.BulkProcessor
}

// APIKeyTransport defiens elasticsearch api key.
type APIKeyTransport struct {
	APIKey   string
	APIKeyID string
}

// RoundTrip for APIKeyTransport auth.
func (t *APIKeyTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	auth := t.APIKeyID + ":" + t.APIKey
	key := base64.StdEncoding.EncodeToString([]byte(auth))

	r.Header.Set("Authorization", "ApiKey "+key)

	return http.DefaultTransport.RoundTrip(r)
}

func getOperator(ctx context.Context, conf ElasticsearchConf) (ElasticsearchOperator, error) {
	var err error
	urls := strings.Split(conf.ElasticsearchURL, ",")
	httpClient := http.DefaultClient
	if conf.AuthAPIKey != "" && conf.AuthAPIKeyID != "" {
		conf.Username = ""
		conf.Password = ""
		httpClient = &http.Client{Transport: &APIKeyTransport{APIKey: conf.AuthAPIKey, APIKeyID: conf.AuthAPIKeyID}}
	}

	e := new(Elasticsearch7Operator)

	e.esClient, err = elastic.NewClient(
		elastic.SetURL(urls...),
		elastic.SetSniff(conf.EnableSniffing),
		elastic.SetBasicAuth(conf.Username, conf.Password),
		elastic.SetHttpClient(httpClient),
	)

	if err != nil {
		return e, errors.Wrap(err, "failed to new es client")
	}
	// Setup a bulk processor
	p := e.esClient.BulkProcessor().Name("LEONAGStashESv6BackgroundProcessor")
	if conf.BulkConfig.Workers != 0 {
		p = p.Workers(conf.BulkConfig.Workers)
	}

	if conf.BulkConfig.FlushInterval != 0 {
		p = p.FlushInterval(time.Duration(conf.BulkConfig.FlushInterval) * time.Second)
	}

	if conf.BulkConfig.BulkActions != 0 {
		p = p.BulkActions(conf.BulkConfig.BulkActions)
	}

	if conf.BulkConfig.BulkSize != 0 {
		p = p.BulkSize(conf.BulkConfig.BulkSize)
	}

	e.bulkProcessor, err = p.Do(ctx)

	return e, errors.Wrap(err, "failed to start bulk processor")
}

// New create a elasticsearch gstash instance.
func (e *ElasticsearchGstash) New() Gstash {
	newGStash := ElasticsearchGstash{}

	return &newGStash
}

// GetName returns the elasticsearch gstash name.
func (e *ElasticsearchGstash) GetName() string {
	return "Elasticsearch Gstash"
}

// Init initialize the elasticsearch gstash instance.
func (e *ElasticsearchGstash) Init(config interface{}) error {
	e.esConf = &ElasticsearchConf{}
	loadConfigErr := mapstructure.Decode(config, &e.esConf)

	if loadConfigErr != nil {
		log.Fatalf("Failed to decode configuration: %s", loadConfigErr.Error())
	}

	if e.esConf.IndexName == "" {
		e.esConf.IndexName = "iam_analytics"
	}

	if e.esConf.ElasticsearchURL == "" {
		e.esConf.ElasticsearchURL = "http://localhost:9200"
	}

	if e.esConf.DocumentType == "" {
		e.esConf.DocumentType = "iam_analytics"
	}

	re := regexp.MustCompile(`(.*)\/\/(.*):(.*)\@(.*)`)
	printableURL := re.ReplaceAllString(e.esConf.ElasticsearchURL, `$1//***:***@$4`)

	log.Infof("Elasticsearch URL: %s", printableURL)
	log.Infof("Elasticsearch Index: %s", e.esConf.IndexName)
	if e.esConf.RollingIndex {
		log.Infof("Index will have date appended to it in the format %s -YYYY.MM.DD", e.esConf.IndexName)
	}

	e.connect(context.Background())

	return nil
}

func (e *ElasticsearchGstash) connect(ctx context.Context) {
	var err error

	e.operator, err = getOperator(ctx, *e.esConf)
	if err != nil {
		log.Errorf("Elasticsearch connection failed: %s", err.Error())
		time.Sleep(5 * time.Second)
		e.connect(ctx)
	}
}

// WriteData write analyzed data to elasticsearch persistent back-end storage.
func (e *ElasticsearchGstash) WriteData(ctx context.Context, data []interface{}) error {
	log.Infof("Writing %d records", len(data))

	if e.operator == nil {
		log.Debug("Connecting to analytics store")
		e.connect(ctx)
		_ = e.WriteData(ctx, data)
	} else if len(data) > 0 {
		_ = e.operator.processData(ctx, data, e.esConf)
	}

	return nil
}

func getIndexName(esConf *ElasticsearchConf) string {
	indexName := esConf.IndexName

	if esConf.RollingIndex {
		currentTime := time.Now()
		// This formats the date to be YYYY.MM.DD but Golang makes you use a specific date for its date formatting
		indexName += "-" + currentTime.Format("2006.01.02")
	}

	return indexName
}

func getMapping(datum analytics.AnalyticsRecord) (map[string]interface{}, string) {
	record := datum
	mapping := map[string]interface{}{
		"@timestamp": record.TimeStamp,
		"username":   record.Username,
		"effect":     record.Effect,
		"conclusion": record.Conclusion,
		"request":    record.Request,
		"policies":   record.Policies,
		"deciders":   record.Deciders,
		"expireAt":   record.ExpireAt,
	}

	return mapping, ""
}

func (e Elasticsearch7Operator) processData(ctx context.Context, data []interface{}, esConf *ElasticsearchConf) error {
	index := e.esClient.Index().Index(getIndexName(esConf))

	for dataIndex := range data {
		if ctxErr := ctx.Err(); ctxErr != nil {
			continue
		}

		d, ok := data[dataIndex].(analytics.AnalyticsRecord)
		if !ok {
			log.Errorf("Error while writing %s: data not of type analytics.AnalyticsRecord", data[dataIndex])

			continue
		}

		mapping, id := getMapping(d)

		if !esConf.DisableBulk {
			r := elastic.NewBulkIndexRequest().Index(getIndexName(esConf)).Type(esConf.DocumentType).Id(id).Doc(mapping)
			e.bulkProcessor.Add(r)
		} else {
			//nolint: staticcheck
			_, err := index.BodyJson(mapping).Type(esConf.DocumentType).Id(id).Do(ctx)
			if err != nil {
				log.Errorf("Error while writing %s %s", data[dataIndex], err.Error())
			}
		}
	}

	return nil
}
