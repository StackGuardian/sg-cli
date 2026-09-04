// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newListByOwnerCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-by-owner",
		Short: "List templates owned by an organization",
		Long:  `List templates owned by an organization`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			ownerOrg := common.String(cmd, "owner-org")
			response, err := c.Templates.ListAllTemplatesBasedOnOwnerOrg(ctx, ownerOrg, &sggosdk.ListAllTemplatesBasedOnOwnerOrgRequest{SgOrgid: org, IsPublic: common.OptString(cmd, "is-public"), SearchQuery: common.OptString(cmd, "search-query"), IsSharedTemplate: common.OptBool(cmd, "shared"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("owner-org", "", "Organization that owns the templates.")
	_ = cmd.MarkFlagRequired("owner-org")
	cmd.Flags().String("is-public", "", "0 for private, 1 for public templates.")
	cmd.Flags().String("search-query", "", "Search by name, tags or description.")
	cmd.Flags().Bool("shared", false, "List shared templates.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	return cmd
}
