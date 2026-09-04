// Code generated for sg-cli; edit the generator spec or this file directly.

package connectorgroup

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newConnectorsAuthenticateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "authenticate",
		Short: "Verify that a connector in the group can authenticate",
		Long:  `Verify that a connector in the group can authenticate`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			connectorGroup := common.String(cmd, "connector-group")
			connectorId := common.String(cmd, "connector-id")
			response, err := c.ConnectorGroups.AuthenticateChildInACloudConnectorGroup(ctx, connectorId, connectorGroup, org)
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
