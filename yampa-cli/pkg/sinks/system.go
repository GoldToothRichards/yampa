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
	"os"

	"yampa-cli/pkg/core"
)

func WriteToStdout(
	trades chan core.TradeJson,
) {
	// Infinite loop writing trades from the channel to stdout
	for {
		trade := <-trades
		json_bytes, _ := json.Marshal(trade)
		json_str := string(json_bytes)
		core.TradeLogger.Println(json_str)
	}
}

func WriteToFile(
	trades chan core.TradeJson,
	file string,
) {
	f, _ := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()

	// Infinite loop writing trades from the channel to a local log file
	for {
		trade := <-trades
		json_bytes, _ := json.Marshal(trade)
		json_str := string(json_bytes)
		f.WriteString(json_str + "\n")
	}
}
