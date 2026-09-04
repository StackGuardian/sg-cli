// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowrun

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newReadByKsuidCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read-by-ksuid",
		Short: "Get a workflow run by KSUIDs (runner and machine tokens only)",
		Long:  `Get a workflow run by KSUIDs (runner and machine tokens only)`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			parentKsuid := common.String(cmd, "parent-ksuid")
			resourceKsuid := common.String(cmd, "resource-ksuid")
			response, err := c.WorkflowRuns.ReadWorkflowRunByKsuid(ctx, org, parentKsuid, resourceKsuid)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("parent-ksuid", "", "KSUID of the parent workflow.")
	_ = cmd.MarkFlagRequired("parent-ksuid")
	cmd.Flags().String("resource-ksuid", "", "KSUID of the workflow run.")
	_ = cmd.MarkFlagRequired("resource-ksuid")
	return cmd
}
