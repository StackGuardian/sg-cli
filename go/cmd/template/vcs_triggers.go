// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/workflows"
	"github.com/spf13/cobra"
)

func newVcsTriggersCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vcs-triggers -- payload.json",
		Short: "Create VCS webhook triggers for a template",
		Long:  `Create VCS webhook triggers for a template`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			templateId := common.String(cmd, "template-id")
			templateType := common.String(cmd, "template-type")
			ownerOrg := common.String(cmd, "owner-org")
			var req workflows.CreateVcsTriggersRequest
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			if err := common.CheckEnum("template-type", templateType, "IAC", "IAC_GROUP", "IAC_POLICY", "WORKFLOW_STEP"); err != nil {
				return err
			}
			if ownerOrg == "" {
				ownerOrg = org
			}
			response, err := c.Templates.CreateTemplateVcsTriggers(ctx, ownerOrg, templateType, templateId, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "VCS triggers created successfully.")
			return nil
		},
	}
	cmd.Flags().String("template-id", "", "Template name, optionally with revision as NAME:REVISION.")
	_ = cmd.MarkFlagRequired("template-id")
	cmd.Flags().String("template-type", "", "Template type (IAC, IAC_GROUP, IAC_POLICY, WORKFLOW_STEP).")
	_ = cmd.MarkFlagRequired("template-type")
	cmd.Flags().String("owner-org", "", "Organization that owns the template (defaults to --org).")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
