// Code generated for sg-cli; edit the generator spec or this file directly.

package connector

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newListCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List connectors",
		Long:  `List connectors`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.Connectors.ListAllConnectors(ctx, org, &sggosdk.ListAllConnectorsRequest{ResourceNames: common.OptString(cmd, "resource-names"), Kind: common.OptString(cmd, "kind"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key"), Limit: common.OptInt(cmd, "limit")})
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
	cmd.Flags().String("kind", "", "Filter by kind (comma-separated), e.g. AWS_RBAC,AZURE_OIDC.")
	return cmd
}
