// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowgroup

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get details of a workflow group",
		Long:  `Get details of a workflow group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			response, err := c.WorkflowGroups.ReadWorkflowGroup(ctx, org, workflowGroup)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("workflow-group", "", "The workflow group under the organization.")
	_ = cmd.MarkFlagRequired("workflow-group")
	return cmd
}
