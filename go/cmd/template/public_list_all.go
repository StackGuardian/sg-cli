// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/templates"
	"github.com/spf13/cobra"
)

func newPublicListAllCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-all",
		Short: "List public templates across all types",
		Long:  `List public templates across all types`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			response, err := c.Templates.ListAllPublicTemplatesV2(ctx, &templates.ListAllPublicTemplatesV2Request{Limit: common.OptInt(cmd, "limit"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key"), IsActive: common.OptString(cmd, "is-active"), TemplateName: common.OptString(cmd, "template-name"), GitHubComRepoID: common.OptString(cmd, "github-repo-id"), SearchQuery: common.OptString(cmd, "search-query"), Tags: common.OptString(cmd, "tags"), TemplateTypes: common.OptString(cmd, "template-types"), SourceConfigKinds: common.OptString(cmd, "source-config-kinds"), ContextTags: common.OptString(cmd, "context-tags"), OwnerOrgs: common.OptString(cmd, "owner-orgs")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().Int("limit", 0, "Limit the number of results returned.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	cmd.Flags().String("is-active", "", "Filter by active state (0 or 1).")
	cmd.Flags().String("template-name", "", "Filter by template name.")
	cmd.Flags().String("github-repo-id", "", "Filter by GitHub repository id.")
	cmd.Flags().String("search-query", "", "Search by name, description, tags or id.")
	cmd.Flags().String("tags", "", "Filter by tags (comma-separated).")
	cmd.Flags().String("template-types", "", "Filter by template types (comma-separated).")
	cmd.Flags().String("source-config-kinds", "", "Filter by source config kind (comma-separated).")
	cmd.Flags().String("context-tags", "", "Filter by context tags as key or key:value (comma-separated).")
	cmd.Flags().String("owner-orgs", "", "Owner organizations (comma-separated).")
	return cmd
}
