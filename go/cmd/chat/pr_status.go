// Code generated for sg-cli; edit the generator spec or this file directly.

package chat

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newPrStatusCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pr-status",
		Short: "Show the status of the chat's pull request",
		Long:  `Show the status of the chat's pull request`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			chatId := common.String(cmd, "chat-id")
			response, err := c.Chats.GetChatPrStatus(ctx, org, chatId)
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
