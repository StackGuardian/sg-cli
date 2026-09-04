// Code generated for sg-cli; edit the generator spec or this file directly.

package stacktemplate

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/stacktemplates"
	"github.com/spf13/cobra"
)

func newUpdateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update -- payload.json",
		Short: "Update a stack template",
		Long:  `Update a stack template`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			templateId := common.String(cmd, "template-id")
			var req stacktemplates.UpdateStackTemplateRequest
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.StackTemplates.UpdateStackTemplate(ctx, org, templateId, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Stack template updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("template-id", "", "The template ID.")
	_ = cmd.MarkFlagRequired("template-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
