// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflowrun

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Cancel a stack workflow run through the DELETE endpoint",
		Long:  `Cancel a stack workflow run through the DELETE endpoint`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			workflowId := common.String(cmd, "workflow-id")
			runId := common.String(cmd, "run-id")
			response, err := c.StackWorkflowRuns.DeleteStackWorkflowRun(ctx, org, stackId, workflowId, workflowGroup, runId)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Stack workflow run cancelled.")
			return nil
		},
	}
	cmd.PersistentFlags().String("run-id", "", "The workflow run ID.")
	_ = cmd.MarkPersistentFlagRequired("run-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
