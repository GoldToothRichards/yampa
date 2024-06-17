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

package sources

import (
	"encoding/json"
	"sync"
	"time"

	"yampa-cli/pkg/core"

	"github.com/gorilla/websocket"
)

type CoincapCrypto struct {
	Exchange  string  `json:"exchange"`
	Base      string  `json:"base"`
	Quote     string  `json:"quote"`
	Direction string  `json:"direction"`
	Price     float64 `json:"price"`
	Volume    float64 `json:"volume"`
	Timestamp int64   `json:"timestamp"`
	PriceUsd  float64 `json:"priceUsd"`
}

func StreamCoincapCrypto(
	trades chan core.TradeJson,
) {
	exchanges := []string{
		"binance", "hitbtc", "gdax", "huobi", "bitfinex", "bitstamp",
		"gemini", "poloniex", "bitso", "luno", "therocktrading",
		"coinmate",
	}

	// Stream trades from each exchange on a separate thread
	var wg sync.WaitGroup
	for i := 0; i < len(exchanges); i++ {
		var exchange = exchanges[i]
		wg.Add(1)
		go func() {
			StreamCoincapExchange(exchange, trades)
			wg.Done()
		}()
	}
	wg.Wait()
}

func StreamCoincapExchange(
	exchange string,
	trades chan core.TradeJson,
) {
	// Infinite loop refreshing websocket connections
	for {
		url := "wss://ws.coincap.io/trades/" + exchange
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			core.WarnLogger.Println("Could not connect to", url, "-", err)
			time.Sleep(time.Millisecond * 10)
			continue
		}
		core.InfoLogger.Println("Connected to", url)

		// Infinite loop reading trades
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				core.WarnLogger.Println("Could not read message from", exchange, "-", err)
				conn.Close()
				break
			}

			// Convert input message bytes to json
			var trade_in CoincapCrypto
			json.Unmarshal(message, &trade_in)

			// Convert input message to standard json schema
			trade_out := core.TradeJson{
				Source:      "coincap-crypto",
				Base:        trade_in.Base,
				Quote:       trade_in.Quote,
				Exchange:    trade_in.Exchange,
				VolumeBase:  trade_in.Volume,
				VolumeQuote: trade_in.Price * trade_in.Volume,
				Price:       trade_in.Price,
				Timestamp:   time.UnixMilli(trade_in.Timestamp).String(),
			}

			// Push trade through the channel
			trades <- trade_out

		}
	}
}
