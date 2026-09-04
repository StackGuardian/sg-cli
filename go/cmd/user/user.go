// Code generated for sg-cli; edit the generator spec or this file directly.

package user

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewUserCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users and SSO groups",
		Long:  `Manage users and SSO groups`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  read                   Get a user or SSO group
  create                 Invite a user or SSO group to the organization
  delete                 Remove a user or SSO group from the organization
  update                 Update a user or SSO group
  list                   List users and SSO groups`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newReadCmd(c),
		newCreateCmd(c),
		newDeleteCmd(c),
		newUpdateCmd(c),
		newListCmd(c),
	)
	return cmd
}
