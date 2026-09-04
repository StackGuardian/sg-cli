// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newListCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List templates of a type",
		Long:  `List templates of a type`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			templateType := common.String(cmd, "template-type")
			if err := common.CheckEnum("template-type", templateType, "IAC", "IAC_GROUP", "IAC_POLICY", "WORKFLOW_STEP"); err != nil {
				return err
			}
			response, err := c.Templates.ListAllTemplates(ctx, sggosdk.ListAllTemplatesRequestTemplateType(templateType), &sggosdk.ListAllTemplatesRequest{SgOrgid: org, IsPublic: common.OptString(cmd, "is-public"), OwnerOrgs: common.OptString(cmd, "owner-orgs"), SearchQuery: common.OptString(cmd, "search-query"), TemplateId: common.OptString(cmd, "template-id"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("template-type", "", "Template type (IAC, IAC_GROUP, IAC_POLICY, WORKFLOW_STEP).")
	_ = cmd.MarkFlagRequired("template-type")
	cmd.Flags().String("is-public", "", "0 for private, 1 for public templates.")
	cmd.Flags().String("owner-orgs", "", "Organizations that own the templates (comma-separated).")
	cmd.Flags().String("search-query", "", "Search by name, tags or description.")
	cmd.Flags().String("template-id", "", "Parent template id to list its revisions.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	return cmd
}
