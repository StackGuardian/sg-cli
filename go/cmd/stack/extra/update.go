// Code generated for sg-cli; edit the generator spec or this file directly.

package extra

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
		Short: "Update a stack",
		Long:  `Update a stack`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			var req sggosdk.PatchedStack
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.Stacks.UpdateStack(ctx, org, stackId, workflowGroup, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Stack updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("stack-id", "", "The stack ID.")
	_ = cmd.MarkFlagRequired("stack-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
