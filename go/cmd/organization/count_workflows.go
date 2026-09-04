// Code generated for sg-cli; edit the generator spec or this file directly.

package organization

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/organizations"
	"github.com/spf13/cobra"
)

func newCountWorkflowsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "count-workflows",
		Short: "Count workflows across every workflow group of an organization",
		Long:  `Count workflows across every workflow group of an organization`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.Organizations.CountWorkflows(ctx, org, &organizations.ListAllOrganizationWorkflowsRequest{
				ResourceNames: common.OptString(cmd, "resource-names"), Description: common.OptString(cmd, "description"), Tags: common.OptString(cmd, "tags"), IsActive: common.OptString(cmd, "is-active"),
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
	cmd.Flags().String("resource-names", "", "Filter by resource names (comma-separated).")
	cmd.Flags().String("description", "", "Filter by description.")
	cmd.Flags().String("tags", "", "Filter by tags (comma-separated).")
	cmd.Flags().String("is-active", "", "Filter by active state (0 or 1).")
	return cmd
}
