// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflow

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/workflows"
	"github.com/spf13/cobra"
)

func newFileUploadUrlCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file-upload-url",
		Short: "Get a signed URL to upload a file to a stack workflow",
		Long:  `Get a signed URL to upload a file to a stack workflow`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			stackId := common.String(cmd, "stack-id")
			workflowId := common.String(cmd, "workflow-id")
			filename := common.String(cmd, "filename")
			response, err := c.StackWorkflows.GetFileUploadUrlForStackWorkflow(ctx, org, stackId, workflowId, workflowGroup, &workflows.GetFileUploadUrlRequest{Filename: filename, Folder: common.OptString(cmd, "folder")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("workflow-id", "", "The workflow ID.")
	_ = cmd.MarkFlagRequired("workflow-id")
	cmd.Flags().String("filename", "", "Name of the file to upload.")
	_ = cmd.MarkFlagRequired("filename")
	cmd.Flags().String("folder", "", "Folder to upload the file into.")
	return cmd
}
