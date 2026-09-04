// Code generated for sg-cli; edit the generator spec or this file directly.

package connectorgroup

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newConnectorsListCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the connectors in the group",
		Long:  `List the connectors in the group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			connectorGroup := common.String(cmd, "connector-group")
			response, err := c.ConnectorGroups.ListAllConnectorsInAGroup(ctx, connectorGroup, org, &sggosdk.ListAllConnectorsInAGroupRequest{ResourceNames: common.OptString(cmd, "resource-names"), Kind: common.OptString(cmd, "kind"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key"), Limit: common.OptInt(cmd, "limit"), Scope: common.OptString(cmd, "scope")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().Int("limit", 0, "Limit the number of results returned.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	cmd.Flags().String("resource-names", "", "Filter by resource name.")
	cmd.Flags().String("kind", "", "Filter by kind.")
	cmd.Flags().String("scope", "", "Filter by scope.")
	return cmd
}
