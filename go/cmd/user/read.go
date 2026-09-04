// Code generated for sg-cli; edit the generator spec or this file directly.

package user

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get a user or SSO group",
		Long:  `Get a user or SSO group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.AccessManagement.ReadUser(ctx, org, &sggosdk.GetorRemoveUserFromOrganization{UserId: common.OptString(cmd, "user-id"), Alias: common.OptString(cmd, "alias")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("user-id", "", "User id (email) or SSO group id.")
	cmd.Flags().String("alias", "", "Alias of an SSO group.")
	cmd.MarkFlagsOneRequired("user-id", "alias")
	return cmd
}
