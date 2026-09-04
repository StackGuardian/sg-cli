// Code generated for sg-cli; edit the generator spec or this file directly.

package runnergroup

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewRunnerGroupCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runner-group",
		Short: "Manage private runner groups",
		Long:  `Manage private runner groups`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create a new runner group
  read                   Get details of a runner group
  update                 Update a runner group
  delete                 Delete a runner group
  list                   List runner groups
  register               Register a runner node (use the runner group's own token as the API key)
  storage-backend-auth   Generate short-lived credentials for the runner group's storage backend
  deregister-runner      Deregister a runner from a runner group
  update-runner-state    Update the state of a runner (ACTIVE or DRAINING)`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newCreateCmd(c),
		newReadCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newListCmd(c),
		newRegisterCmd(c),
		newStorageBackendAuthCmd(c),
		newDeregisterRunnerCmd(c),
		newUpdateRunnerStateCmd(c),
	)
	return cmd
}
