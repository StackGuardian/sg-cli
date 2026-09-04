// Code generated for sg-cli; edit the generator spec or this file directly.

package stacktemplate

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewStackTemplateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stack-template",
		Short: "Manage stack (IaC group) templates",
		Long:  `Manage stack (IaC group) templates`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  read                   Get a stack template
  create                 Create a stack template
  update                 Update a stack template
  delete                 Delete a stack template
  revisions              Manage revisions of stack (IaC group) templates`)
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
		Short: "Manage revisions of stack (IaC group) templates",
		Long:  `Manage revisions of stack (IaC group) templates`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create a new revision of a stack template
  read                   Get a stack template revision
  update                 Update a stack template revision
  delete                 Delete a stack template revision`)
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
