// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflowrun

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newLogsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Get the logs of a stack workflow run",
		Long:  `Get the logs of a stack workflow run`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			workflowId := common.String(cmd, "workflow-id")
			runId := common.String(cmd, "run-id")
			response, err := c.StackWorkflowRuns.ReadStackWorkflowRunLogs(ctx, org, stackId, workflowId, workflowGroup, runId)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.PersistentFlags().String("run-id", "", "The workflow run ID.")
	_ = cmd.MarkPersistentFlagRequired("run-id")
	return cmd
}
