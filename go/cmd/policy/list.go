// Code generated for sg-cli; edit the generator spec or this file directly.

package policy

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
		Short: "List policies",
		Long:  `List policies`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.Policies.ListAllPolicies(ctx, org, &sggosdk.ListAllPoliciesRequest{Description: common.OptString(cmd, "description"), PolicyTypes: common.OptString(cmd, "policy-types"), ResourceNames: common.OptString(cmd, "resource-names"), Tags: common.OptString(cmd, "tags"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key"), Limit: common.OptInt(cmd, "limit"), SearchQuery: common.OptString(cmd, "search-query")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().Int("limit", 0, "Limit the number of results returned.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	cmd.Flags().String("description", "", "Filter by description.")
	cmd.Flags().String("policy-types", "", "Filter by policy types (default GENERAL).")
	cmd.Flags().String("resource-names", "", "Filter by resource names.")
	cmd.Flags().String("tags", "", "Filter by tags.")
	cmd.Flags().String("search-query", "", "Search by resource name, tags or description.")
	return cmd
}
