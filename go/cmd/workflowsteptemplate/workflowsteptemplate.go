// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowsteptemplate

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewWorkflowStepTemplateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow-step-template",
		Short: "Manage workflow step templates",
		Long:  `Manage workflow step templates`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  read                   Get a workflow step template
  create                 Create a workflow step template
  update                 Update a workflow step template
  delete                 Delete a workflow step template
  revisions              Manage revisions of workflow step templates`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newReadCmd(c),
		newCreateCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newRevisionsCmd(c),
	)
	return cmd
}
func newRevisionsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revisions",
		Short: "Manage revisions of workflow step templates",
		Long:  `Manage revisions of workflow step templates`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create a new revision of a workflow step template
  read                   Get a workflow step template revision
  update                 Update a workflow step template revision
  delete                 Delete a workflow step template revision`)
		},
	}
	cmd.AddCommand(
		newRevisionsCreateCmd(c),
		newRevisionsReadCmd(c),
		newRevisionsUpdateCmd(c),
		newRevisionsDeleteCmd(c),
	)
	return cmd
}
