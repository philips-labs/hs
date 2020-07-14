/*
Copyright © 2020 Andy Lo-A-Foe <andy.lo-a-foe@philips.com>

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

	"github.com/philips-software/go-hsdp-api/iron"

	"github.com/spf13/cobra"
)

// ironQueueCmd represents the queue command
var ironQueueCmd = &cobra.Command{
	Use:     "queue <code>",
	Aliases: []string{"q"},
	Short:   "Queues tasks on a cluster",
	Long:    `Queues tasks on a cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
		config, err := readIronConfig()
		if err != nil {
			fmt.Printf("error reading iron config: %v\n", err)
			return
		}
		client, err := iron.NewClient(config)
		if err != nil {
			fmt.Printf("error configuring iron client: %v\n", err)
			return
		}
		client.Close()
	},
}

func init() {
	ironCmd.AddCommand(ironQueueCmd)
	ironQueueCmd.Flags().StringP("payload", "p", "", "Payload to use")
	ironQueueCmd.Flags().IntP("timeout", "t", 3600, "Timeout to use in seconds")
}
