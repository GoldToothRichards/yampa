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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yampa",
	Short: "A CLI utility for streaming real-time trading data",
	Long:  `A CLI utility for streaming real-time trading data`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.yampa.env)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".yampa" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("env")
		viper.SetConfigName(".yampa")
	}

	viper.SetDefault("KAFKA_BROKER_URLS", "localhost:9092")
	viper.SetDefault("SCHEMA_REGISTRY_URL", "localhost:8081")
	viper.SetDefault("ALPACA_KEY", "")
	viper.SetDefault("ALPACA_SECRET", "")
	viper.SetDefault("KAFKA_MESSAGE_FORMAT", "JSON")
	viper.SetDefault("KAFKA_MESSAGE_INCLUDE_KEY", false)
	viper.SetDefault("PUSHER_APP_ID", "")
	viper.SetDefault("PUSHER_KEY", "")
	viper.SetDefault("PUSHER_SECRET", "")
	viper.SetDefault("PUSHER_CLUSTER", "us2")
	viper.SetDefault("PUSHER_CHANNEL", "trades")
	viper.SetDefault("PUSHER_EVENT_NAME", "new-trade")
	viper.SetDefault("PUSHER_BATCH_SIZE", 500)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

}
