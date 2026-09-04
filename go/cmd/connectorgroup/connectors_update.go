// Code generated for sg-cli; edit the generator spec or this file directly.

package connectorgroup

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newConnectorsUpdateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update -- payload.json",
		Short: "Update a connector in the group",
		Long:  `Update a connector in the group`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			connectorGroup := common.String(cmd, "connector-group")
			connectorId := common.String(cmd, "connector-id")
			var req sggosdk.PatchedIntegration
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.ConnectorGroups.UpdateChildInACloudConnectorGroup(ctx, connectorId, connectorGroup, org, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Connector updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("connector-id", "", "The connector ID.")
	_ = cmd.MarkFlagRequired("connector-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
