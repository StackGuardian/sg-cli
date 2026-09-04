// Code generated for sg-cli; edit the generator spec or this file directly.

package extra

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newTfstateUploadUrlCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tfstate-upload-url",
		Short: "Get a signed URL to upload a Terraform state file",
		Long:  `Get a signed URL to upload a Terraform state file`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			workflowGroup := common.String(cmd, "workflow-group")
			workflowId := common.String(cmd, "workflow-id")
			response, err := c.Workflows.GetSignedUrlToUploadTfstateFile(ctx, org, workflowId, workflowGroup, &sggosdk.GetSignedUrlToUploadTfstateFileRequest{Filename: common.OptString(cmd, "filename")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("workflow-id", "", "The workflow ID.")
	_ = cmd.MarkFlagRequired("workflow-id")
	cmd.Flags().String("filename", "", "Name of the file to upload (default tfstate.json).")
	return cmd
}
