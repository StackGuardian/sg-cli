// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/templates"
	"github.com/spf13/cobra"
)

func newArtifactsDownloadUrlCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-url",
		Short: "Get a presigned download URL for an artifact",
		Long:  `Get a presigned download URL for an artifact`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			templateId := common.String(cmd, "template-id")
			templateType := common.String(cmd, "template-type")
			ownerOrg := common.String(cmd, "owner-org")
			artifactId := common.String(cmd, "artifact-id")
			if err := common.CheckEnum("template-type", templateType, "IAC", "IAC_GROUP", "IAC_POLICY", "WORKFLOW_STEP"); err != nil {
				return err
			}
			if ownerOrg == "" {
				ownerOrg = org
			}
			response, err := c.Templates.GetTemplateArtifactDownloadUrl(ctx, ownerOrg, templateType, templateId, artifactId, &templates.TemplateArtifactRequest{SgOrgid: org})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("artifact-id", "", "The artifact ID.")
	_ = cmd.MarkFlagRequired("artifact-id")
	return cmd
}
