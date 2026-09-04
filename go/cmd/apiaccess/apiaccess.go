// Code generated for sg-cli; edit the generator spec or this file directly.

package apiaccess

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewApiAccessCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-access",
		Short: "Manage API accesses (API keys and OIDC identities)",
		Long:  `Manage API accesses (API keys and OIDC identities)`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create an API access
  read                   Get an API access
  update                 Update an API access
  delete                 Delete an API access
  regenerate-key         Regenerate the key of an API access
  list                   List API accesses`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newCreateCmd(c),
		newReadCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newRegenerateKeyCmd(c),
		newListCmd(c),
	)
	return cmd
}
