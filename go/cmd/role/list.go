// Code generated for sg-cli; edit the generator spec or this file directly.

package role

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
		Short: "List roles",
		Long:  `List roles`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.AccessManagement.ListAllRoles(ctx, org, &sggosdk.ListAllRolesRequest{IncludePredefinedRoles: common.OptBool(cmd, "include-predefined-roles")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().Bool("include-predefined-roles", false, "Include predefined system roles.")
	return cmd
}
