// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowgroup

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newCreateChildCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-child -- payload.json",
		Short: "Create a child workflow group inside a parent group",
		Long:  `Create a child workflow group inside a parent group`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			parent := common.String(cmd, "parent")
			var req sggosdk.WorkflowGroup
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.WorkflowGroups.CreateChildWorkflowGroup(ctx, org, parent, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Child workflow group created successfully.")
			return nil
		},
	}
	cmd.Flags().String("parent", "", "The parent workflow group.")
	_ = cmd.MarkFlagRequired("parent")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
