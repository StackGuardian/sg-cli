// Code generated for sg-cli; edit the generator spec or this file directly.

package chat

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newPrSyncStatusCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pr-sync-status",
		Short: "Compare the local bundle with the pull request",
		Long:  `Compare the local bundle with the pull request`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			chatId := common.String(cmd, "chat-id")
			response, err := c.Chats.GetChatPrSyncStatus(ctx, org, chatId)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("chat-id", "", "The chat ID.")
	_ = cmd.MarkFlagRequired("chat-id")
	return cmd
}
