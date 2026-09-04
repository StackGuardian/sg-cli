// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflowrun

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newFactsReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get the facts of a stack workflow run",
		Long:  `Get the facts of a stack workflow run`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			workflowId := common.String(cmd, "workflow-id")
			runId := common.String(cmd, "run-id")
			factsId := common.String(cmd, "facts-id")
			response, err := c.StackWorkflowRunFacts.GetStackWorkflowRunFacts(ctx, org, stackId, workflowId, workflowGroup, runId, factsId)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("facts-id", "", "The workflow run facts ID.")
	_ = cmd.MarkFlagRequired("facts-id")
	return cmd
}
