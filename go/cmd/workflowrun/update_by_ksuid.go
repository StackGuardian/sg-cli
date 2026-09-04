// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowrun

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newUpdateByKsuidCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-by-ksuid -- payload.json",
		Short: "Update a workflow run by KSUIDs (runner and machine tokens only)",
		Long:  `Update a workflow run by KSUIDs (runner and machine tokens only)`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			parentKsuid := common.String(cmd, "parent-ksuid")
			resourceKsuid := common.String(cmd, "resource-ksuid")
			var req sggosdk.PatchedWorkflowRun
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.WorkflowRuns.UpdateWorkflowRunByKsuid(ctx, org, parentKsuid, resourceKsuid, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Workflow run updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("parent-ksuid", "", "KSUID of the parent workflow.")
	_ = cmd.MarkFlagRequired("parent-ksuid")
	cmd.Flags().String("resource-ksuid", "", "KSUID of the workflow run.")
	_ = cmd.MarkFlagRequired("resource-ksuid")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
