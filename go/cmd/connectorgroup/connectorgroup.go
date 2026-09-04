// Code generated for sg-cli; edit the generator spec or this file directly.

package connectorgroup

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewConnectorGroupCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connector-group",
		Short: "Manage connector groups (cloud integration groups)",
		Long:  `Manage connector groups (cloud integration groups)`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create a new connector group
  read                   Get details of a connector group
  update                 Update a connector group
  delete                 Delete a connector group
  authenticate           Authenticate a connector group
  list                   List connector groups
  discovery-scan         Trigger discovery scans across the group's child connectors
  connectors             Manage the connectors inside a connector group`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newCreateCmd(c),
		newReadCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newAuthenticateCmd(c),
		newListCmd(c),
		newDiscoveryScanCmd(c),
		newConnectorsCmd(c),
	)
	return cmd
}
func newConnectorsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connectors",
		Short: "Manage the connectors inside a connector group",
		Long:  `Manage the connectors inside a connector group`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  read                   Get a connector in the group
  update                 Update a connector in the group
  delete                 Delete a connector from the group
  list                   List the connectors in the group
  authenticate           Verify that a connector in the group can authenticate`)
		},
	}
	cmd.PersistentFlags().String("connector-group", "", "The connector group name.")
	_ = cmd.MarkPersistentFlagRequired("connector-group")
	cmd.AddCommand(
		newConnectorsReadCmd(c),
		newConnectorsUpdateCmd(c),
		newConnectorsDeleteCmd(c),
		newConnectorsListCmd(c),
		newConnectorsAuthenticateCmd(c),
	)
	return cmd
}
