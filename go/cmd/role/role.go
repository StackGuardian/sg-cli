// Code generated for sg-cli; edit the generator spec or this file directly.

package role

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewRoleCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "Manage roles",
		Long:  `Manage roles`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create a new role
  read                   Get details of a role
  update                 Update a role
  delete                 Delete a role
  list                   List roles`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newCreateCmd(c),
		newReadCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newListCmd(c),
	)
	return cmd
}
