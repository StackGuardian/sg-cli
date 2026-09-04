// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowgroup

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewWorkflowGroupCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow-group",
		Short: "Manage workflow groups (nested groups are addressed as parent/child)",
		Long:  `Manage workflow groups (nested groups are addressed as parent/child)`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create a new workflow group
  read                   Get details of a workflow group
  update                 Update a workflow group
  delete                 Delete a workflow group
  create-child           Create a child workflow group inside a parent group
  list-children          List the child workflow groups of a parent group
  list                   List workflow groups
  list-resources         List the workflows, stacks and child groups inside a workflow group`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newCreateCmd(c),
		newReadCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newCreateChildCmd(c),
		newListChildrenCmd(c),
		newListCmd(c),
		newListResourcesCmd(c),
	)
	return cmd
}
