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

package cmd

import (
	"context"
	"strings"

	"yampa-cli/pkg/core"

	"github.com/hamba/avro"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

type Healthcheck struct {
	Message string `avro:"message"`
}

var HealthcheckSchema = `{
	"type": "record",
	"name": "Healthcheck",
	"fields" : [
		{"name": "message", "type": "string"}
	]
}`

// healthcheckCmd represents the healthcheck command
var healthcheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Check the connection to Kafka",
	Long: `
Create a Kafka producer and send a message to the Kafka broker.
Message will be sent to the "healthcheck" topic by default and
will contain the value "healthcheck".
`,
	Run: func(cmd *cobra.Command, args []string) {

		message, _ := cmd.Flags().GetString("message")
		topic, _ := cmd.Flags().GetString("topic")

		// Ensure our schema is registered.
		rcl, err := sr.NewClient(sr.URLs(viper.GetString("SCHEMA_REGISTRY_URL")))
		if err != nil {
			core.ErrorLogger.Fatal("Unable to create schema registry client. ", err)
		}
		ss, err := rcl.CreateSchema(context.Background(), topic+"-value", sr.Schema{
			Schema: HealthcheckSchema,
			Type:   sr.TypeAvro,
		})
		if err != nil {
			core.ErrorLogger.Fatal("Unable to create avro schema. ", err)
		}
		core.InfoLogger.Printf("Created or reusing schema subject %q version %d id %d\n", ss.Subject, ss.Version, ss.ID)

		// Setup our serializer / deserializer.
		avroSchema, err := avro.Parse(HealthcheckSchema)
		if err != nil {
			core.ErrorLogger.Fatal("Unable to parse avro schema. ", err)
		}
		var serde sr.Serde
		serde.Register(
			ss.ID,
			Healthcheck{},
			sr.EncodeFn(func(v interface{}) ([]byte, error) {
				return avro.Marshal(avroSchema, v)
			}),
			sr.DecodeFn(func(b []byte, v interface{}) error {
				return avro.Unmarshal(avroSchema, b, v)
			}),
		)

		// Setup kafka client
		cl, err := kgo.NewClient(
			kgo.SeedBrokers(strings.Split(viper.GetString("KAFKA_BROKER_URLS"), ",")...),
			kgo.DefaultProduceTopic(topic),
			kgo.ConsumeTopics(topic),
		)
		if err != nil {
			core.ErrorLogger.Fatal("Unable to initialize kafka client. ", err)
		}

		// Send a message to the Kafka broker
		cl.ProduceSync(
			context.Background(),
			&kgo.Record{
				Value: serde.MustEncode(Healthcheck{
					Message: message,
				}),
			},
		)
	},
}

func init() {
	healthcheckCmd.Flags().StringP("message", "m", "healthcheck", "message to send in the healthcheck")
	healthcheckCmd.Flags().StringP("topic", "t", "healthcheck", "kafka topic to stream into")
	rootCmd.AddCommand(healthcheckCmd)
}
