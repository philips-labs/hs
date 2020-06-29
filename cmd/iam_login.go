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
	"context"
	"fmt"
	"github.com/philips-software/go-hsdp-api/iam"
	"net/http"
	"os"
	"time"

	"github.com/pkg/browser"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var iamLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// IAM
		iamClient, err := iam.NewClient(http.DefaultClient, &iam.Config{
			IAMURL: "https://iam-client-test.us-east.philips-healthsuite.com",
			IDMURL: "https://idm-client-test.us-east.philips-healthsuite.com",
			OAuth2ClientID: clientID,
			OAuth2Secret: clientSecret,
			Debug: true,
			DebugLog: "/tmp/hs.log",
		})
		if err != nil {
			fmt.Printf("error initializing IAM client: %v\n", err)
			os.Exit(1)
		}
		e := echo.New()
		e.HideBanner = true
		redirectURI := "http://localhost:35444/callback"
		loginSuccess := false
		e.GET("/callback", func(c echo.Context) error {
			code := c.QueryParam("code")
			err := iamClient.CodeLogin(code, redirectURI)
			if err != nil {
				c.HTML(http.StatusForbidden, "<html><body>Login failed</body></html>")
				go func() {
					time.Sleep(1 * time.Second)
					_ = e.Shutdown(context.Background())
				}()
				return err
			}
			c.HTML(http.StatusOK, "<html><body>You are now logged in! Feel free to close this window...</body></html>")
			loginSuccess = true
			go func() {
				time.Sleep(2 * time.Second)
				_ = e.Shutdown(context.Background())
			}()
			return nil
		})
		fmt.Printf("login using your browser ...\n")
		err = browser.OpenURL("https://iam-client-test.us-east.philips-healthsuite.com/authorize/oauth2/authorize?response_type=code&client_id=hsappclient&redirect_uri=http://localhost:35444/callback")
		if err != nil {
			fmt.Printf("failed to open browser login: %v\n", err)
			os.Exit(1)
		}
		_ = e.Start(":35444")
		if !loginSuccess {
			fmt.Printf("login failed. Please try again ...\n")
			os.Exit(1)
		}

		introspect, _, err := iamClient.Introspect()
		if err != nil {
			fmt.Printf("error performing introspect: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("logged in as: %s\n", introspect.Username)
	},
}

func init() {
	iamCmd.AddCommand(iamLoginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}