package read

import (
	"context"
	"os"

	"github.com/StackGuardian/sg-cli/cmd/output"
	"github.com/StackGuardian/sg-cli/cmd/tui"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewReadCmd(c *client.Client) *cobra.Command {
	var outputJson bool
	var wfId string

	var readCmd = &cobra.Command{
		Use:   "read",
		Short: "Get details of a workflow",
		Long:  "Fetch and display the full details of a workflow. If --workflow-id is omitted, an interactive picker opens.",
		Run: func(cmd *cobra.Command, args []string) {
			org := cmd.Parent().Flags().Lookup("org").Value.String()
			wfGrp := cmd.Parent().Flags().Lookup("workflow-group").Value.String()

			if wfId == "" {
				id, err := pickWorkflow(c, org, wfGrp, "Read Workflow")
				if err != nil {
					output.Error(err.Error())
					os.Exit(1)
				}
				wfId = id
			}

			var response interface{}
			err := output.WithSpinner("Fetching workflow "+wfId+"...", func() error {
				var apiErr error
				response, apiErr = c.Workflows.ReadWorkflow(
					context.Background(),
					org,
					wfId,
					wfGrp,
				)
				return apiErr
			})
			if err != nil {
				output.Error(err.Error())
				os.Exit(1)
			}

			if outputJson {
				cmd.Println(response)
				return
			}
			cmd.Println(response)
		},
	}

	readCmd.Flags().StringVar(&wfId, "workflow-id", "", "The workflow ID to retrieve. Omit to pick interactively.")
	readCmd.Flags().BoolVar(&outputJson, "output-json", false, "Output API response as JSON.")

	return readCmd
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
