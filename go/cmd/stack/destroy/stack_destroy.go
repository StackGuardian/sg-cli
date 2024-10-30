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
	Stack      string
}

func NewDestroyCmd(c *client.Client) *cobra.Command {
	opts := &RunOptions{}
	DASHBOARD_URL := "https://app.stackguardian.io/orchestrator"
	// destroyCmd represents the apply command
	var destroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Execute \"Destroy\" on existing stack",
		Long:  `Execute "Destroy" on existing stack`,
		Run: func(cmd *cobra.Command, args []string) {
			opts.Org = cmd.Parent().PersistentFlags().Lookup("org").Value.String()
			opts.WfgGrp = cmd.Parent().PersistentFlags().Lookup("workflow-group").Value.String()
			opts.Stack = cmd.Flags().Lookup("stack-id").Value.String()
			response, err := c.StackWorkflowRuns.CreateStackRun(
				context.Background(),
				opts.Org,
				opts.Stack,
				opts.WfgGrp,
				&sggosdk.StackAction{
					ActionType: sggosdk.ActionTypeEnumDestroy,
				},
			)
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			if opts.OutputJson {
				cmd.Println(response)
			}
			workflowRunPath := DASHBOARD_URL +
				"/orgs/" +
				opts.Org +
				"/wfgrps/" +
				opts.WfgGrp +
				"/stacks/" +
				opts.Stack +
				"?tab=runs"
			cmd.Println("To view the workflow run, please visit the following URL:")
			cmd.Println(workflowRunPath)
			cmd.Println("Stack Workflow destroy run successfully.")

		},
	}

	destroyCmd.Flags().String("stack-id", "", "The stack ID to retrieve.")
	destroyCmd.MarkFlagRequired("stack-id")

	destroyCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output execution response as json to STDIN.")

	return destroyCmd
}
