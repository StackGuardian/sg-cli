// Code generated for sg-cli; edit the generator spec or this file directly.

package runnergroup

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newStorageBackendAuthCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage-backend-auth",
		Short: "Generate short-lived credentials for the runner group's storage backend",
		Long:  `Generate short-lived credentials for the runner group's storage backend`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			runnerGroup := common.String(cmd, "runner-group")
			response, err := c.RunnerGroups.GenerateStorageBackendAuth(ctx, org, runnerGroup)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("runner-group", "", "The runner group name.")
	_ = cmd.MarkFlagRequired("runner-group")
	return cmd
}
