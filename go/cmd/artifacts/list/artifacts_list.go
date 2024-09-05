package list

import (
	"context"
	"os"
	"strings"

	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

type RunOptions struct {
	OutputJson bool
}

func NewListCmd(c *client.Client) *cobra.Command {
	opts := &RunOptions{}
	// listCmd represents the list command
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List Artifacts",
		Long:  `List Artifacts`,
		Run: func(cmd *cobra.Command, args []string) {
			response, err := c.Workflows.ListAllArtifacts(
				context.Background(),
				cmd.Parent().Flags().Lookup("org").Value.String(),
				cmd.Parent().Flags().Lookup("workflow-id").Value.String(),
				cmd.Parent().Flags().Lookup("workflow-group").Value.String(),
			)
			if err != nil {
				cmd.PrintErrln("== Failed To List All Artifacts From Workflow ==")
				if strings.Contains(err.Error(), "the server responded with nothing") {
					cmd.PrintErrln("No artifacts found for this workflow")
					os.Exit(-1)
				}
				cmd.Println(err)
				os.Exit(-1)
			}

			if opts.OutputJson {
				cmd.Println(response)
			}

			cmd.Println(response)
		},
	}
	listCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output execution response as json to STDIN.")
	return listCmd
}
