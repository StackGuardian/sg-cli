// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflow

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newArtifactsReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Read an artifact's content (or a signed URL on private runners)",
		Long:  `Read an artifact's content (or a signed URL on private runners)`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			workflowId := common.String(cmd, "workflow-id")
			artifactId := common.String(cmd, "artifact-id")
			response, err := c.StackWorkflows.ReadStackWorkflowArtifact(ctx, org, stackId, workflowId, workflowGroup, artifactId)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("artifact-id", "", "The artifact ID (path).")
	_ = cmd.MarkFlagRequired("artifact-id")
	return cmd
}
