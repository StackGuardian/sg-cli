// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflow

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newArtifactsLockCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Lock the stack workflow on an artifact",
		Long:  `Lock the stack workflow on an artifact`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			workflowId := common.String(cmd, "workflow-id")
			artifactId := common.String(cmd, "artifact-id")
			response, err := c.StackWorkflows.LockStackWorkflowArtifact(ctx, org, stackId, workflowId, workflowGroup, artifactId)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Artifact locked.")
			return nil
		},
	}
	cmd.Flags().String("artifact-id", "", "The artifact ID (path).")
	_ = cmd.MarkFlagRequired("artifact-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
