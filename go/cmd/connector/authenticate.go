// Code generated for sg-cli; edit the generator spec or this file directly.

package connector

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newAuthenticateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "authenticate",
		Short: "Verify that a connector can authenticate against its provider",
		Long:  `Verify that a connector can authenticate against its provider`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			connectorId := common.String(cmd, "connector-id")
			response, err := c.Connectors.AuthenticateConnector(ctx, connectorId, org)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("connector-id", "", "The connector ID.")
	_ = cmd.MarkFlagRequired("connector-id")
	return cmd
}
