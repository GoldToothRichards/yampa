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
	"github.com/spf13/cobra"
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream all trades available from a given data source",
	Long: `
Stream all trades available from a given data source. All messages will be
converted to a standard schema before being written to the user specified output
destination(s). Multiple destinations can be chosen using the optional flags.`,
}

func init() {
	rootCmd.AddCommand(streamCmd)
	streamCmd.PersistentFlags().Bool("quiet", false, "do not print stream to stdout")
	streamCmd.PersistentFlags().String("file", "", "stream to a local log file")
	streamCmd.PersistentFlags().Bool("kafka", false, "stream to kafka")
	streamCmd.PersistentFlags().StringP("topic", "t", "", "kafka topic to stream into")
	streamCmd.PersistentFlags().Bool("pusher", false, "stream to pusher")
	streamCmd.PersistentFlags().String("channel", "", "pusher channel to stream into")
	streamCmd.PersistentFlags().String("event", "", "pusher event name")
}
