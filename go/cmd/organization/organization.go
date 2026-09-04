// Code generated for sg-cli; edit the generator spec or this file directly.

package organization

import (
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewOrganizationCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "organization",
		Short: "Manage organizations",
		Long:  `Manage organizations`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  read                   Get details of an organization
  create                 Create a new organization (requires a user token, not an API key)
  list                   List the organizations the caller belongs to (requires a user token)
  update                 Update the settings of an organization
  delete                 Archive an organization and all of its resources
  workflows              List workflows across every workflow group of an organization
  count-workflows        Count workflows across every workflow group of an organization
  bulk-action            Activate or deactivate workflows and stacks in bulk`)
		},
	}
	cmd.AddCommand(
		newReadCmd(c),
		newCreateCmd(c),
		newListCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newWorkflowsCmd(c),
		newCountWorkflowsCmd(c),
		newBulkActionCmd(c),
	)
	return cmd
}
