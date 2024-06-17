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
	"fmt"
	"strings"
	"time"

	"yampa-cli/pkg/core"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

type AlpacaStock struct {
	MessageType    string   `json:"T"`
	Symbol         string   `json:"S"`
	TradeID        int64    `json:"i"`
	Exchange       string   `json:"x"`
	Price          float64  `json:"p"`
	Size           int64    `json:"s"`
	Timestamp      string   `json:"t"`
	TradeCondition []string `json:"c"`
	Tape           string   `json:"z"`
}

type AlpacaCrypto struct {
	MessageType string  `json:"T"`
	Pair        string  `json:"S"`
	Price       float64 `json:"p"`
	Size        float64 `json:"s"`
	Timestamp   string  `json:"t"`
	TradeID     int64   `json:"i"`
	Tks         string  `json:"tks"`
}

func StreamAlpacaStocks(
	trades chan core.TradeJson,
) {
	symbols := []string{
		"AAPL", "IBM", "AMD", "TSLA", "AMZN", "GOOGL", "NVDA", "NU",
		"BAC", "F", "NFLX", "MSFT", "T", "V", "XOM", "UNH", "JNJ",
		"WMT", "META", "MA", "JPM", "HD", "KO", "COST", "MCD", "PEP",
		"GOOG",
	}

	// Infinite loop refreshing websocket connections
	for {
		url := "wss://stream.data.alpaca.markets/v2/iex"
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			core.WarnLogger.Println("Could not connect to", url, "-", err)
			time.Sleep(time.Millisecond * 10)
			continue
		}
		core.InfoLogger.Println("Connected to", url)

		// Initial welcome message
		_, welcome_message, err := conn.ReadMessage()
		if err != nil {
			core.WarnLogger.Println("Unable to read welcome message", err)
			conn.Close()
			time.Sleep(time.Millisecond * 10)
			continue
		} else {
			core.InfoLogger.Println(string(welcome_message))
		}

		// Authenticate
		type AuthMessage struct {
			Action string `json:"action"`
			Key    string `json:"key"`
			Secret string `json:"secret"`
		}
		auth_message := AuthMessage{
			Action: "auth",
			Key:    viper.GetString("ALPACA_KEY"),
			Secret: viper.GetString("ALPACA_SECRET"),
		}
		core.InfoLogger.Println(viper.GetString("ALPACA_KEY"))
		core.InfoLogger.Println(viper.GetString("ALPACA_SECRET"))
		conn.WriteJSON(auth_message)

		// Authentication response
		_, auth_response_message, err := conn.ReadMessage()
		if err != nil {
			core.WarnLogger.Println("Unable to read welcome message", err)
			conn.Close()
			time.Sleep(time.Millisecond * 10)
			continue
		} else {
			core.InfoLogger.Println(string(auth_response_message))
		}

		// Subscribe
		type SubscribeMessage struct {
			Action string   `json:"action"`
			Trades []string `json:"trades"`
		}
		subscribe_message := SubscribeMessage{
			Action: "subscribe",
			Trades: symbols,
		}
		conn.WriteJSON(subscribe_message)

		// Subscription response
		_, sub_response_message, err := conn.ReadMessage()
		if err != nil {
			core.WarnLogger.Println("Unable to read welcome message", err)
			conn.Close()
			time.Sleep(time.Millisecond * 10)
			continue
		} else {
			core.InfoLogger.Println(string(sub_response_message))
		}

		// Infinite loop reading trades
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				core.WarnLogger.Println("Could not read message", err)
				conn.Close()
				break
			}

			fmt.Println("MESSAGE", string(message))

			// Convert input message bytes to json
			var trades_in []AlpacaStock
			json.Unmarshal(message, &trades_in)
			trade_in := trades_in[0]

			// Convert input message to standard json schema
			trade_out := core.TradeJson{
				Source:      "alpaca-stocks",
				Base:        trade_in.Symbol,
				Quote:       "USD",
				Exchange:    trade_in.Exchange,
				VolumeBase:  float64(trade_in.Size),
				VolumeQuote: trade_in.Price * float64(trade_in.Size),
				Price:       trade_in.Price,
				Timestamp:   trade_in.Timestamp,
			}

			// Push trade through the channel
			trades <- trade_out
		}
	}
}

