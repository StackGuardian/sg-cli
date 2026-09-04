// Code generated for sg-cli; edit the generator spec or this file directly.

package role

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a role",
		Long:  `Delete a role`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			roleId := common.String(cmd, "role-id")
			if err := c.AccessManagement.DeleteRole(ctx, org, roleId); err != nil {
				return err
			}
			common.Done(cmd, nil, "Role deleted successfully.")
			return nil
		},
	}
	cmd.Flags().String("role-id", "", "The role ID.")
	_ = cmd.MarkFlagRequired("role-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
