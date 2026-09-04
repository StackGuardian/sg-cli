// Code generated for sg-cli; edit the generator spec or this file directly.

package statebackend

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/statebackends"
	"github.com/spf13/cobra"
)

func newListStatefilesCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-statefiles",
		Short: "List the state files in a storage backend",
		Long:  `List the state files in a storage backend`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			config := common.String(cmd, "config")
			response, err := c.StateBackends.ListStatefiles(ctx, org, &statebackends.ListStatefilesRequest{Config: config})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("config", "", "Base64-encoded JSON StateBackendConfig.")
	_ = cmd.MarkFlagRequired("config")
	return cmd
}
