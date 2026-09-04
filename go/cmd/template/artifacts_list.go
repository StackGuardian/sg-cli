// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/templates"
	"github.com/spf13/cobra"
)

func newArtifactsListCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the artifacts of a template",
		Long:  `List the artifacts of a template`,
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
			response, err := c.Templates.ListAllTemplateArtifacts(ctx, ownerOrg, templateType, templateId, &templates.ListAllTemplateArtifactsRequest{SgOrgid: org, ArtifactPrefix: common.OptString(cmd, "prefix"), MaxKeys: common.OptInt(cmd, "max-keys"), StartAfterKey: common.OptString(cmd, "start-after-key")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("prefix", "", "Only list artifacts under this prefix.")
	cmd.Flags().Int("max-keys", 0, "Maximum number of keys.")
	cmd.Flags().String("start-after-key", "", "Pagination: start after this key.")
	return cmd
}
