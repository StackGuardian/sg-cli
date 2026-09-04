// Code generated for sg-cli; edit the generator spec or this file directly.

package chat

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/chats"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newArtifactsListCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Get a presigned URL that lists the artifacts",
		Long:  `Get a presigned URL that lists the artifacts`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			chatId := common.String(cmd, "chat-id")
			response, err := c.Chats.ListAllChatArtifacts(ctx, org, chatId, &chats.ListAllChatArtifactsRequest{Prefix: common.OptString(cmd, "prefix"), MaxKeys: common.OptInt(cmd, "max-keys"), StartAfterKey: common.OptString(cmd, "start-after-key")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("prefix", "", "Only list artifacts under this prefix.")
	cmd.Flags().Int("max-keys", 0, "Maximum number of keys.")
	cmd.Flags().String("start-after-key", "", "Full S3 key to start after.")
	return cmd
}
