// Code generated for sg-cli; edit the generator spec or this file directly.

package connectorgroup

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a connector group",
		Long:  `Delete a connector group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			connectorGroup := common.String(cmd, "connector-group")
			response, err := c.ConnectorGroups.DeleteConnectorGroup(ctx, connectorGroup, org)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Connector group deleted successfully.")
			return nil
		},
	}
	cmd.Flags().String("connector-group", "", "The connector group name.")
	_ = cmd.MarkFlagRequired("connector-group")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
