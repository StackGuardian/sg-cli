// Code generated for sg-cli; edit the generator spec or this file directly.

package chat

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/chats"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newUpdateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update -- payload.json",
		Short: "Update a chat",
		Long:  `Update a chat`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			chatId := common.String(cmd, "chat-id")
			var req chats.PatchedChat
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.Chats.UpdateChat(ctx, org, chatId, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Chat updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("chat-id", "", "The chat ID.")
	_ = cmd.MarkFlagRequired("chat-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
