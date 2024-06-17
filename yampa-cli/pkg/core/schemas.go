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

package core

type TradeJson struct {
	Source      string  `json:"source"`
	Base        string  `json:"base"`
	Quote       string  `json:"quote"`
	Exchange    string  `json:"exchange"`
	VolumeBase  float64 `json:"volume_base"`
	VolumeQuote float64 `json:"volume_quote"`
	Price       float64 `json:"price"`
	Timestamp   string  `json:"timestamp"`
}

type TradeJsonKey struct {
	Source   string `json:"source"`
	Base     string `json:"base"`
	Quote    string `json:"quote"`
	Exchange string `json:"exchange"`
}

type TradeAvro struct {
	Source      string  `avro:"source"`
	Base        string  `avro:"base"`
	Quote       string  `avro:"quote"`
	Exchange    string  `avro:"exchange"`
	VolumeBase  float64 `avro:"volume_base"`
	VolumeQuote float64 `avro:"volume_quote"`
	Price       float64 `avro:"price"`
	Timestamp   string  `avro:"timestamp"`
}

type TradeAvroKey struct {
	Source   string `avro:"source"`
	Base     string `avro:"base"`
	Quote    string `avro:"quote"`
	Exchange string `avro:"exchange"`
}

var TradeAvroSchema = `{
	"type": "record",
	"name": "Trade",
	"fields" : [
		{"name": "source", "type": "string"},
		{"name": "base", "type": "string"},
		{"name": "quote", "type": "string"},
		{"name": "exchange", "type": "string"},
		{"name": "volume_base", "type": "double"},
		{"name": "volume_quote", "type": "double"},
		{"name": "price", "type": "double"},
		{"name": "timestamp", "type": "string"}
	]
}`

var TradeAvroKeySchema = `{
	"type": "record",
	"name": "Trade",
	"fields" : [
		{"name": "source", "type": "string"},
		{"name": "base", "type": "string"},
		{"name": "quote", "type": "string"},
		{"name": "exchange", "type": "string"}
	]
}`
