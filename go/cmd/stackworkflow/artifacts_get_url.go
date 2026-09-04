// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflow

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/workflows"
	"github.com/spf13/cobra"
)

func newArtifactsGetUrlCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-url",
		Short: "Get a signed download URL for an artifact",
		Long:  `Get a signed download URL for an artifact`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			workflowId := common.String(cmd, "workflow-id")
			artifactPath := common.String(cmd, "artifact-path")
			response, err := c.StackWorkflows.GetStackWorkflowArtifactUrl(ctx, org, stackId, workflowId, workflowGroup, &workflows.GetWorkflowArtifactRequest{ArtifactPath: artifactPath})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("artifact-path", "", "Path of the artifact.")
	_ = cmd.MarkFlagRequired("artifact-path")
	return cmd
}
