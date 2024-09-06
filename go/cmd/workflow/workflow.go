package workflow

import (
	"fmt"

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
	// workflowCmd represents the workflow command
	var workflowCmd = &cobra.Command{
		Use:   "workflow",
		Short: "Manage workflows",
		Long:  `Manage workflows in Stackguardian platform.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(`Sub-commands:
  create      Create new workflow
  delete      Delete the workflow from workflow group
  apply       Execute "Apply" on existing workflow
  destroy     Execute "Destroy" on existing workflow
  read        Read, get details of a workflow
  list        List workflows`)
		},
	}

	workflowCmd.PersistentFlags().String("org", "", "The organization name on Stackguardian platform.")
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
