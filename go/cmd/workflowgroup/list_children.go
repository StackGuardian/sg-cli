// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowgroup

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newListChildrenCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-children",
		Short: "List the child workflow groups of a parent group",
		Long:  `List the child workflow groups of a parent group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			parent := common.String(cmd, "parent")
			response, err := c.WorkflowGroups.ListAllChildWorkflowGroups(ctx, org, parent, &sggosdk.ListAllChildWorkflowGroupsRequest{Description: common.OptString(cmd, "description"), ResourceNames: common.OptString(cmd, "resource-names"), Tags: common.OptString(cmd, "tags"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key"), Limit: common.OptInt(cmd, "limit")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("parent", "", "The parent workflow group.")
	_ = cmd.MarkFlagRequired("parent")
	cmd.Flags().Int("limit", 0, "Limit the number of results returned.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	cmd.Flags().String("description", "", "Filter by description.")
	cmd.Flags().String("resource-names", "", "Filter by resource names.")
	cmd.Flags().String("tags", "", "Filter by tags.")
	return cmd
}
