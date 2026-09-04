// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflow

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
		Short: "List the workflows of a stack",
		Long:  `List the workflows of a stack`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			response, err := c.StackWorkflows.ListAllStackWorkflows(ctx, org, stackId, workflowGroup, &sggosdk.ListAllStackWorkflowsRequest{Description: common.OptString(cmd, "description"), IacTemplateId: common.OptString(cmd, "iac-template-id"), LatestWfRunStatuses: common.OptInt(cmd, "latest-wf-run-statuses"), ResourceNames: common.OptString(cmd, "resource-names"), RunnerNames: common.OptString(cmd, "runner-names"), Tags: common.OptString(cmd, "tags"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key"), Limit: common.OptInt(cmd, "limit")})
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
	cmd.Flags().String("iac-template-id", "", "Filter by IaC template id.")
	cmd.Flags().Int("latest-wf-run-statuses", 0, "Filter by latest workflow run status.")
	cmd.Flags().String("resource-names", "", "Filter by resource names.")
	cmd.Flags().String("runner-names", "", "Filter by runner names.")
	cmd.Flags().String("tags", "", "Filter by tags.")
	return cmd
}
