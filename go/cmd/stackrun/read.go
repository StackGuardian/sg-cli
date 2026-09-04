// Code generated for sg-cli; edit the generator spec or this file directly.

package stackrun

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get a stack run",
		Long:  `Get a stack run`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			runId := common.String(cmd, "run-id")
			response, err := c.StackRuns.ReadStackRun(ctx, org, stackId, runId, workflowGroup)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("run-id", "", "The stack run ID.")
	_ = cmd.MarkFlagRequired("run-id")
	return cmd
}
