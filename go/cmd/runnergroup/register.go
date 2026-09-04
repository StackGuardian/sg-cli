// Code generated for sg-cli; edit the generator spec or this file directly.

package runnergroup

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newRegisterCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a runner node (use the runner group's own token as the API key)",
		Long:  `Register a runner node (use the runner group's own token as the API key)`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			runnerGroup := common.String(cmd, "runner-group")
			response, err := c.RunnerGroups.RegisterRunner(ctx, org, runnerGroup)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Runner registered.")
			return nil
		},
	}
	cmd.Flags().String("runner-group", "", "The runner group name.")
	_ = cmd.MarkFlagRequired("runner-group")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
