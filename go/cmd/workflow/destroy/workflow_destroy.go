package destroy

import (
	"context"
	"os"

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

func NewDestroyCmd(c *client.Client) *cobra.Command {
	opts := &RunOptions{}
	DASHBOARD_URL := "https://app.stackguardian.io/orchestrator"

	// destroyCmd represents the apply command
	var destroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Execute \"Destroy\" on existing workflow",
		Long:  `Execute "Destroy" on existing workflow`,
		Run: func(cmd *cobra.Command, args []string) {
			opts.Org = cmd.Parent().PersistentFlags().Lookup("org").Value.String()
			opts.WfgGrp = cmd.Parent().PersistentFlags().Lookup("workflow-group").Value.String()
			opts.WfId = cmd.Flags().Lookup("workflow-id").Value.String()
			response, err := c.WorkflowRuns.CreateWorkflowRun(
				context.Background(),
				opts.Org,
				opts.WfId,
				opts.WfgGrp,
				&sggosdk.WorkflowRun{
					TerraformAction: &sggosdk.TerraformAction{
						Action: sggosdk.ActionEnumDestroy.Ptr().Ptr(),
					},
				},
			)
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			if opts.OutputJson {
				cmd.Println(response)
			}
			cmd.Println("Workflow destroy run successfully.")
			workflowRunPath := DASHBOARD_URL +
				"/orgs/" +
				opts.Org +
				"/wfgrps/" +
				opts.WfgGrp +
				"/wfs/" +
				opts.WfId +
				"?tab=runs"
			cmd.Println("To view the workflow runs, please visit the following URL:")
			cmd.Println(workflowRunPath)
		},
	}

	destroyCmd.Flags().String("workflow-id", "", "The workflow ID to retrieve.")
	destroyCmd.MarkFlagRequired("workflow-id")

	destroyCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output execution response as json to STDIN.")

	return destroyCmd
}
