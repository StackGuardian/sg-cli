// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewTemplateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage templates and template revisions across all template types",
		Long: `Manage templates and template revisions across all template types.

A template revision is addressed as NAME or NAME:REVISION (for example my-template:5). Without a revision the base template is used.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  list                   List templates of a type
  list-by-owner          List templates owned by an organization
  create                 Create a template or a new revision of it (payload must include TemplateType: IAC, IAC_GROUP or IAC_POLICY)
  read                   Get a template or one of its revisions
  update                 Update a template or one of its revisions
  delete                 Delete a template or one of its revisions
  group-iac              Get an IaC template that belongs to an IaC group template
  input-schema           Decode one of a template's input schemas (requires a user token)
  vcs-triggers           Create VCS webhook triggers for a template
  subscriptions          Manage the organization's template subscriptions
  artifacts              Work with template artifacts through presigned URLs
  public                 Browse public marketplace templates (requires a user token)`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newListCmd(c),
		newListByOwnerCmd(c),
		newCreateCmd(c),
		newReadCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newGroupIacCmd(c),
		newInputSchemaCmd(c),
		newVcsTriggersCmd(c),
		newSubscriptionsCmd(c),
		newArtifactsCmd(c),
		newPublicCmd(c),
	)
	return cmd
}
func newSubscriptionsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "Manage the organization's template subscriptions",
		Long:  `Manage the organization's template subscriptions`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  read                   Get the subscriptions of a type
  create                 Create a subscription set (all four subscription maps are required)
  update                 Subscribe to or unsubscribe from templates`)
		},
	}
	cmd.AddCommand(
		newSubscriptionsReadCmd(c),
		newSubscriptionsCreateCmd(c),
		newSubscriptionsUpdateCmd(c),
	)
	return cmd
}
func newArtifactsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artifacts",
		Short: "Work with template artifacts through presigned URLs",
		Long:  `Work with template artifacts through presigned URLs`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  list                   List the artifacts of a template
  download-url           Get a presigned download URL for an artifact
  upload-url             Get a presigned upload URL for an artifact
  delete                 Delete an artifact`)
		},
	}
	cmd.PersistentFlags().String("template-id", "", "Template name, optionally with revision as NAME:REVISION.")
	_ = cmd.MarkPersistentFlagRequired("template-id")
	cmd.PersistentFlags().String("template-type", "", "Template type (IAC, IAC_GROUP, IAC_POLICY, WORKFLOW_STEP).")
	_ = cmd.MarkPersistentFlagRequired("template-type")
	cmd.PersistentFlags().String("owner-org", "", "Organization that owns the template (defaults to --org).")
	cmd.AddCommand(
		newArtifactsListCmd(c),
		newArtifactsDownloadUrlCmd(c),
		newArtifactsUploadUrlCmd(c),
		newArtifactsDeleteCmd(c),
	)
	return cmd
}
func newPublicCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "public",
		Short: "Browse public marketplace templates (requires a user token)",
		Long:  `Browse public marketplace templates (requires a user token)`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  list                   List public templates of a type
  list-all               List public templates across all types
  read                   Get a public template`)
		},
	}
	cmd.AddCommand(
		newPublicListCmd(c),
		newPublicListAllCmd(c),
		newPublicReadCmd(c),
	)
	return cmd
}
