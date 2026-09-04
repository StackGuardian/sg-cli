// Code generated for sg-cli; edit the generator spec or this file directly.

package runnergroup

import (
	"context"
	"fmt"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newUpdateRunnerStateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-runner-state",
		Short: "Update the state of a runner (ACTIVE or DRAINING)",
		Long:  `Update the state of a runner (ACTIVE or DRAINING)`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			runnerGroup := common.String(cmd, "runner-group")
			status := common.String(cmd, "status")
			if status != "ACTIVE" && status != "DRAINING" {
				return fmt.Errorf("--status must be ACTIVE or DRAINING, got %q", status)
			}
			body := sggosdk.RunnerStatus{Status: sggosdk.StatusEnum(status)}
			if v := common.OptString(cmd, "runner-id"); v != nil {
				body.RunnerId = sggosdk.Optional(*v)
			}
			if v := common.OptStringSlice(cmd, "container-instance-ids"); v != nil {
				body.ContainerInstanceIds = sggosdk.Optional(v)
			}
			if err := c.RunnerGroups.UpdateRunnerState(ctx, org, runnerGroup, &body); err != nil {
				return err
			}
			common.Done(cmd, nil, "Runner state updated.")
			return nil
		},
	}
	cmd.Flags().String("runner-group", "", "The runner group name.")
	_ = cmd.MarkFlagRequired("runner-group")
	cmd.Flags().String("status", "", "New runner status (ACTIVE, DRAINING).")
	_ = cmd.MarkFlagRequired("status")
	cmd.Flags().String("runner-id", "", "Runner id.")
	cmd.Flags().StringSlice("container-instance-ids", nil, "Container instance ids (comma-separated).")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
