// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowsteptemplate

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a workflow step template",
		Long:  `Delete a workflow step template`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			templateId := common.String(cmd, "template-id")
			if err := c.WorkflowStepTemplate.DeleteWorkflowStepTemplate(ctx, org, templateId); err != nil {
				return err
			}
			common.Done(cmd, nil, "Workflow step template deleted successfully.")
			return nil
		},
	}
	cmd.Flags().String("template-id", "", "The template ID.")
	_ = cmd.MarkFlagRequired("template-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
