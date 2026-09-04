// Code generated for sg-cli; edit the generator spec or this file directly.

package statebackend

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get details of a state backend",
		Long:  `Get details of a state backend`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			stateBackend := common.String(cmd, "state-backend")
			response, err := c.StateBackends.ReadStateBackend(ctx, org, stateBackend)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("state-backend", "", "The state backend name.")
	_ = cmd.MarkFlagRequired("state-backend")
	return cmd
}
