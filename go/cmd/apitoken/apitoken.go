// Code generated for sg-cli; edit the generator spec or this file directly.

package apitoken

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewApiTokenCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-token",
		Short: "Manage legacy user API keys",
		Long:  `Manage legacy user API keys`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create or retrieve the caller's API key, or look up a runner group's key
  delete                 Revoke the API key of a principal`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newCreateCmd(c),
		newDeleteCmd(c),
	)
	return cmd
}
