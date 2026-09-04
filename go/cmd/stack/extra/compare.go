// Code generated for sg-cli; edit the generator spec or this file directly.

package extra

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/stacks"
	"github.com/spf13/cobra"
)

func newCompareCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare -- payload.json",
		Short: "Compare a stack against a stack template without applying changes",
		Long:  `Compare a stack against a stack template without applying changes`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			var req stacks.CompareStackRequest
			if len(args) > 0 {
				if err := common.ReadPayload(cmd, args, &req); err != nil {
					if errors.Is(err, common.ErrDryRun) {
						return nil
					}
					return err
				}
			}
			if v := common.OptString(cmd, "target-template-group-id"); v != nil {
				req.TargetTemplateGroupId = *v
			}
			if v := common.OptString(cmd, "upgrade-mode"); v != nil {
				req.UpgradeMode = v
			}
			if req.TargetTemplateGroupId == "" {
				return errors.New("provide --target-template-group-id or a payload with targetTemplateGroupId")
			}
			response, err := c.Stacks.CompareStack(ctx, org, stackId, workflowGroup, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Comparison completed.")
			return nil
		},
	}
	cmd.Flags().String("stack-id", "", "The stack ID.")
	_ = cmd.MarkFlagRequired("stack-id")
	cmd.Flags().String("target-template-group-id", "", "Stack template to compare against.")
	cmd.Flags().String("upgrade-mode", "", "Upgrade mode to simulate.")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
