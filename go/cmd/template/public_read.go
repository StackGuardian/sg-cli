// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newPublicReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get a public template",
		Long:  `Get a public template`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			templateId := common.String(cmd, "template-id")
			templateType := common.String(cmd, "template-type")
			ownerOrg := common.String(cmd, "owner-org")
			if err := common.CheckEnum("template-type", templateType, "IAC", "IAC_GROUP", "IAC_POLICY", "WORKFLOW_STEP"); err != nil {
				return err
			}
			response, err := c.Templates.ReadPublicTemplate(ctx, ownerOrg, templateType, templateId)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.PersistentFlags().String("template-id", "", "Template name, optionally with revision as NAME:REVISION.")
	_ = cmd.MarkPersistentFlagRequired("template-id")
	cmd.PersistentFlags().String("template-type", "", "Template type (IAC, IAC_GROUP, IAC_POLICY, WORKFLOW_STEP).")
	_ = cmd.MarkPersistentFlagRequired("template-type")
	cmd.Flags().String("owner-org", "", "Organization that owns the template.")
	_ = cmd.MarkFlagRequired("owner-org")
	return cmd
}
