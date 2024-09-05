package list

import (
	"context"
	"os"

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
		Short: "List all workflows",
		Long:  `List all workflows`,
		Run: func(cmd *cobra.Command, args []string) {
			response, err := c.Workflows.ListAll(
				context.Background(),
				cmd.Parent().Flags().Lookup("org").Value.String(),
				cmd.Parent().Flags().Lookup("workflow-group").Value.String(),
			)
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}

			if opts.OutputJson {
				cmd.Println(response)
			}

			for _, workflow := range response.Msg {
				cmd.Println("> Workflow Name: ", workflow.ResourceName)
				cmd.Println("Description: ", workflow.Description)
				//New line for formatting
				cmd.Println()
			}
		},
	}
	listCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output execution response as json to STDIN.")
	return listCmd
}
