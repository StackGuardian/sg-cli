// Code generated for sg-cli; edit the generator spec or this file directly.

package apiaccess

import (
	"context"
	"fmt"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newListCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List API accesses",
		Long:  `List API accesses`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			var accessType *sggosdk.ListAllApiAccessesRequestAccessType
			if v := common.OptString(cmd, "access-type"); v != nil {
				if *v != "APIKEY" && *v != "OIDC" {
					return fmt.Errorf("--access-type must be APIKEY or OIDC, got %q", *v)
				}
				t := sggosdk.ListAllApiAccessesRequestAccessType(*v)
				accessType = &t
			}
			response, err := c.AccessManagement.ListAllApiAccesses(ctx, org, &sggosdk.ListAllApiAccessesRequest{Limit: common.OptInt(cmd, "limit"), LastEvaluatedKey: common.OptString(cmd, "last-evaluated-key"), ResourceName: common.OptString(cmd, "resource-name"), Roles: common.OptString(cmd, "roles"), Status: common.OptString(cmd, "status"), AccessType: accessType})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().Int("limit", 0, "Limit the number of results returned.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	cmd.Flags().String("access-type", "", "Filter by access type (APIKEY, OIDC).")
	cmd.Flags().String("resource-name", "", "Filter by resource name.")
	cmd.Flags().String("roles", "", "Filter by roles (comma-separated).")
	cmd.Flags().String("status", "", "Filter by status (active, expired).")
	return cmd
}
