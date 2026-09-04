// Code generated for sg-cli; edit the generator spec or this file directly.

package organization

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Archive an organization and all of its resources",
		Long:  `Archive an organization and all of its resources`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.Organizations.DeleteOrganization(ctx, org)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Organization archived successfully.")
			return nil
		},
	}
	cmd.Flags().String("org", "", "The organization name on Stackguardian platform.")
	_ = cmd.MarkFlagRequired("org")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
