// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowrun

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newFactsUpdateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update -- payload.json",
		Short: "Update the facts of a workflow run",
		Long:  `Update the facts of a workflow run`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
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
			if err := common.Require(cmd, "workflow-group", "workflow-id"); err != nil {
				return err
			}
			response, err := c.WorkflowRunFacts.UpdateWorkflowRunFacts(ctx, org, workflowId, workflowGroup, runId, factsId, req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Workflow run facts updated.")
			return nil
		},
	}
	cmd.Flags().String("facts-id", "", "The workflow run facts ID.")
	_ = cmd.MarkFlagRequired("facts-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
