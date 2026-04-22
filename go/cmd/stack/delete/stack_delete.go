package delete

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/StackGuardian/sg-cli/cmd/output"
	"github.com/StackGuardian/sg-cli/cmd/tui"
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

	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a stack from a workflow group",
		Long:  "Permanently delete a stack. If --stack-id is omitted, an interactive picker opens. Use --force-delete to also remove all workflows within the stack.",
		Run: func(cmd *cobra.Command, args []string) {
			opts.Org = cmd.Parent().PersistentFlags().Lookup("org").Value.String()
			opts.WfgGrp = cmd.Parent().PersistentFlags().Lookup("workflow-group").Value.String()

			if opts.StackId == "" {
				id, err := pickStack(c, opts.Org, opts.WfgGrp, "Delete Stack")
				if err != nil {
					output.Error(err.Error())
					os.Exit(1)
				}
				opts.StackId = id
			}

			var response interface{}
			err := output.WithSpinner("Deleting stack "+opts.StackId+"...", func() error {
				var apiErr error
				response, apiErr = executeStackDeletion(c, cmd, opts)
				return apiErr
			})
			if err != nil {
				output.Error(err.Error())
				os.Exit(1)
			}

			if opts.OutputJson && response != nil {
				cmd.Println(response)
			}

			output.Success("Stack deleted successfully.")
		},
	}

	deleteCmd.Flags().StringVar(&opts.StackId, "stack-id", "", "The stack ID to delete. Omit to pick interactively.")
	deleteCmd.Flags().BoolVar(&opts.ForceDelete, "force-delete", false, "Delete the stack along with all of its workflows.")
	deleteCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output API response as JSON.")

	return deleteCmd
}

func pickStack(c *client.Client, org, wfGrp, title string) (string, error) {
	var response *sggosdk.GeneratedStackListAllResponse
	err := output.WithSpinner("Fetching stacks...", func() error {
		var apiErr error
		response, apiErr = c.Stacks.ListAllStacks(
			context.Background(),
			org,
			wfGrp,
			&sggosdk.ListAllStacksRequest{},
		)
		return apiErr
	})
	if err != nil {
		return "", err
	}

	if len(response.Msg) == 0 {
		return "", nil
	}

	items := make([]tui.Item, len(response.Msg))
	for i, s := range response.Msg {
		items[i] = tui.Item{
			ID:          s.ResourceName,
			Label:       s.ResourceName,
			Description: s.Description,
			Badge:       s.LatestWfStatus,
		}
	}

	subtitle := org + " / " + wfGrp
	return tui.NewPicker(title, subtitle, items)
}

func deleteStack(c *client.Client, opts *RunOptions) (interface{}, error) {
	return c.Stacks.DeleteStack(
		context.Background(),
		opts.Org,
		opts.StackId,
		opts.WfgGrp,
	)
}

func deleteAllStackWorkflows(c *client.Client, cmd *cobra.Command, opts *RunOptions) {
	stackWorkflows, err := c.StackWorkflows.ListAllStackWorkflows(
		context.Background(),
		opts.Org,
		opts.StackId,
		opts.WfgGrp,
		&sggosdk.ListAllStackWorkflowsRequest{},
	)
	if err != nil {
		output.Error("Failed to list stack workflows: " + err.Error())
		os.Exit(1)
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
			output.Error("Failed to delete stack workflow " + stackWf.ResourceId + ": " + err.Error())
			os.Exit(1)
		}
		output.Info("Deleted stack workflow: " + stackWf.ResourceId)
	}
}

func executeStackDeletion(c *client.Client, cmd *cobra.Command, opts *RunOptions) (interface{}, error) {
	response, err := deleteStack(c, opts)
	if err == nil {
		return response, nil
	}

	if !strings.Contains(err.Error(), "Stack is not empty") {
		return nil, err
	}

	if !opts.ForceDelete {
		return nil, fmt.Errorf("stack contains workflows — use --force-delete to remove them first")
	}

	output.Warning("Force deletion enabled. Removing all workflows in the stack first...")
	deleteAllStackWorkflows(c, cmd, opts)
	output.Info("All workflows deleted. Removing the stack...")

	return deleteStack(c, opts)
}
