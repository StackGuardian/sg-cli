package outputs

import (
	"context"
	"os"

	"github.com/StackGuardian/sg-cli/cmd/output"
	"github.com/StackGuardian/sg-cli/cmd/tui"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewOutputsCmd(c *client.Client) *cobra.Command {
	var stackId string

	var outputsCmd = &cobra.Command{
		Use:   "outputs",
		Short: "Get outputs from a stack",
		Long:  "Fetch and display the outputs produced by a stack. If --stack-id is omitted, an interactive picker opens.",
		Run: func(cmd *cobra.Command, args []string) {
			org := cmd.Parent().Flags().Lookup("org").Value.String()
			wfGrp := cmd.Parent().Flags().Lookup("workflow-group").Value.String()

			if stackId == "" {
				id, err := pickStack(c, org, wfGrp, "Stack Outputs")
				if err != nil {
					output.Error(err.Error())
					os.Exit(1)
				}
				stackId = id
			}

			var response interface{}
			err := output.WithSpinner("Fetching outputs for stack "+stackId+"...", func() error {
				var apiErr error
				response, apiErr = c.Stacks.ReadStackOutputs(
					context.Background(),
					org,
					stackId,
					wfGrp,
				)
				return apiErr
			})
			if err != nil {
				output.Error(err.Error())
				os.Exit(1)
			}

			cmd.Println(response)
		},
	}

	outputsCmd.Flags().StringVar(&stackId, "stack-id", "", "The stack ID to retrieve outputs for. Omit to pick interactively.")

	return outputsCmd
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
