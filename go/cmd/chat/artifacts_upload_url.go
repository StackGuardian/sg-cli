// Code generated for sg-cli; edit the generator spec or this file directly.

package chat

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/chats"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newArtifactsUploadUrlCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-url",
		Short: "Get a presigned upload URL for an artifact",
		Long:  `Get a presigned upload URL for an artifact`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			chatId := common.String(cmd, "chat-id")
			artifactPath := common.String(cmd, "artifact-path")
			response, err := c.Chats.SaveChatArtifact(ctx, org, chatId, &chats.SaveChatArtifactRequest{ArtifactPath: artifactPath, ContentType: common.OptString(cmd, "content-type")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("artifact-path", "", "Artifact path relative to the chat's artifact store.")
	_ = cmd.MarkFlagRequired("artifact-path")
	cmd.Flags().String("content-type", "", "Content type the upload must use (default application/octet-stream).")
	return cmd
}
