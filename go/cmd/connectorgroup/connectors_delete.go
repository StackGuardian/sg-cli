// Code generated for sg-cli; edit the generator spec or this file directly.

package connectorgroup

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newConnectorsDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a connector from the group",
		Long:  `Delete a connector from the group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			connectorGroup := common.String(cmd, "connector-group")
			connectorId := common.String(cmd, "connector-id")
			response, err := c.ConnectorGroups.DeleteChildInACloudConnectorGroup(ctx, connectorId, connectorGroup, org)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Connector deleted successfully.")
			return nil
		},
	}
	cmd.Flags().String("connector-id", "", "The connector ID.")
	_ = cmd.MarkFlagRequired("connector-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
