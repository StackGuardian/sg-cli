package delete

import (
	"context"
	"os"

	"github.com/StackGuardian/sg-cli/cmd/output"
	"github.com/StackGuardian/sg-cli/cmd/tui"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

type RunOptions struct {
	OutputJson bool
	Org        string
	WfgGrp     string
	WfId       string
}

func NewDeleteCmd(c *client.Client) *cobra.Command {
	opts := &RunOptions{}

	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a workflow from a workflow group",
		Long:  "Permanently delete a workflow. If --workflow-id is omitted, an interactive picker opens.",
		Run: func(cmd *cobra.Command, args []string) {
			opts.Org = cmd.Parent().PersistentFlags().Lookup("org").Value.String()
			opts.WfgGrp = cmd.Parent().PersistentFlags().Lookup("workflow-group").Value.String()
			opts.WfId = cmd.Flags().Lookup("workflow-id").Value.String()

			if opts.WfId == "" {
				id, err := pickWorkflow(c, opts.Org, opts.WfgGrp, "Delete Workflow")
				if err != nil {
					output.Error(err.Error())
					os.Exit(1)
				}
				opts.WfId = id
			}

			var response interface{}
			err := output.WithSpinner("Deleting workflow "+opts.WfId+"...", func() error {
				var apiErr error
				response, apiErr = c.Workflows.DeleteWorkflow(
					context.Background(),
					opts.Org,
					opts.WfId,
					opts.WfgGrp,
				)
				return apiErr
			})
			if err != nil {
				output.Error(err.Error())
				os.Exit(1)
			}

			if opts.OutputJson {
				cmd.Println(response)
			}

			output.Success("Workflow deleted successfully.")
		},
	}

	deleteCmd.Flags().StringVar(&opts.WfId, "workflow-id", "", "The workflow ID to delete. Omit to pick interactively.")
	deleteCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output API response as JSON.")

	return deleteCmd
}

func pickWorkflow(c *client.Client, org, wfGrp, title string) (string, error) {
	var response *sggosdk.WorkflowsListAll
	err := output.WithSpinner("Fetching workflows...", func() error {
		var apiErr error
		response, apiErr = c.Workflows.ListAllWorkflows(
			context.Background(),
			org,
			wfGrp,
			&sggosdk.ListAllWorkflowsRequest{},
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
	for i, wf := range response.Msg {
		items[i] = tui.Item{
			ID:          wf.ResourceName,
			Label:       wf.ResourceName,
			Description: wf.Description,
			Badge:       wf.LatestWfrunStatus,
		}
	}

	subtitle := org + " / " + wfGrp
	return tui.NewPicker(title, subtitle, items)
}
