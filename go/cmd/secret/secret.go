// Code generated for sg-cli; edit the generator spec or this file directly.

package secret

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewSecretCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Manage secrets",
		Long:  `Manage secrets`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create a new secret
  update                 Update a secret
  delete                 Delete a secret
  list                   List secrets
  read-bulk              Read several secrets by name`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newCreateCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newListCmd(c),
		newReadBulkCmd(c),
	)
	return cmd
}
