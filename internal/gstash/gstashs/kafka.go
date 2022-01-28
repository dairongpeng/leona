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
	"crypto/tls"
	"time"

	"github.com/dairongpeng/leona/pkg/errors"
	"github.com/dairongpeng/leona/pkg/json"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/segmentio/kafka-go/snappy"

	"github.com/dairongpeng/leona/internal/gstash/analytics"
	"github.com/dairongpeng/leona/pkg/log"
)

// KafkaGstash defines a kafka gstash with kafka specific options and common options.
type KafkaGstash struct {
	kafkaConf    *KafkaConf
	writerConfig kafka.WriterConfig
	CommonGstashConfig
}

// Message contains the messages need to push to gstash.
type Message map[string]interface{}

// KafkaConf defines kafka specific options.
type KafkaConf struct {
	Broker                []string          `mapstructure:"broker"`
	ClientID              string            `mapstructure:"client_id"`
	Topic                 string            `mapstructure:"topic"`
	SSLCertFile           string            `mapstructure:"ssl_cert_file"`
	SSLKeyFile            string            `mapstructure:"ssl_key_file"`
	SASLMechanism         string            `mapstructure:"sasl_mechanism"`
	Username              string            `mapstructure:"sasl_username"`
	Password              string            `mapstructure:"sasl_password"`
	Algorithm             string            `mapstructure:"sasl_algorithm"`
	Timeout               time.Duration     `mapstructure:"timeout"`
	MetaData              map[string]string `mapstructure:"meta_data"`
	Compressed            bool              `mapstructure:"compressed"`
	UseSSL                bool              `mapstructure:"use_ssl"`
	SSLInsecureSkipVerify bool              `mapstructure:"ssl_insecure_skip_verify"`
}

// New create a kafka gstash instance.
func (k *KafkaGstash) New() Gstash {
	newGstash := KafkaGstash{}

	return &newGstash
}

// GetName returns the kafka gstash name.
func (k *KafkaGstash) GetName() string {
	return "Kafka Gstash"
}

// Init initialize the kafka gstash instance.
func (k *KafkaGstash) Init(config interface{}) error {
	// Read configuration file
	k.kafkaConf = &KafkaConf{}
	err := mapstructure.Decode(config, &k.kafkaConf)
	if err != nil {
		log.Fatalf("Failed to decode configuration: %s", err.Error())
	}

	var tlsConfig *tls.Config
	// nolint: nestif
	if k.kafkaConf.UseSSL {
		if k.kafkaConf.SSLCertFile != "" && k.kafkaConf.SSLKeyFile != "" {
			var cert tls.Certificate
			log.Debug("Loading certificates for mTLS.")
			cert, err = tls.LoadX509KeyPair(k.kafkaConf.SSLCertFile, k.kafkaConf.SSLKeyFile)
			if err != nil {
				log.Debugf("Error loading mTLS certificates: %s", err.Error())

				return errors.Wrap(err, "failed loading mTLS certificates")
			}
			tlsConfig = &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: k.kafkaConf.SSLInsecureSkipVerify,
			}
		} else if k.kafkaConf.SSLCertFile != "" || k.kafkaConf.SSLKeyFile != "" {
			log.Error("Only one of ssl_cert_file and ssl_cert_key configuration option is setted, you should set both to enable mTLS.")
		} else {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: k.kafkaConf.SSLInsecureSkipVerify,
			}
		}
	} else if k.kafkaConf.SASLMechanism != "" {
		log.Warn("SASL-Mechanism is setted but use_ssl is false.", log.String("SASL-Mechanism", k.kafkaConf.SASLMechanism))
	}

	var mechanism sasl.Mechanism

	switch k.kafkaConf.SASLMechanism {
	case "":
		break
	case "PLAIN", "plain":
		mechanism = plain.Mechanism{Username: k.kafkaConf.Username, Password: k.kafkaConf.Password}
	case "SCRAM", "scram":
		algorithm := scram.SHA256
		if k.kafkaConf.Algorithm == "sha-512" || k.kafkaConf.Algorithm == "SHA-512" {
			algorithm = scram.SHA512
		}
		var mechErr error
		mechanism, mechErr = scram.Mechanism(algorithm, k.kafkaConf.Username, k.kafkaConf.Password)
		if mechErr != nil {
			log.Fatalf("Failed initialize kafka mechanism: %s", mechErr.Error())
		}
	default:
		log.Warn(
			"LEONA gstash doesn't support this SASL mechanism.",
			log.String("SASL-Mechanism", k.kafkaConf.SASLMechanism),
		)
	}

	// Kafka writer connection config
	dialer := &kafka.Dialer{
		Timeout:       k.kafkaConf.Timeout,
		ClientID:      k.kafkaConf.ClientID,
		TLS:           tlsConfig,
		SASLMechanism: mechanism,
	}

	k.writerConfig.Brokers = k.kafkaConf.Broker
	k.writerConfig.Topic = k.kafkaConf.Topic
	k.writerConfig.Balancer = &kafka.LeastBytes{}
	k.writerConfig.Dialer = dialer
	k.writerConfig.WriteTimeout = k.kafkaConf.Timeout
	k.writerConfig.ReadTimeout = k.kafkaConf.Timeout
	if k.kafkaConf.Compressed {
		k.writerConfig.CompressionCodec = snappy.NewCompressionCodec()
	}

	log.Infof("Kafka config: %s", k.writerConfig)

	return nil
}

// WriteData write analyzed data to kafka persistent back-end storage.
func (k *KafkaGstash) WriteData(ctx context.Context, data []interface{}) error {
	startTime := time.Now()
	log.Infof("Writing %d records ...", len(data))
	kafkaMessages := make([]kafka.Message, len(data))
	for i, v := range data {
		// Build message format
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
		// Add static metadata to json
		for key, value := range k.kafkaConf.MetaData {
			message[key] = value
		}

		// Transform object to json string
		json, jsonError := json.Marshal(message)
		if jsonError != nil {
			log.Error("unable to marshal message", log.String("error", jsonError.Error()))
		}

		// Kafka message structure
		kafkaMessages[i] = kafka.Message{
			Time:  time.Now(),
			Value: json,
		}
	}
	// Send kafka message
	kafkaError := k.write(ctx, kafkaMessages)
	if kafkaError != nil {
		log.Error("unable to write message", log.String("error", kafkaError.Error()))
	}
	log.Debugf("ElapsedTime in seconds for %d records %v", len(data), time.Since(startTime))

	return nil
}

// write kafka写入插件写入动作
func (k *KafkaGstash) write(ctx context.Context, messages []kafka.Message) error {
	kafkaWriter := kafka.NewWriter(k.writerConfig)
	defer kafkaWriter.Close()

	return kafkaWriter.WriteMessages(ctx, messages...)
}
