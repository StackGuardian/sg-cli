package delete

import (
	"context"
	"fmt"
	"os"
	"strings"

	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

type RunOptions struct {
	OutputJson  bool
	Org         string
	WfgGrp      string
	StackId     string
	ForceDelete bool
}

func NewDeleteCmd(c *client.Client) *cobra.Command {
	opts := &RunOptions{}
	// deleteCmd represents the delete command
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete the Stack from workflow group",
		Long:  `Delete the Stack from workflow group. Use option --force-delete to delete the Stack along with all of its workflows.`,
		Run: func(cmd *cobra.Command, args []string) {
			opts.Org = cmd.Parent().PersistentFlags().Lookup("org").Value.String()
			opts.WfgGrp = cmd.Parent().PersistentFlags().Lookup("workflow-group").Value.String()

			response, err := executeStackDeletion(c, cmd, opts)
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}

			if opts.OutputJson && response != nil {
				cmd.Println(response)
			}
			cmd.Println("Stack deleted successfully.")
		},
	}

	deleteCmd.Flags().StringVar(&opts.StackId, "stack-id", "", "The Stack ID to delete.")
	deleteCmd.Flags().BoolVar(&opts.ForceDelete, "force-delete", false, "The force-delete flag will delete the Stack along with all of its workflows. Use with caution.")
	deleteCmd.MarkFlagRequired("stack-id")

	deleteCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output execution response as json to STDIN.")

	return deleteCmd
}

// deleteStack attempts to delete a stack and returns the response and any error
func deleteStack(c *client.Client, opts *RunOptions) (interface{}, error) {
	return c.Stacks.DeleteStack(
		context.Background(),
		opts.Org,
		opts.StackId,
		opts.WfgGrp,
	)
}

// deleteAllStackWorkflows will find and delete all the Workflows that are part of this Stack
func deleteAllStackWorkflows(c *client.Client, cmd *cobra.Command, opts *RunOptions) {
	stackWorkflows, err := c.StackWorkflows.ListAllStackWorkflows(
		context.Background(),
		opts.Org,
		opts.StackId,
		opts.WfgGrp,
		&sggosdk.ListAllStackWorkflowsRequest{},
	)
	if err != nil {
		cmd.Println("An error occured while listing all the Stack Workflows to delete.")
		cmd.Println(err)
		os.Exit(-1)
	}
	for _, stackWf := range stackWorkflows.Msg {
		stackWfResourceIdSplit := strings.Split(stackWf.ResourceId, "/")
		err = c.StackWorkflows.DeleteStackWorkflow(
			context.Background(),
			opts.Org,
			opts.StackId,
			stackWfResourceIdSplit[len(stackWfResourceIdSplit)-1],
			opts.WfgGrp,
		)
		if err != nil {
			cmd.Println("An error occured while deleting Stack workflow " + stackWf.ResourceId)
			cmd.Println(err)
			os.Exit(-1)
		}
		cmd.Println("Stack workflow " + stackWf.ResourceId + " deleted successfully.")
	}
}

// executeStackDeletion handles the stack deletion logic and returns the response and any errors
func executeStackDeletion(c *client.Client, cmd *cobra.Command, opts *RunOptions) (interface{}, error) {
	response, err := deleteStack(c, opts)
	if err == nil {
		return response, nil
	}

	// Check if error is due to non-empty stack
	if !strings.Contains(err.Error(), "Stack is not empty") {
		return nil, err
	}

	// Handle non-empty stack error
	if !opts.ForceDelete {
		return nil, fmt.Errorf("this stack cannot be deleted since it contains workflows.\n" +
			"You can use the --force-delete flag to force the deletion of the stack along with all of its workflows")
	}

	// Force delete is enabled, delete all workflows first
	cmd.Println("Force deletion is enabled. Deleting the Stack's Workflows...")
	deleteAllStackWorkflows(c, cmd, opts)
	cmd.Println("All the Workflows in the Stack have been deleted. Deleting the Stack..")

	// Try deleting the stack again
	return deleteStack(c, opts)
}
