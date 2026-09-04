// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowrun

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newUpdateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update -- payload.json",
		Short: "Update a workflow run",
		Long:  `Update a workflow run`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			workflowId := common.String(cmd, "workflow-id")
			runId := common.String(cmd, "run-id")
			var req sggosdk.PatchedWorkflowRun
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			if err := common.Require(cmd, "workflow-group", "workflow-id"); err != nil {
				return err
			}
			response, err := c.WorkflowRuns.UpdateWorkflowRun(ctx, org, workflowId, workflowGroup, runId, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Workflow run updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("run-id", "", "The workflow run ID.")
	_ = cmd.MarkFlagRequired("run-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
