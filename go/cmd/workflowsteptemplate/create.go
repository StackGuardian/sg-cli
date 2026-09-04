// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowsteptemplate

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/workflowsteptemplate"
	"github.com/spf13/cobra"
)

func newCreateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create -- payload.json",
		Short: "Create a workflow step template",
		Long:  `Create a workflow step template`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			createFirstRevision := common.Bool(cmd, "create-first-revision")
			var req workflowsteptemplate.CreateWorkflowStepTemplate
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.WorkflowStepTemplate.CreateWorkflowStepTemplate(ctx, org, createFirstRevision, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Workflow step template created successfully.")
			return nil
		},
	}
	cmd.Flags().Bool("create-first-revision", true, "Also create the first revision of the template.")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
