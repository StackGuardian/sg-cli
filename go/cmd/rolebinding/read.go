// Code generated for sg-cli; edit the generator spec or this file directly.

package rolebinding

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get a role binding",
		Long:  `Get a role binding`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			roleBinding := common.String(cmd, "role-binding")
			response, err := c.AccessManagement.ReadRoleBinding(ctx, org, roleBinding)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("role-binding", "default", "The role binding name.")
	return cmd
}
