package read

import (
	"context"
	"os"

	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewReadCmd(c *client.Client) *cobra.Command {
	// readCmd represents the read command
	var readCmd = &cobra.Command{
		Use:   "read",
		Short: "Get details of a workflow",
		Long:  `Get details of a workflow.`,
		Run: func(cmd *cobra.Command, args []string) {
			response, err := c.Workflows.Get(
				context.Background(),
				cmd.Parent().Flags().Lookup("org").Value.String(),
				cmd.Flags().Lookup("workflow-id").Value.String(),
				cmd.Parent().Flags().Lookup("workflow-group").Value.String(),
			)
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			cmd.Println(response)
		},
	}

	readCmd.Flags().String("workflow-id", "", "The workflow ID to retrieve.")
	readCmd.MarkFlagRequired("workflow-id")

	return readCmd
}
