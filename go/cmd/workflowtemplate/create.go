// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowtemplate

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/workflowtemplates"
	"github.com/spf13/cobra"
)

func newCreateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create -- payload.json",
		Short: "Create a workflow template",
		Long:  `Create a workflow template`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			createFirstRevision := common.Bool(cmd, "create-first-revision")
			var req workflowtemplates.CreateWorkflowTemplateRequest
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.WorkflowTemplates.CreateWorkflowTemplate(ctx, org, createFirstRevision, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Workflow template created successfully.")
			return nil
		},
	}
	cmd.Flags().Bool("create-first-revision", true, "Also create the first revision of the template.")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
