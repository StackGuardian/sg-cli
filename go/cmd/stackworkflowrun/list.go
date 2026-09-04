// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflowrun

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/stackworkflowruns"
	"github.com/spf13/cobra"
)

func newListCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the runs of a stack workflow",
		Long:  `List the runs of a stack workflow`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			workflowId := common.String(cmd, "workflow-id")
			response, err := c.StackWorkflowRuns.ListAllStackWorkflowRuns(ctx, org, stackId, workflowId, workflowGroup, &stackworkflowruns.ListAllStackWorkflowRunsRequest{
				Limit: common.OptInt(cmd, "limit"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key"), CreatedAfter: common.OptString(cmd, "created-after"), ParentKSUID: common.OptString(cmd, "parent-ksuid"),
				ResourceNames: common.OptString(cmd, "resource-names"), TemplateId: common.OptString(cmd, "template-id"), UseMarketplaceTemplate: common.OptString(cmd, "use-marketplace-template"), Author: common.OptString(cmd, "author"),
				Approver: common.OptString(cmd, "approver"), LatestStatus: common.OptString(cmd, "latest-status"), RepoId: common.OptString(cmd, "repo-id"), Ref: common.OptString(cmd, "ref"), TerraformAction: common.OptString(cmd, "terraform-action"),
				RunnerType: common.OptString(cmd, "runner-type"), ContextTags: common.OptString(cmd, "context-tags"),
			})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().Int("limit", 0, "Limit the number of results returned.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	cmd.Flags().String("created-after", "", "Only runs created after this Unix timestamp in milliseconds.")
	cmd.Flags().String("parent-ksuid", "", "KSUID of the parent workflow.")
	cmd.Flags().String("resource-names", "", "Filter by resource names (comma-separated).")
	cmd.Flags().String("template-id", "", "Filter by IaC template id.")
	cmd.Flags().String("use-marketplace-template", "", "Filter by marketplace template usage (true/false).")
	cmd.Flags().String("author", "", "Filter by authors (comma-separated).")
	cmd.Flags().String("approver", "", "Filter by approvers (comma-separated).")
	cmd.Flags().String("latest-status", "", "Filter by latest status (comma-separated).")
	cmd.Flags().String("repo-id", "", "Filter by VCS repository.")
	cmd.Flags().String("ref", "", "Filter by VCS ref.")
	cmd.Flags().String("terraform-action", "", "Filter by Terraform action (comma-separated).")
	cmd.Flags().String("runner-type", "", "Filter by runner type.")
	cmd.Flags().String("context-tags", "", "Filter by context tags as key or key:value (comma-separated).")
	return cmd
}
