// Code generated for sg-cli; edit the generator spec or this file directly.

package runnergroup

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get details of a runner group",
		Long:  `Get details of a runner group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			runnerGroup := common.String(cmd, "runner-group")
			response, err := c.RunnerGroups.ReadRunnerGroup(ctx, org, runnerGroup, &sggosdk.ReadRunnerGroupRequest{GetActiveWorkflows: common.OptBool(cmd, "get-active-workflows"), GetActiveWorkflowsDetails: common.OptBool(cmd, "get-active-workflows-details")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("runner-group", "", "The runner group name.")
	_ = cmd.MarkFlagRequired("runner-group")
	cmd.Flags().Bool("get-active-workflows", false, "Include the active workflows.")
	cmd.Flags().Bool("get-active-workflows-details", false, "Include details of the active workflows.")
	return cmd
}
