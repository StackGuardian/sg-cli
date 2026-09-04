// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflowrun

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newFactsCreateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create -- payload.json",
		Short: "Create or replace the facts of a stack workflow run",
		Long:  `Create or replace the facts of a stack workflow run`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			workflowId := common.String(cmd, "workflow-id")
			runId := common.String(cmd, "run-id")
			factsId := common.String(cmd, "facts-id")
			var req map[string]interface{}
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.StackWorkflowRunFacts.CreateStackWorkflowRunFacts(ctx, org, stackId, workflowId, workflowGroup, runId, factsId, req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Stack workflow run facts created.")
			return nil
		},
	}
	cmd.Flags().String("facts-id", "", "The workflow run facts ID.")
	_ = cmd.MarkFlagRequired("facts-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
