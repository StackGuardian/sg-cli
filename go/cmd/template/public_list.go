// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/templates"
	"github.com/spf13/cobra"
)

func newPublicListCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List public templates of a type",
		Long:  `List public templates of a type`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			templateType := common.String(cmd, "template-type")
			if err := common.CheckEnum("template-type", templateType, "IAC", "IAC_GROUP", "IAC_POLICY", "WORKFLOW_STEP"); err != nil {
				return err
			}
			response, err := c.Templates.ListAllPublicTemplates(ctx, templateType, &templates.ListAllPublicTemplatesRequest{Limit: common.OptInt(cmd, "limit"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key"), TemplateId: common.OptString(cmd, "template-id"), OwnerOrgs: common.OptString(cmd, "owner-orgs"), SearchQuery: common.OptString(cmd, "search-query"), ContextTags: common.OptString(cmd, "context-tags"), SourceConfigKinds: common.OptString(cmd, "source-config-kinds")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.PersistentFlags().String("template-type", "", "Template type (IAC, IAC_GROUP, IAC_POLICY, WORKFLOW_STEP).")
	_ = cmd.MarkPersistentFlagRequired("template-type")
	cmd.Flags().Int("limit", 0, "Limit the number of results returned.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	cmd.Flags().String("template-id", "", "Parent template id (with owner org) to list its revisions.")
	cmd.Flags().String("owner-orgs", "", "Owner organizations (comma-separated).")
	cmd.Flags().String("search-query", "", "Search by name, description or tags.")
	cmd.Flags().String("context-tags", "", "Filter by context tags as key or key:value (comma-separated).")
	cmd.Flags().String("source-config-kinds", "", "Filter by source config kind (comma-separated).")
	return cmd
}
