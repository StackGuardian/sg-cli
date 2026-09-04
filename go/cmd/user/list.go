// Code generated for sg-cli; edit the generator spec or this file directly.

package user

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
		Short: "List users and SSO groups",
		Long:  `List users and SSO groups`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.AccessManagement.ListAllUsers(ctx, org, &sggosdk.ListAllUsersRequest{Alias: common.OptString(cmd, "alias"), EntityType: common.OptString(cmd, "entity-type"), LoginMethod: common.OptString(cmd, "login-method"), Roles: common.OptString(cmd, "roles"), UserId: common.OptString(cmd, "user-id")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("alias", "", "Filter by alias.")
	cmd.Flags().String("entity-type", "", "Filter by entity type (EMAIL, GROUP).")
	cmd.Flags().String("login-method", "", "Filter by login method.")
	cmd.Flags().String("roles", "", "Filter by roles (comma-separated).")
	cmd.Flags().String("user-id", "", "Filter by user id.")
	return cmd
}
