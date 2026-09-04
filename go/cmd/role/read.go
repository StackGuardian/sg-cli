// Code generated for sg-cli; edit the generator spec or this file directly.

package role

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get details of a role",
		Long:  `Get details of a role`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			roleId := common.String(cmd, "role-id")
			response, err := c.AccessManagement.ReadRole(ctx, org, roleId)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("role-id", "", "The role ID.")
	_ = cmd.MarkFlagRequired("role-id")
	return cmd
}
