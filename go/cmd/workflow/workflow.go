package workflow

import (
	"github.com/StackGuardian/sg-cli/cmd/workflow/apply"
	"github.com/StackGuardian/sg-cli/cmd/workflow/create"
	"github.com/StackGuardian/sg-cli/cmd/workflow/delete"
	"github.com/StackGuardian/sg-cli/cmd/workflow/destroy"
	"github.com/StackGuardian/sg-cli/cmd/workflow/list"
	"github.com/StackGuardian/sg-cli/cmd/workflow/read"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewWorkflowCmd(c *client.Client) *cobra.Command {
	var workflowCmd = &cobra.Command{
		Use:   "workflow",
		Short: "Manage workflows",
		Long:  "Manage workflows in the StackGuardian platform.",
	}

	workflowCmd.PersistentFlags().String("org", "", "The organization name on StackGuardian platform.")
	workflowCmd.MarkPersistentFlagRequired("org")

	workflowCmd.PersistentFlags().String("workflow-group", "", "The workflow group under the organization.")
	workflowCmd.MarkPersistentFlagRequired("workflow-group")

	workflowCmd.AddCommand(read.NewReadCmd(c))
	workflowCmd.AddCommand(delete.NewDeleteCmd(c))
	workflowCmd.AddCommand(create.NewCreateCmd(c))
	workflowCmd.AddCommand(apply.NewApplyCmd(c))
	workflowCmd.AddCommand(list.NewListCmd(c))
	workflowCmd.AddCommand(destroy.NewDestroyCmd(c))

	return workflowCmd
}
