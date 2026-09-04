// Code generated for sg-cli; edit the generator spec or this file directly.

package chat

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/chats"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newArtifactsGetUrlCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-url",
		Short: "Get a presigned download URL for an artifact",
		Long:  `Get a presigned download URL for an artifact`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			chatId := common.String(cmd, "chat-id")
			artifactPath := common.String(cmd, "artifact-path")
			response, err := c.Chats.GetChatArtifact(ctx, org, chatId, &chats.GetChatArtifactRequest{ArtifactPath: artifactPath})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("artifact-path", "", "Artifact path relative to the chat's artifact store.")
	_ = cmd.MarkFlagRequired("artifact-path")
	return cmd
}
