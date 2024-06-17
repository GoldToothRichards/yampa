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
	"yampa-cli/pkg/app"
	"yampa-cli/pkg/sources"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// coincapCmd represents the coincap command
var coincapCmd = &cobra.Command{
	Use:   "coincap",
	Short: "Stream all trades available from Coincap",
	Long: `
Stream all trades available from Coincap. All messages will be converted to
a standard schema before being written to the user specified output destination(s).
Multiple destinations can be chosen using the optional flags.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Read cmd arguments
		quiet, _ := cmd.Flags().GetBool("quiet")
		file, _ := cmd.Flags().GetString("file")
		kafka, _ := cmd.Flags().GetBool("kafka")
		topic, _ := cmd.Flags().GetString("topic")
		pusher, _ := cmd.Flags().GetBool("pusher")
		channel, _ := cmd.Flags().GetString("channel")
		event, _ := cmd.Flags().GetString("event")

		// Define the kafka topic from cmd argument if present, otherwise use a default
		if topic != "" {
			viper.Set("topic", topic)
		} else {
			viper.Set("topic", "trades")
		}

		// Define the pusher config from cmd arguments if present, otherwise use the environment
		if channel != "" {
			viper.Set("PUSHER_CHANNEL", channel)
		}

		if event != "" {
			viper.Set("PUSHER_EVENT_NAME", event)
		}

		// Write trades to the user specified destination(s)
		app.StreamTrades(sources.StreamCoincapCrypto, quiet, file, kafka, pusher)

	},
}

func init() {
	streamCmd.AddCommand(coincapCmd)
}
