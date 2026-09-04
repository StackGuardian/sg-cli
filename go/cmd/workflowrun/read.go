// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowrun

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get a workflow run",
		Long:  `Get a workflow run`,
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
			response, err := c.WorkflowRuns.ReadWorkflowRun(ctx, org, workflowId, workflowGroup, runId)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("run-id", "", "The workflow run ID.")
	_ = cmd.MarkFlagRequired("run-id")
	return cmd
}
