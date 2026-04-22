package stack

import (
	"github.com/StackGuardian/sg-cli/cmd/stack/apply"
	"github.com/StackGuardian/sg-cli/cmd/stack/create"
	"github.com/StackGuardian/sg-cli/cmd/stack/delete"
	"github.com/StackGuardian/sg-cli/cmd/stack/destroy"
	"github.com/StackGuardian/sg-cli/cmd/stack/outputs"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewStackCmd(c *client.Client) *cobra.Command {
	var stackCmd = &cobra.Command{
		Use:   "stack",
		Short: "Manage stacks",
		Long:  "Manage stacks in the StackGuardian platform.",
	}

	stackCmd.PersistentFlags().String("org", "", "The organization name on StackGuardian platform.")
	stackCmd.MarkPersistentFlagRequired("org")

	stackCmd.PersistentFlags().String("workflow-group", "", "The workflow group under the organization.")
	stackCmd.MarkPersistentFlagRequired("workflow-group")

	stackCmd.AddCommand(outputs.NewOutputsCmd(c))
	stackCmd.AddCommand(destroy.NewDestroyCmd(c))
	stackCmd.AddCommand(create.NewCreateCmd(c))
	stackCmd.AddCommand(delete.NewDeleteCmd(c))
	stackCmd.AddCommand(apply.NewApplyCmd(c))

	return stackCmd
}
