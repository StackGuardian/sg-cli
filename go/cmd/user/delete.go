// Code generated for sg-cli; edit the generator spec or this file directly.

package user

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Remove a user or SSO group from the organization",
		Long:  `Remove a user or SSO group from the organization`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.AccessManagement.DeleteUser(ctx, org, &sggosdk.GetorRemoveUserFromOrganization{UserId: common.OptString(cmd, "user-id"), Alias: common.OptString(cmd, "alias")})
			if err != nil {
				return err
			}
			common.Done(cmd, response, "User removed successfully.")
			return nil
		},
	}
	cmd.Flags().String("user-id", "", "User id (email) or SSO group id.")
	cmd.Flags().String("alias", "", "Alias of an SSO group.")
	common.AddOutputJSONFlag(cmd)
	cmd.MarkFlagsOneRequired("user-id", "alias")
	return cmd
}
