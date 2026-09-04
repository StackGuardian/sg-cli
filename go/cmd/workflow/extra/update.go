// Code generated for sg-cli; edit the generator spec or this file directly.

package extra

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/workflows"
	"github.com/spf13/cobra"
)

func newUpdateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update -- payload.json",
		Short: "Update an existing workflow",
		Long:  `Update an existing workflow`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			workflowId := common.String(cmd, "workflow-id")
			var req workflows.PatchedWorkflow
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			var mode *workflows.UpgradeModeEnum
			if v := common.OptString(cmd, "upgrade-mode"); v != nil {
				if err := common.CheckEnum("upgrade-mode", *v, "PRESERVE_SETTINGS", "RESET_TO_TEMPLATE", "MANUAL"); err != nil {
					return err
				}
				m := workflows.UpgradeModeEnum(*v)
				mode = &m
			}
			response, err := c.Workflows.UpdateWorkflow(ctx, org, workflowId, workflowGroup, mode, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Workflow updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("workflow-id", "", "The workflow ID.")
	_ = cmd.MarkFlagRequired("workflow-id")
	cmd.Flags().String("upgrade-mode", "", "Template upgrade mode (PRESERVE_SETTINGS, RESET_TO_TEMPLATE, MANUAL).")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
