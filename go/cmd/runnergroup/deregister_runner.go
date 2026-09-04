// Code generated for sg-cli; edit the generator spec or this file directly.

package runnergroup

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeregisterRunnerCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deregister-runner",
		Short: "Deregister a runner from a runner group",
		Long:  `Deregister a runner from a runner group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			runnerGroup := common.String(cmd, "runner-group")
			var body sggosdk.RunnerDeregister
			if v := common.OptString(cmd, "runner-id"); v != nil {
				body.RunnerId = sggosdk.Optional(*v)
			}
			if v := common.OptStringSlice(cmd, "container-instance-ids"); v != nil {
				body.ContainerInstanceIds = sggosdk.Optional(v)
			}
			if v := common.OptBool(cmd, "force"); v != nil {
				body.ForceDeregister = sggosdk.Optional(*v)
			}
			if err := c.RunnerGroups.DeregisterRunner(ctx, org, runnerGroup, &body); err != nil {
				return err
			}
			common.Done(cmd, nil, "Runner deregistered.")
			return nil
		},
	}
	cmd.Flags().String("runner-group", "", "The runner group name.")
	_ = cmd.MarkFlagRequired("runner-group")
	cmd.Flags().String("runner-id", "", "Runner id.")
	cmd.Flags().StringSlice("container-instance-ids", nil, "Container instance ids (comma-separated).")
	cmd.Flags().Bool("force", false, "Force deregistration.")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
