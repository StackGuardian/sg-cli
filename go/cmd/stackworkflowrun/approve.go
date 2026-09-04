// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflowrun

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
		Short: "Approve (or, with --reject, reject) a stack workflow run waiting for approval",
		Long:  `Approve (or, with --reject, reject) a stack workflow run waiting for approval`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			workflowId := common.String(cmd, "workflow-id")
			runId := common.String(cmd, "run-id")
			reject := common.Bool(cmd, "reject")
			approve := !reject
			body := sggosdk.WorkflowRunApproval{Approve: &approve, Message: common.OptString(cmd, "message"), ApprovalStep: common.OptString(cmd, "approval-step"), ReasonForApprovalRequired: common.OptString(cmd, "reason")}
			if err := c.StackWorkflowRuns.ApproveStackWorkflowRun(ctx, org, stackId, workflowId, workflowGroup, runId, &body); err != nil {
				return err
			}
			common.Done(cmd, nil, "Approval decision submitted.")
			return nil
		},
	}
	cmd.PersistentFlags().String("run-id", "", "The workflow run ID.")
	_ = cmd.MarkPersistentFlagRequired("run-id")
	cmd.Flags().Bool("reject", false, "Reject instead of approve.")
	cmd.Flags().String("message", "", "Message to attach to the decision.")
	cmd.Flags().String("approval-step", "", "Name of the approval step.")
	cmd.Flags().String("reason", "", "Reason the approval was required.")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
