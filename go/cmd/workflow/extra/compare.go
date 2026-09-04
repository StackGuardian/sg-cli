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

func newCompareCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare -- payload.json",
		Short: "Compare a workflow against a template without applying changes",
		Long:  `Compare a workflow against a template without applying changes`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			workflowId := common.String(cmd, "workflow-id")
			var req workflows.CompareWorkflowRequest
			if len(args) > 0 {
				if err := common.ReadPayload(cmd, args, &req); err != nil {
					if errors.Is(err, common.ErrDryRun) {
						return nil
					}
					return err
				}
			}
			if v := common.OptString(cmd, "target-template-id"); v != nil {
				req.TargetTemplateId = *v
			}
			if v := common.OptString(cmd, "upgrade-mode"); v != nil {
				req.UpgradeMode = v
			}
			if req.TargetTemplateId == "" {
				return errors.New("provide --target-template-id or a payload with targetTemplateId")
			}
			response, err := c.Workflows.CompareWorkflow(ctx, org, workflowId, workflowGroup, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Comparison completed.")
			return nil
		},
	}
	cmd.Flags().String("workflow-id", "", "The workflow ID.")
	_ = cmd.MarkFlagRequired("workflow-id")
	cmd.Flags().String("target-template-id", "", "Template to compare against.")
	cmd.Flags().String("upgrade-mode", "", "Upgrade mode to simulate.")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
