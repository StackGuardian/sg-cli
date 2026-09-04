// Code generated for sg-cli; edit the generator spec or this file directly.

package extra

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/workflows"
	"github.com/spf13/cobra"
)

func newGetUrlCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-url",
		Short: "Get a signed download URL for an artifact",
		Long:  `Get a signed download URL for an artifact`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			workflowId := common.String(cmd, "workflow-id")
			artifactPath := common.String(cmd, "artifact-path")
			response, err := c.Workflows.GetWorkflowArtifactUrl(ctx, org, workflowId, workflowGroup, &workflows.GetWorkflowArtifactRequest{ArtifactPath: artifactPath, VersionId: common.OptString(cmd, "version-id")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("artifact-path", "", "Path of the artifact.")
	_ = cmd.MarkFlagRequired("artifact-path")
	cmd.Flags().String("version-id", "", "Specific version of the artifact.")
	return cmd
}
