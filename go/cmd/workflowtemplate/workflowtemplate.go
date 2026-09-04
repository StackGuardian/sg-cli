// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowtemplate

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewWorkflowTemplateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow-template",
		Short: "Manage workflow (IaC) templates",
		Long:  `Manage workflow (IaC) templates`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  read                   Get a workflow template
  create                 Create a workflow template
  update                 Update a workflow template
  delete                 Delete a workflow template
  revisions              Manage revisions of workflow (IaC) templates`)
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
		Short: "Manage revisions of workflow (IaC) templates",
		Long:  `Manage revisions of workflow (IaC) templates`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create a new revision of a workflow template
  read                   Get a workflow template revision
  update                 Update a workflow template revision
  delete                 Delete a workflow template revision`)
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
