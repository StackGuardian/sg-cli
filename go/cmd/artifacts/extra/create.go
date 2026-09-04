// Code generated for sg-cli; edit the generator spec or this file directly.

package extra

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newCreateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create -- payload.json",
		Short: "Create or replace an artifact with the JSON payload as its content",
		Long:  `Create or replace an artifact with the JSON payload as its content`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			workflowId := common.String(cmd, "workflow-id")
			artifactId := common.String(cmd, "artifact-id")
			var req interface{}
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.Workflows.CreateArtifact(ctx, org, workflowId, workflowGroup, artifactId, req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Artifact created successfully.")
			return nil
		},
	}
	cmd.Flags().String("artifact-id", "", "The artifact ID (path).")
	_ = cmd.MarkFlagRequired("artifact-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
