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

package app

import (
	"sync"

	"yampa-cli/pkg/core"
	"yampa-cli/pkg/sinks"

	"github.com/spf13/viper"
)

type Stream func(chan core.TradeJson)

func StreamTrades(
	source Stream,
	quiet bool,
	file string,
	kafka bool,
	pusher bool,
) {
	// Create a channel to transmit trades
	trades := make(chan core.TradeJson)

	// Create a threadpool to read/write trades
	var wg sync.WaitGroup

	// Add trades from a source stream to the channel on a background thread
	wg.Add(1)
	go func() {
		source(trades)
		wg.Done()
	}()

	// Write trades from the channel to stdout on a background thread
	if !quiet {
		wg.Add(1)
		go func() {
			sinks.WriteToStdout(trades)
			wg.Done()
		}()
	}

	// Write trades from the channel to a local log file on a background thread
	if file != "" {
		wg.Add(1)
		go func() {
			sinks.WriteToFile(trades, file)
			wg.Done()
		}()
	}

	// Write trades from the channel to kafka on a background thread
	if kafka {
		config := sinks.NewKafkaMessageConfig()
		config.Format = viper.GetString("KAFKA_MESSAGE_FORMAT")
		config.IncludeKey = viper.GetBool("KAFKA_MESSAGE_INCLUDE_KEY")
		wg.Add(1)
		go func() {
			sinks.WriteToKafka(trades, config)
			wg.Done()
		}()
	}

	// Write trades from the channel to Pusher on a background thread
	if pusher {
		wg.Add(1)
		go func() {
			sinks.WriteToPusher(trades)
			wg.Done()
		}()
	}

	// Now that all threads have started working, wait forever
	wg.Wait()
}
