// Copyright © 2015 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package cmd

import (
	"fmt"
	"os"

	"github.com/bep/alfn/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/tylerb/graceful.v1"
	"log"
	"net/http"
	"time"
)

var cfgFile string
var app *lib.App

// flags
var serverPort int
var serverInterface string

var rootCmd = &cobra.Command{
	Use:   "alfn",
	Short: "Runs a web server and a feed aggregator with filters and a limiter.",
	Long:  `Runs a web server and a feed aggregator with filters and a limiter.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return startup()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.alfn/config.toml)")
	rootCmd.Flags().IntVarP(&serverPort, "port", "p", 1926, "port on which the server will listen")
	rootCmd.Flags().StringVarP(&serverInterface, "bind", "", "127.0.0.1", "interface to which the server will bind")
}

// Read in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.alfn")
	viper.AddConfigPath("/etc/alfn/")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("error: Failed to read config:", err)
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())
}

func startup() error {

	var config lib.Config

	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	serverAndPort := fmt.Sprintf("%s:%d", serverInterface, serverPort)

	if config.Feed.Link == "" {
		config.Feed.Link = "http://" + serverAndPort
	}

	if config.Feed.MaxItems <= 0 {
		config.Feed.MaxItems = 10
	}

	app := lib.NewApp(config).Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
		fmt.Fprintf(w, app.GetFeed())
	})

	fmt.Printf("\nStarting server on %s ...\n\n", serverAndPort)

	graceful.Run(serverAndPort, 10*time.Second, mux)

	fmt.Println("\nStopped ...")
	return nil
}
