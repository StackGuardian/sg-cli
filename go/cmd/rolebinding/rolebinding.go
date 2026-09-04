// Code generated for sg-cli; edit the generator spec or this file directly.

package rolebinding

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewRoleBindingCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role-binding",
		Short: "Inspect role bindings",
		Long:  `Inspect role bindings`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  read                   Get a role binding`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newReadCmd(c),
	)
	return cmd
}
