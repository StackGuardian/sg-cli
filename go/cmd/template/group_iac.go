// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newGroupIacCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group-iac",
		Short: "Get an IaC template that belongs to an IaC group template",
		Long:  `Get an IaC template that belongs to an IaC group template`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			templateId := common.String(cmd, "template-id")
			iacTemplateId := common.String(cmd, "iac-template-id")
			ownerOrg := common.String(cmd, "owner-org")
			if ownerOrg == "" {
				ownerOrg = org
			}
			response, err := c.Templates.ReadIacGroupsIacTemplate(ctx, ownerOrg, iacTemplateId, templateId, &sggosdk.ReadIacGroupsIacTemplateRequest{SgOrgid: org})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("template-id", "", "The IaC group template.")
	_ = cmd.MarkFlagRequired("template-id")
	cmd.Flags().String("iac-template-id", "", "The IaC template inside the group.")
	_ = cmd.MarkFlagRequired("iac-template-id")
	cmd.Flags().String("owner-org", "", "Organization that owns the template (defaults to --org).")
	return cmd
}
