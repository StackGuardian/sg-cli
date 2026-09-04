// Code generated for sg-cli; edit the generator spec or this file directly.

package organization

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/organizations"
	"github.com/spf13/cobra"
)

func newWorkflowsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflows",
		Short: "List workflows across every workflow group of an organization",
		Long:  `List workflows across every workflow group of an organization`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.Organizations.ListAllWorkflows(ctx, org, &organizations.ListAllOrganizationWorkflowsRequest{
				Limit: common.OptInt(cmd, "limit"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key"), FetchSgOwnedWfs: common.OptBool(cmd, "fetch-sg-owned-wfs"),
				LatestWfRunStatuses: common.OptString(cmd, "latest-wf-run-statuses"), LatestDriftRunStatuses: common.OptString(cmd, "latest-drift-run-statuses"),
				RepoId: common.OptString(cmd, "repo-id"), UseMarketplaceTemplate: common.OptString(cmd, "use-marketplace-template"), IacTemplateId: common.OptString(cmd, "iac-template-id"),
				ParentIacTemplateId: common.OptString(cmd, "parent-iac-template-id"), TerraformConfigTerraformVersion: common.OptString(cmd, "terraform-version"),
				RunnerNames: common.OptString(cmd, "runner-names"), TerraformConfigDriftCheck: common.OptString(cmd, "drift-check"), ResourceNames: common.OptString(cmd, "resource-names"),
				Description: common.OptString(cmd, "description"), Tags: common.OptString(cmd, "tags"), IsActive: common.OptString(cmd, "is-active"),
			})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("org", "", "The organization name on Stackguardian platform.")
	_ = cmd.MarkFlagRequired("org")
	cmd.Flags().Int("limit", 0, "Limit the number of results returned.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	cmd.Flags().Bool("fetch-sg-owned-wfs", true, "Include workflows owned by StackGuardian.")
	cmd.Flags().String("latest-wf-run-statuses", "", "Filter by latest run status (comma-separated).")
	cmd.Flags().String("latest-drift-run-statuses", "", "Filter by latest drift run status (comma-separated).")
	cmd.Flags().String("repo-id", "", "Filter by VCS repository.")
	cmd.Flags().String("use-marketplace-template", "", "Filter by marketplace template usage (true/false).")
	cmd.Flags().String("iac-template-id", "", "Filter by IaC template id.")
	cmd.Flags().String("parent-iac-template-id", "", "Filter by parent IaC template id.")
	cmd.Flags().String("terraform-version", "", "Filter by Terraform version.")
	cmd.Flags().String("runner-names", "", "Filter by runner names (comma-separated).")
	cmd.Flags().String("drift-check", "", "Any value restricts to workflows with drift information.")
	cmd.Flags().String("resource-names", "", "Filter by resource names (comma-separated).")
	cmd.Flags().String("description", "", "Filter by description.")
	cmd.Flags().String("tags", "", "Filter by tags (comma-separated).")
	cmd.Flags().String("is-active", "", "Filter by active state (0 or 1).")
	return cmd
}
