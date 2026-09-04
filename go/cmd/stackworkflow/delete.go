// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflow

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a stack workflow",
		Long:  `Delete a stack workflow`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			workflowId := common.String(cmd, "workflow-id")
			if err := c.StackWorkflows.DeleteStackWorkflow(ctx, org, stackId, workflowId, workflowGroup); err != nil {
				return err
			}
			common.Done(cmd, nil, "Stack workflow deleted successfully.")
			return nil
		},
	}
	cmd.Flags().String("workflow-id", "", "The workflow ID.")
	_ = cmd.MarkFlagRequired("workflow-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
