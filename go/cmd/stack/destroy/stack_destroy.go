package destroy

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
	Stack      string
}

func NewDestroyCmd(c *client.Client) *cobra.Command {
	opts := &RunOptions{}
	const dashboardURL = "https://app.stackguardian.io/orchestrator"

	var destroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Execute \"Destroy\" on an existing stack",
		Long:  "Trigger a Destroy run on an existing stack. If --stack-id is omitted, an interactive picker opens.",
		Run: func(cmd *cobra.Command, args []string) {
			opts.Org = cmd.Parent().PersistentFlags().Lookup("org").Value.String()
			opts.WfgGrp = cmd.Parent().PersistentFlags().Lookup("workflow-group").Value.String()
			opts.Stack = cmd.Flags().Lookup("stack-id").Value.String()

			if opts.Stack == "" {
				id, err := pickStack(c, opts.Org, opts.WfgGrp, "Destroy Stack")
				if err != nil {
					output.Error(err.Error())
					os.Exit(1)
				}
				opts.Stack = id
			}

			var response interface{}
			err := output.WithSpinner("Running destroy on stack "+opts.Stack+"...", func() error {
				var apiErr error
				response, apiErr = c.StackRuns.CreateStackRun(
					context.Background(),
					opts.Org,
					opts.Stack,
					opts.WfgGrp,
					&sggosdk.StackAction{
						ActionType: string(sggosdk.ActionEnumDestroy),
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

			output.Success("Stack destroy triggered successfully.")
			output.URL("View run at:", dashboardURL+
				"/orgs/"+opts.Org+
				"/wfgrps/"+opts.WfgGrp+
				"/stacks/"+opts.Stack+"?tab=runs")
		},
	}

	destroyCmd.Flags().StringVar(&opts.Stack, "stack-id", "", "The stack ID to destroy. Omit to pick interactively.")
	destroyCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output API response as JSON.")

	return destroyCmd
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
