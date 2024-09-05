package stack

import (
	"fmt"

	"github.com/StackGuardian/sg-cli/cmd/stack/apply"
	"github.com/StackGuardian/sg-cli/cmd/stack/create"
	"github.com/StackGuardian/sg-cli/cmd/stack/destroy"
	"github.com/StackGuardian/sg-cli/cmd/stack/outputs"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewStackCmd(c *client.Client) *cobra.Command {
	// stackCmd represents the Stack command
	var stackCmd = &cobra.Command{
		Use:   "stack",
		Short: "Manage stacks",
		Long:  `Manage stacks in Stackguardian platform.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(`Sub-commands:
  create      Create new stack
  apply       Execute "Apply" on existing stack
  destroy     Execute "Destroy" on existing stack
  outputs     Get outputs from stack`)
		},
	}

	stackCmd.PersistentFlags().String("org", "", "The organization name on Stackguardian platform.")
	stackCmd.MarkPersistentFlagRequired("org")

	stackCmd.PersistentFlags().String("workflow-group", "", "The workflow group under the organization.")
	stackCmd.MarkPersistentFlagRequired("workflow-group")

	stackCmd.AddCommand(outputs.NewOutputsCmd(c))
	stackCmd.AddCommand(destroy.NewDestroyCmd(c))
	stackCmd.AddCommand(create.NewCreateCmd(c))
	stackCmd.AddCommand(apply.NewApplyCmd(c))

	return stackCmd
}
