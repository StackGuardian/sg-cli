// Code generated for sg-cli; edit the generator spec or this file directly.

package extra

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
		Short: "List the stacks in a workflow group",
		Long:  `List the stacks in a workflow group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			response, err := c.Stacks.ListAllStacks(ctx, org, workflowGroup, &sggosdk.ListAllStacksRequest{Description: common.OptString(cmd, "description"), ResourceNames: common.OptString(cmd, "resource-names"), Tags: common.OptString(cmd, "tags"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key"), Limit: common.OptInt(cmd, "limit")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().Int("limit", 0, "Limit the number of results returned.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	cmd.Flags().String("description", "", "Filter by description.")
	cmd.Flags().String("resource-names", "", "Filter by resource names.")
	cmd.Flags().String("tags", "", "Filter by tags.")
	return cmd
}
