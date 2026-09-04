// Code generated for sg-cli; edit the generator spec or this file directly.

package runnergroup

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
		Short: "Update a runner group",
		Long:  `Update a runner group`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			runnerGroup := common.String(cmd, "runner-group")
			var req sggosdk.PatchedRunnerGroup
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.RunnerGroups.UpdateRunnerGroup(ctx, org, runnerGroup, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Runner group updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("runner-group", "", "The runner group name.")
	_ = cmd.MarkFlagRequired("runner-group")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
