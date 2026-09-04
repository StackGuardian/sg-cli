// Code generated for sg-cli; edit the generator spec or this file directly.

package statebackend

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewStateBackendCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state-backend",
		Short: "Manage state backends (S3 or Azure Blob storage for state files)",
		Long:  `Manage state backends (S3 or Azure Blob storage for state files)`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create a new state backend
  read                   Get details of a state backend
  update                 Update a state backend
  delete                 Delete a state backend
  list                   List state backends
  list-statefiles        List the state files in a storage backend`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newCreateCmd(c),
		newReadCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newListCmd(c),
		newListStatefilesCmd(c),
	)
	return cmd
}
