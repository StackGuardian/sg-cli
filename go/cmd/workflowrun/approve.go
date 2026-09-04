// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowrun

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newApproveCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Approve (or, with --reject, reject) a workflow run waiting for approval",
		Long:  `Approve (or, with --reject, reject) a workflow run waiting for approval`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			workflowId := common.String(cmd, "workflow-id")
			runId := common.String(cmd, "run-id")
			reject := common.Bool(cmd, "reject")
			if err := common.Require(cmd, "workflow-group", "workflow-id"); err != nil {
				return err
			}
			approve := !reject
			body := sggosdk.WorkflowRunApproval{Approve: &approve, Message: common.OptString(cmd, "message"), ApprovalStep: common.OptString(cmd, "approval-step"), ReasonForApprovalRequired: common.OptString(cmd, "reason")}
			response, err := c.WorkflowRuns.ApproveWorkflowRun(ctx, org, workflowId, workflowGroup, runId, &body)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Approval decision submitted.")
			return nil
		},
	}
	cmd.Flags().String("run-id", "", "The workflow run ID.")
	_ = cmd.MarkFlagRequired("run-id")
	cmd.Flags().Bool("reject", false, "Reject instead of approve.")
	cmd.Flags().String("message", "", "Message to attach to the decision.")
	cmd.Flags().String("approval-step", "", "Name of the approval step.")
	cmd.Flags().String("reason", "", "Reason the approval was required.")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
