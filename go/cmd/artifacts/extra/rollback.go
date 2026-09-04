// Code generated for sg-cli; edit the generator spec or this file directly.

package extra

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/workflows"
	"github.com/spf13/cobra"
)

func newRollbackCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "Roll an artifact back to a previous version",
		Long:  `Roll an artifact back to a previous version`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			workflowId := common.String(cmd, "workflow-id")
			artifactId := common.String(cmd, "artifact-id")
			versionId := common.String(cmd, "version-id")
			response, err := c.Workflows.RollbackArtifactVersion(ctx, org, workflowId, workflowGroup, artifactId, &workflows.RollbackArtifactVersionRequest{VersionId: versionId})
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Artifact rolled back.")
			return nil
		},
	}
	cmd.Flags().String("artifact-id", "", "The artifact ID (path).")
	_ = cmd.MarkFlagRequired("artifact-id")
	cmd.Flags().String("version-id", "", "Version to roll back to.")
	_ = cmd.MarkFlagRequired("version-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
