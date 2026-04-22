package apply

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

func NewApplyCmd(c *client.Client) *cobra.Command {
	opts := &RunOptions{}
	const dashboardURL = "https://app.stackguardian.io/orchestrator"

	var applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Execute \"Apply\" on an existing workflow",
		Long:  "Trigger an Apply run on an existing workflow. If --workflow-id is omitted, an interactive picker opens.",
		Run: func(cmd *cobra.Command, args []string) {
			opts.Org = cmd.Parent().PersistentFlags().Lookup("org").Value.String()
			opts.WfgGrp = cmd.Parent().PersistentFlags().Lookup("workflow-group").Value.String()
			opts.WfId = cmd.Flags().Lookup("workflow-id").Value.String()

			if opts.WfId == "" {
				id, err := pickWorkflow(c, opts.Org, opts.WfgGrp, "Apply Workflow")
				if err != nil {
					output.Error(err.Error())
					os.Exit(1)
				}
				opts.WfId = id
			}

			var response interface{}
			err := output.WithSpinner("Running apply on "+opts.WfId+"...", func() error {
				var apiErr error
				response, apiErr = c.WorkflowRuns.CreateWorkflowRun(
					context.Background(),
					opts.Org,
					opts.WfId,
					opts.WfgGrp,
					&sggosdk.WorkflowRun{
						TerraformAction: &sggosdk.TerraformAction{
							Action: sggosdk.ActionEnumApply.Ptr(),
						},
					},
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

			output.Success("Workflow apply triggered successfully.")
			output.URL("View run at:", dashboardURL+
				"/orgs/"+opts.Org+
				"/wfgrps/"+opts.WfgGrp+
				"/wfs/"+opts.WfId+"?tab=runs")
		},
	}

	applyCmd.Flags().StringVar(&opts.WfId, "workflow-id", "", "The workflow ID to apply. Omit to pick interactively.")
	applyCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output API response as JSON.")

	return applyCmd
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
