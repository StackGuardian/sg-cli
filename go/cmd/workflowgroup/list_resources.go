// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowgroup

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/workflowgroups"
	"github.com/spf13/cobra"
)

func newListResourcesCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-resources",
		Short: "List the workflows, stacks and child groups inside a workflow group",
		Long:  `List the workflows, stacks and child groups inside a workflow group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			response, err := c.WorkflowGroups.ListAllResourcesInWorkflowGroup(ctx, org, workflowGroup, &workflowgroups.ListAllResourcesRequest{Limit: common.OptInt(cmd, "limit"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key"), ResourceNames: common.OptString(cmd, "resource-names"), Description: common.OptString(cmd, "description"), Tags: common.OptString(cmd, "tags"), SearchQuery: common.OptString(cmd, "search-query"), LatestStatuses: common.OptString(cmd, "latest-statuses"), ResourceTypes: common.OptString(cmd, "resource-types"), ContextTags: common.OptString(cmd, "context-tags"), IsActive: common.OptString(cmd, "is-active")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("workflow-group", "", "The workflow group under the organization.")
	_ = cmd.MarkFlagRequired("workflow-group")
	cmd.Flags().Int("limit", 0, "Limit the number of results returned.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	cmd.Flags().String("resource-names", "", "Filter by resource names (comma-separated).")
	cmd.Flags().String("description", "", "Filter by description.")
	cmd.Flags().String("tags", "", "Filter by tags (comma-separated).")
	cmd.Flags().String("search-query", "", "Search by resource name, description or tags.")
	cmd.Flags().String("latest-statuses", "", "Filter by latest status (comma-separated).")
	cmd.Flags().String("resource-types", "", "Filter by resource type: WORKFLOW, STACK, WORKFLOW_GROUP (comma-separated).")
	cmd.Flags().String("context-tags", "", "Filter by context tags as key or key:value (comma-separated).")
	cmd.Flags().String("is-active", "", "Filter by active state (0 or 1).")
	return cmd
}
