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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ustrajunior/cfgo"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the dns record with current ip",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var currentIP string
		client, err := cfgo.NewClient(viper.GetString("cfKey"), viper.GetString("cfEmail"))
		if err != nil {
			log.Fatal(err)
		}

		if len(ip) > 7 {
			currentIP = ip
		} else {
			currentIP, err = cfgo.GetIPV4IP()
			if err != nil {
				log.Fatal(err)
			}
		}
		record, err := client.GetDNSRecord(zone, dnsRecord)
		if err != nil {
			log.Fatal(err)
		}

		if record.Content != currentIP {
			record.Content = strings.TrimSpace(currentIP)
			client.UpdateDNSRecord(record)
		}
		fmt.Println("DNS updated to:", record.Name, record.Content)
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
