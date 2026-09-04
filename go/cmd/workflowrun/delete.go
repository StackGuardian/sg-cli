// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowrun

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Cancel a workflow run through the DELETE endpoint",
		Long:  `Cancel a workflow run through the DELETE endpoint`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			workflowId := common.String(cmd, "workflow-id")
			runId := common.String(cmd, "run-id")
			if err := common.Require(cmd, "workflow-group", "workflow-id"); err != nil {
				return err
			}
			response, err := c.WorkflowRuns.DeleteWorkflowRun(ctx, org, workflowId, workflowGroup, runId)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Workflow run cancelled.")
			return nil
		},
	}
	cmd.Flags().String("run-id", "", "The workflow run ID.")
	_ = cmd.MarkFlagRequired("run-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
