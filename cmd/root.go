// Copyright © 2016 Jose Carlos Ustra Junior <dev@ustrajunior.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ustrajunior/keepup/cfgo"
)

var (
	api       *cloudflare.API
	cfgFile   string
	account   string
	ip        string
	dnsRecord string

	cfClient *cfgo.CloudflareClient
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "keepup",
	Short: "Update your cloudflare Records",
	Long: `A command line to update your cloudflare.com dns records
	with your current ip or with other given ip`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.keepup.yaml)")
	RootCmd.PersistentFlags().StringVar(&account, "account", "", "the account to work on")
	RootCmd.PersistentFlags().StringVar(&dnsRecord, "dns", "", "the dns record to be updated (domain.com / my.domain.com)")
	RootCmd.PersistentFlags().StringVar(&ip, "ip", "", "the ip that will be used to update the dns record (127.0.0.1)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$HOME")
		viper.SetConfigName(".keepup")
	}
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("could not read the config file: %v", err)
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())

	if account == "" {
		if viper.GetString("default") == "" {
			log.Fatal("you must define a --account flag or set the default option on the config file")
		}
		account = viper.GetString("default")
	}

	if dnsRecord == "" {
		log.Fatal("you must define a --dns flag")
	}

	var err error

	cfClient, err = cfgo.NewCloudflareClient(getKey("cfKey"), getKey("cfEmail"))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func getKey(k string) string {
	return viper.GetString(fmt.Sprintf("%s.%s", account, k))
}
