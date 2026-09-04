// Code generated for sg-cli; edit the generator spec or this file directly.

package statebackend

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/statebackends"
	"github.com/spf13/cobra"
)

func newUpdateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update -- payload.json",
		Short: "Update a state backend",
		Long:  `Update a state backend`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			stateBackend := common.String(cmd, "state-backend")
			var req statebackends.PatchedStateBackend
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.StateBackends.UpdateStateBackend(ctx, org, stateBackend, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "State backend updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("state-backend", "", "The state backend name.")
	_ = cmd.MarkFlagRequired("state-backend")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
