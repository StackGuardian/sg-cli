// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a template or one of its revisions",
		Long:  `Delete a template or one of its revisions`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			templateId := common.String(cmd, "template-id")
			templateType := common.String(cmd, "template-type")
			ownerOrg := common.String(cmd, "owner-org")
			if err := common.CheckEnum("template-type", templateType, "IAC", "IAC_GROUP", "IAC_POLICY", "WORKFLOW_STEP"); err != nil {
				return err
			}
			if ownerOrg == "" {
				ownerOrg = org
			}
			if err := c.Templates.DeleteTemplateRevision(ctx, ownerOrg, templateId, templateType, &sggosdk.DeleteTemplateRevisionRequest{SgOrgid: org}); err != nil {
				return err
			}
			common.Done(cmd, nil, "Template deleted successfully.")
			return nil
		},
	}
	cmd.Flags().String("template-id", "", "Template name, optionally with revision as NAME:REVISION.")
	_ = cmd.MarkFlagRequired("template-id")
	cmd.Flags().String("template-type", "", "Template type (IAC, IAC_GROUP, IAC_POLICY, WORKFLOW_STEP).")
	_ = cmd.MarkFlagRequired("template-type")
	cmd.Flags().String("owner-org", "", "Organization that owns the template (defaults to --org).")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
