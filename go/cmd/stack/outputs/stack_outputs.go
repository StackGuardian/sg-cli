package outputs

import (
	"context"
	"os"

	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewOutputsCmd(c *client.Client) *cobra.Command {
	// outputsCmd represents the output command
	var outputsCmd = &cobra.Command{
		Use:   "outputs",
		Short: "Get outputs from stack",
		Long:  `Get outputs from stack.`,
		Run: func(cmd *cobra.Command, args []string) {
			response, err := c.Stacks.GetStackOutputs(
				context.Background(),
				cmd.Parent().Flags().Lookup("org").Value.String(),
				cmd.Flags().Lookup("stack-id").Value.String(),
				cmd.Parent().Flags().Lookup("workflow-group").Value.String(),
			)
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			cmd.Println(response)
		},
	}

	outputsCmd.Flags().String("stack-id", "", "The stack ID to retrieve.")
	outputsCmd.MarkFlagRequired("stack-id")

	return outputsCmd
}
