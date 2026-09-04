// Code generated for sg-cli; edit the generator spec or this file directly.

package connector

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewConnectorCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connector",
		Short: "Manage connectors (cloud and VCS integrations)",
		Long:  `Manage connectors (cloud and VCS integrations)`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create a new connector
  read                   Get details of a connector
  update                 Update a connector
  delete                 Delete a connector
  list                   List connectors
  authenticate           Verify that a connector can authenticate against its provider
  list-accounts          List the cloud accounts or subscriptions reachable with a connector's credentials
  github-repos           List the repositories accessible to a GITHUB_COM connector
  repos                  List the repositories reachable through a VCS connector`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newCreateCmd(c),
		newReadCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newListCmd(c),
		newAuthenticateCmd(c),
		newListAccountsCmd(c),
		newGithubReposCmd(c),
		newReposCmd(c),
	)
	return cmd
}
