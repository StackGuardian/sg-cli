// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowtemplate

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/workflowtemplaterevisions"
	"github.com/spf13/cobra"
)

func newRevisionsCreateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create -- payload.json",
		Short: "Create a new revision of a workflow template",
		Long:  `Create a new revision of a workflow template`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			templateId := common.String(cmd, "template-id")
			var req workflowtemplaterevisions.CreateWorkflowTemplateRevisionsRequest
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.WorkflowTemplatesRevisions.CreateWorkflowTemplateRevision(ctx, org, templateId, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Revision created successfully.")
			return nil
		},
	}
	cmd.Flags().String("template-id", "", "The template ID.")
	_ = cmd.MarkFlagRequired("template-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
