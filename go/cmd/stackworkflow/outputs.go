// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflow

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newOutputsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outputs",
		Short: "Get the outputs of a stack workflow",
		Long:  `Get the outputs of a stack workflow`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			workflowId := common.String(cmd, "workflow-id")
			response, err := c.StackWorkflows.StackWorkflowOutputs(ctx, org, stackId, workflowId, workflowGroup)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("workflow-id", "", "The workflow ID.")
	_ = cmd.MarkFlagRequired("workflow-id")
	return cmd
}
