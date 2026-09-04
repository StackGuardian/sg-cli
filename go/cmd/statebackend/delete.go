// Code generated for sg-cli; edit the generator spec or this file directly.

package statebackend

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a state backend",
		Long:  `Delete a state backend`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			stateBackend := common.String(cmd, "state-backend")
			response, err := c.StateBackends.DeleteStateBackend(ctx, org, stateBackend)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "State backend deleted successfully.")
			return nil
		},
	}
	cmd.Flags().String("state-backend", "", "The state backend name.")
	_ = cmd.MarkFlagRequired("state-backend")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
