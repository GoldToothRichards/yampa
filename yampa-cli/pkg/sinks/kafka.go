/*
Copyright © 2023 Jacob Crabtree <crabtr26@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package sinks

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"yampa-cli/pkg/core"

	"github.com/hamba/avro"
	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

type KafkaMessageConfig struct {
	Format     string // "AVRO" or "JSON"
	IncludeKey bool   // Whether to include a key in the Kafka messages
}

func NewKafkaMessageConfig() KafkaMessageConfig {
	return KafkaMessageConfig{
		Format:     "JSON", // Default format is JSON
		IncludeKey: false,  // Default is not to include a key
	}
}

const timestampFormat = "2006-01-02 15:04:05.000 -0700 MST"

func WriteToKafka(
	trades chan core.TradeJson,
	config KafkaMessageConfig,
) {

	serde := GetSerde(config)
	producer := GetProducer()

	switch config.Format {
	case "AVRO":
		switch config.IncludeKey {
		case true:
			// Infinite loop writing trades from the channel to a kafka topic
			for {
				trade_json := <-trades

				// Parse the JSON string timestamp to time.Time
				parsedTime, err := time.Parse(timestampFormat, trade_json.Timestamp)
				if err != nil {
					core.ErrorLogger.Printf("Error parsing timestamp: %v", err)
					json_bytes, _ := json.Marshal(trade_json)
					json_str := string(json_bytes)
					core.ErrorLogger.Printf("Trade: %v", json_str)
					continue
				}

				// Convert json to avro
				trade_avro := core.TradeAvro{
					Source:      trade_json.Source,
					Base:        trade_json.Base,
					Quote:       trade_json.Quote,
					Exchange:    trade_json.Exchange,
					VolumeBase:  trade_json.VolumeBase,
					VolumeQuote: trade_json.VolumeQuote,
					Price:       trade_json.Price,
					Timestamp:   parsedTime.UnixMilli(),
				}

				// Define avro key from the full avro record
				trade_key := core.TradeAvroKey{
					Source:   trade_avro.Source,
					Base:     trade_avro.Base,
					Quote:    trade_avro.Quote,
					Exchange: trade_avro.Exchange,
				}
				producer.Produce(
					context.Background(),
					&kgo.Record{
						Key:   serde.MustEncode(trade_key),
						Value: serde.MustEncode(trade_avro),
					},
					func(r *kgo.Record, err error) {
						if err != nil {
							core.WarnLogger.Println("Error while producing to Kafka", err)
						}
					},
				)
			}

		case false:
			// Infinite loop writing trades from the channel to a kafka topic
			for {
				trade_json := <-trades

				// Parse the JSON string timestamp to time.Time
				parsedTime, err := time.Parse(timestampFormat, trade_json.Timestamp)
				if err != nil {
					core.ErrorLogger.Printf("Error parsing timestamp: %v", err)
					json_bytes, _ := json.Marshal(trade_json)
					json_str := string(json_bytes)
					core.ErrorLogger.Printf("Trade: %v", json_str)
					continue
				}

				// Convert json to avro
				trade_avro := core.TradeAvro{
					Source:      trade_json.Source,
					Base:        trade_json.Base,
					Quote:       trade_json.Quote,
					Exchange:    trade_json.Exchange,
					VolumeBase:  trade_json.VolumeBase,
					VolumeQuote: trade_json.VolumeQuote,
					Price:       trade_json.Price,
					Timestamp:   parsedTime.UnixMilli(),
				}

				producer.Produce(
					context.Background(),
					&kgo.Record{
						Value: serde.MustEncode(trade_avro),
					},
					func(r *kgo.Record, err error) {
						if err != nil {
							core.WarnLogger.Println("Error while producing to Kafka", err)
						}
					},
				)
			}

		}

	case "JSON":
		switch config.IncludeKey {
		case true:
			// Infinite loop writing trades from the channel to a kafka topic
			for {
				trade_json := <-trades

				// Define JSON key from the JSON record
				trade_key := core.TradeJsonKey{
					Source:   trade_json.Source,
					Base:     trade_json.Base,
					Quote:    trade_json.Quote,
					Exchange: trade_json.Exchange,
				}

				// Convert JSON to bytes
				trade_key_bytes, err := json.Marshal(trade_key)
				if err != nil {
					core.WarnLogger.Println("Error while marshalling key to JSON", err)
				}
				trade_json_bytes, err := json.Marshal(trade_json)
				if err != nil {
					core.WarnLogger.Println("Error while marshalling value to JSON", err)
				}

				producer.Produce(
					context.Background(),
					&kgo.Record{
						Key:   trade_key_bytes,
						Value: trade_json_bytes,
					},
					func(r *kgo.Record, err error) {
						if err != nil {
							core.WarnLogger.Println("Error while producing to Kafka", err)
						}
					},
				)
			}
		case false:
			// Infinite loop writing trades from the channel to a kafka topic
			for {
				trade_json := <-trades

				// Convert JSON to bytes
				trade_json_bytes, err := json.Marshal(trade_json)
				if err != nil {
					core.WarnLogger.Println("Error while marshalling value to JSON", err)
				}

				producer.Produce(
					context.Background(),
					&kgo.Record{
						Value: trade_json_bytes,
					},
					func(r *kgo.Record, err error) {
						if err != nil {
							core.WarnLogger.Println("Error while producing to Kafka", err)
						}
					},
				)
			}
		}
	default:
		core.ErrorLogger.Fatal("Invalid message format '", config.Format, "'. Must be either 'AVRO' or 'JSON'.")

	}
}

func GetSerde(config KafkaMessageConfig) *sr.Serde {
	// Create schema registry client
	rcl, err := sr.NewClient(sr.URLs(viper.GetString("SCHEMA_REGISTRY_URL")))
	if err != nil {
		core.ErrorLogger.Fatal("Unable to create schema registry client. ", err)
	}

	var serde sr.Serde

	if config.Format == "AVRO" {
		if config.IncludeKey {
			registerAvroSchema(rcl, viper.GetString("topic")+"-key", core.TradeAvroKeySchema, true, &serde)
		}
		registerAvroSchema(rcl, viper.GetString("topic")+"-value", core.TradeAvroSchema, false, &serde)
	} else if config.Format == "JSON" {
		if config.IncludeKey {
			// registerJsonSchema(rcl, viper.GetString("topic")+"-key", core.TradeJsonKey, true, &serde)
			core.InfoLogger.Println("JSON key schema not supported yet")
		}
		// registerJsonSchema(rcl, viper.GetString("topic")+"-value", core.TradeJson, false, &serde)
		core.InfoLogger.Println("JSON value schema not supported yet")
	} else {
		core.ErrorLogger.Fatal("Unsupported schema format: ", config.Format)
	}

	return &serde
}

func registerAvroSchema(rcl *sr.Client, subject string, schemaStr string, isKeySchema bool, serde *sr.Serde) {
	schema, err := rcl.CreateSchema(context.Background(), subject, sr.Schema{
		Schema: schemaStr,
		Type:   sr.TypeAvro,
	})
	if err != nil {
		core.ErrorLogger.Fatal("Unable to create schema: ", err)
	}
	core.InfoLogger.Printf("Created or reusing schema subject %q version %d id %d\n", schema.Subject, schema.Version, schema.ID)

	avroSchema, err := avro.Parse(schemaStr)
	if err != nil {
		core.ErrorLogger.Fatal("Unable to parse avro schema. ", err)
	}

	if isKeySchema {
		serde.Register(
			schema.ID,
			core.TradeAvroKey{},
			sr.EncodeFn(func(v interface{}) ([]byte, error) {
				return avro.Marshal(avroSchema, v)
			}),
			sr.DecodeFn(func(b []byte, v interface{}) error {
				return avro.Unmarshal(avroSchema, b, v)
			}),
		)
	} else {
		serde.Register(
			schema.ID,
			core.TradeAvro{},
			sr.EncodeFn(func(v interface{}) ([]byte, error) {
				return avro.Marshal(avroSchema, v)
			}),
			sr.DecodeFn(func(b []byte, v interface{}) error {
				return avro.Unmarshal(avroSchema, b, v)
			}),
		)
	}
}

func GetProducer() *kgo.Client {

	// Setup kafka client
	producer, err := kgo.NewClient(
		kgo.SeedBrokers(strings.Split(viper.GetString("KAFKA_BROKER_URLS"), ",")...),
		kgo.DefaultProduceTopic(viper.GetString("topic")),
		kgo.ProducerBatchCompression(kgo.NoCompression()),
	)
	if err != nil {
		core.ErrorLogger.Fatal("Unable to initialize kafka client. ", err)
	}

	return producer
}
