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
	"encoding/json"
	"fmt"
	"yampa-cli/pkg/core"

	"github.com/pusher/pusher-http-go/v5"
	"github.com/spf13/viper"
)


func WriteToPusher(
	trades chan core.TradeJson,
){
	pusherClient := pusher.Client{
		AppID: viper.GetString("PUSHER_APP_ID"),
		Key: viper.GetString("PUSHER_KEY"),
		Secret: viper.GetString("PUSHER_SECRET"),
		Cluster: viper.GetString("PUSHER_CLUSTER"),
		Secure: true,
	}

	fmt.Println("Pushing '", viper.GetString("PUSHER_EVENT_NAME"), "' events to '", viper.GetString("PUSHER_CHANNEL"), "' channel.")

	// Infinite loop writing trades from the channel to Pusher channel
	for {
		var tradeBatch []core.TradeJson

		// Collect trades until batch size is reached
		for i := 0; i < viper.GetInt("PUSHER_BATCH_SIZE"); i++ {
			trade := <-trades
			tradeBatch = append(tradeBatch, trade)
		}

		// Marshal the entire batch as JSON
		jsonBytes, err := json.Marshal(tradeBatch)
		if err != nil {
			fmt.Println("Error marshaling trade batch:", err)
			continue
		}
		jsonStr := string(jsonBytes)

		// Trigger Pusher event with the batch of trades
		err = pusherClient.Trigger(viper.GetString("PUSHER_CHANNEL"), viper.GetString("PUSHER_EVENT_NAME"), jsonStr)
		if err != nil {
			fmt.Println("Error triggering Pusher event:", err)
		}
	}
}