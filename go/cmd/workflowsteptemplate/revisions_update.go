// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowsteptemplate

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/workflowsteptemplaterevision"
	"github.com/spf13/cobra"
)

func newRevisionsUpdateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update -- payload.json",
		Short: "Update a workflow step template revision",
		Long:  `Update a workflow step template revision`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			revisionId := common.String(cmd, "revision-id")
			var req workflowsteptemplaterevision.UpdateWorkflowStepTemplateRevisionModel
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.WorkflowStepTemplateRevision.UpdateWorkflowStepTemplateRevision(ctx, org, revisionId, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Revision updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("revision-id", "", "The revision ID (NAME:REVISION).")
	_ = cmd.MarkFlagRequired("revision-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
