package delete

import (
	"context"
	"os"

	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

type RunOptions struct {
	OutputJson bool
	Org        string
	WfgGrp     string
	WfId       string
}

func NewDeleteCmd(c *client.Client) *cobra.Command {
	opts := &RunOptions{}
	// deleteCmd represents the delete command
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete the workflow from workflow group",
		Long:  `Delete the workflow from workflow group`,
		Run: func(cmd *cobra.Command, args []string) {
			opts.Org = cmd.Parent().PersistentFlags().Lookup("org").Value.String()
			opts.WfgGrp = cmd.Parent().PersistentFlags().Lookup("workflow-group").Value.String()
			opts.WfId = cmd.Flags().Lookup("workflow-id").Value.String()

			response, err := c.Workflows.DeleteWorkflow(
				context.Background(),
				opts.Org,
				opts.WfId,
				opts.WfgGrp,
			)
			if err != nil {
				cmd.Println(err)
				os.Exit(-1)
			}
			if opts.OutputJson {
				cmd.Println(response)
			}
			cmd.Println("Workflow deleted successfully.")
		},
	}

	deleteCmd.Flags().String("workflow-id", "", "The workflow ID to retrieve.")
	deleteCmd.MarkFlagRequired("workflow-id")

	deleteCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output execution response as json to STDIN.")

	return deleteCmd
}