func StreamAlpacaCrypto(
	trades chan core.TradeJson,
) {
	// all_pairs := []string{
	// 	"LINK/BTC", "SOL/USD", "DAI/USD", "UNI/USDT", "LINK/USDT",
	// 	"DOGE/USDT", "LTC/USDT", "SUSHI/USD", "LINK/USD", "LTC/USD",
	// 	"SOL/BTC", "MATIC/BTC", "ETH/USD", "SOL/USDT", "AVAX/USD",
	// 	"BCH/BTC", "TRX/USDT", "PAXG/USD", "DOGE/USD", "ALGO/USDT",
	// 	"TRX/USD", "BTC/USDT", "AAVE/USD", "PAXG/USDT", "NEAR/USD",
	// 	"BCH/USDT", "AVAX/USDT", "GRT/USD", "ETH/BTC", "LTC/BTC",
	// 	"WBTC/USD", "BTC/USD", "UNI/BTC", "AAVE/USDT", "USDT/USD",
	// 	"BAT/USD", "UNI/USD", "YFI/USD", "NEAR/USDT", "SHIB/USDT",
	// 	"BCH/USD", "SUSHI/USDT", "ETH/USDT", "YFI/USDT", "ALGO/USD",
	// 	"MATIC/USD", "MKR/USD",
	// }
	pairs := []string{
		"SOL/USD", "DAI/USD", "UNI/USDT",
		"DOGE/USDT", "SUSHI/USD",
		"SOL/BTC", "ETH/USD", "SOL/USDT", "AVAX/USD",
		"DOGE/USD", "ALGO/USDT",
		"BTC/USDT", "AAVE/USD", "NEAR/USD",
		"AVAX/USDT", "ETH/BTC",
		"WBTC/USD", "BTC/USD", "UNI/BTC", "AAVE/USDT", "USDT/USD",
		"BAT/USD", "UNI/USD", "NEAR/USDT", "SHIB/USDT",
		"SUSHI/USDT", "ETH/USDT", "ALGO/USD",
		"MKR/USD",
	}

	// Infinite loop refreshing websocket connections
	for {
		url := "wss://stream.data.alpaca.markets/v1beta2/crypto"
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			core.WarnLogger.Println("Could not connect to", url, "-", err)
			time.Sleep(time.Millisecond * 10)
			continue
		}
		core.InfoLogger.Println("Connected to", url)

		// Initial welcome message
		_, welcome_message, err := conn.ReadMessage()
		if err != nil {
			core.WarnLogger.Println("Unable to read welcome message", err)
			conn.Close()
			time.Sleep(time.Millisecond * 10)
			continue
		} else {
			core.InfoLogger.Println(string(welcome_message))
		}

		// Authenticate
		type AuthMessage struct {
			Action string `json:"action"`
			Key    string `json:"key"`
			Secret string `json:"secret"`
		}
		auth_message := AuthMessage{
			Action: "auth",
			Key:    viper.GetString("ALPACA_KEY"),
			Secret: viper.GetString("ALPACA_SECRET"),
		}
		conn.WriteJSON(auth_message)

		// Authentication response
		_, auth_response_message, err := conn.ReadMessage()
		if err != nil {
			core.WarnLogger.Println("Unable to read welcome message", err)
			conn.Close()
			time.Sleep(time.Millisecond * 10)
			continue
		} else {
			core.InfoLogger.Println(string(auth_response_message))
		}

		// Subscribe
		type SubscribeMessage struct {
			Action string   `json:"action"`
			Trades []string `json:"trades"`
		}
		subscribe_message := SubscribeMessage{
			Action: "subscribe",
			Trades: pairs,
		}
		conn.WriteJSON(subscribe_message)

		// Subscription response
		_, sub_response_message, err := conn.ReadMessage()
		if err != nil {
			core.WarnLogger.Println("Unable to read welcome message", err)
			conn.Close()
			time.Sleep(time.Millisecond * 10)
			continue
		} else {
			core.InfoLogger.Println(string(sub_response_message))
		}

		// Infinite loop reading trades
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				core.WarnLogger.Println("Could not read message", err)
				conn.Close()
				break
			}

			// Convert input message bytes to json
			var trades_in []AlpacaCrypto
			json.Unmarshal(message, &trades_in)
			trade_in := trades_in[0]
			pair := strings.Split(trade_in.Pair, "/")
			base := pair[0]
			quote := pair[1]

			// Convert input message to standard json schema
			trade_out := core.TradeJson{
				Source:      "alpaca-crypto",
				Base:        base,
				Quote:       quote,
				Exchange:    "alpaca",
				VolumeBase:  float64(trade_in.Size),
				VolumeQuote: trade_in.Price * float64(trade_in.Size),
				Price:       trade_in.Price,
				Timestamp:   trade_in.Timestamp,
			}

			// Push trade through the channel
			trades <- trade_out
		}
	}
}
