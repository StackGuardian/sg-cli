package list

import (
	"context"
	"os"
	"strings"

	"github.com/StackGuardian/sg-cli/cmd/output"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

type RunOptions struct {
	OutputJson bool
}

func NewListCmd(c *client.Client) *cobra.Command {
	opts := &RunOptions{}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List artifacts for a workflow",
		Long:  "List all artifacts produced by a workflow.",
		Run: func(cmd *cobra.Command, args []string) {
			org := cmd.Parent().Flags().Lookup("org").Value.String()
			wfId := cmd.Parent().Flags().Lookup("workflow-id").Value.String()
			wfGrp := cmd.Parent().Flags().Lookup("workflow-group").Value.String()

			var response interface{}
			err := output.WithSpinner("Fetching artifacts...", func() error {
				var apiErr error
				response, apiErr = c.Workflows.ListAllWorkflowArtifacts(
					context.Background(),
					org,
					wfId,
					wfGrp,
				)
				return apiErr
			})
			if err != nil {
				if strings.Contains(err.Error(), "the server responded with nothing") {
					output.Warning("No artifacts found for this workflow.")
					os.Exit(0)
				}
				output.Error("Failed to list artifacts: " + err.Error())
				os.Exit(1)
			}

			cmd.Println(response)
		},
	}

	listCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output API response as JSON.")
	return listCmd
}
