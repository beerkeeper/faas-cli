// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/openfaas/faas-cli/proxy"
	"github.com/openfaas/faas-cli/schema"
	"github.com/spf13/cobra"
)

// secretListCmd represents the secretCreate command
var secretListCmd = &cobra.Command{
	Use:     `list [--tls-no-verify]`,
	Aliases: []string{"ls"},
	Short:   "List all secrets",
	Long:    `List all secrets`,
	Example: `faas-cli secret list
faas-cli secret list --gateway=http://127.0.0.1:8080`,
	RunE:    runSecretList,
	PreRunE: preRunSecretListCmd,
}

func init() {
	secretListCmd.Flags().StringVarP(&gateway, "gateway", "g", defaultGateway, "Gateway URL starting with http(s)://")
	secretListCmd.Flags().BoolVar(&tlsInsecure, "tls-no-verify", false, "Disable TLS validation")

	secretCmd.AddCommand(secretListCmd)
}

func preRunSecretListCmd(cmd *cobra.Command, args []string) error {
	return nil
}

func runSecretList(cmd *cobra.Command, args []string) error {
	var gatewayAddress string
	gatewayAddress = getGatewayURL(gateway, defaultGateway, "", os.Getenv(openFaaSURLEnvironment))

	if msg := checkTLSInsecure(gatewayAddress, tlsInsecure); len(msg) > 0 {
		fmt.Println(msg)
	}

	secrets, err := proxy.GetSecretList(gatewayAddress, tlsInsecure)
	if err != nil {
		return err
	}

	if len(secrets) == 0 {
		fmt.Printf("No secrets found.\n")
		return nil
	}

	fmt.Printf("%s", renderSecretList(secrets))

	return nil
}

func renderSecretList(secrets []schema.Secret) string {
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "NAME")

	for _, secret := range secrets {
		fmt.Fprintf(w, "%s\n", secret.Name)
	}

	fmt.Fprintln(w)
	w.Flush()
	return b.String()
}
