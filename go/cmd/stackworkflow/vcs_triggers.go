// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflow

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/workflows"
	"github.com/spf13/cobra"
)

func newVcsTriggersCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vcs-triggers -- payload.json",
		Short: "Create VCS webhook triggers for a stack workflow",
		Long:  `Create VCS webhook triggers for a stack workflow`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			workflowId := common.String(cmd, "workflow-id")
			var req workflows.CreateVcsTriggersRequest
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.StackWorkflows.CreateVcsTriggersForStackWorkflow(ctx, org, stackId, workflowId, workflowGroup, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "VCS triggers created successfully.")
			return nil
		},
	}
	cmd.Flags().String("workflow-id", "", "The workflow ID.")
	_ = cmd.MarkFlagRequired("workflow-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
