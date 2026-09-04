// Code generated for sg-cli; edit the generator spec or this file directly.

package stackrun

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newCreateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create -- payload.json",
		Short: "Trigger a stack run (use --action-type or a payload)",
		Long:  `Trigger a stack run (use --action-type or a payload)`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			actionType := common.String(cmd, "action-type")
			var req sggosdk.StackAction
			if len(args) > 0 {
				if err := common.ReadPayload(cmd, args, &req); err != nil {
					if errors.Is(err, common.ErrDryRun) {
						return nil
					}
					return err
				}
			}
			if actionType != "" {
				req.ActionType = actionType
			}
			if req.ActionType == "" {
				return errors.New("provide --action-type or a payload with ActionType")
			}
			response, err := c.StackRuns.CreateStackRun(ctx, org, stackId, workflowGroup, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Stack run created successfully.")
			return nil
		},
	}
	cmd.Flags().String("action-type", "", "Action to run, e.g. APPLY or DESTROY.")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
