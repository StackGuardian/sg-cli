// Code generated for sg-cli; edit the generator spec or this file directly.

package chat

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newMessagesCancelCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Stop a pending or streaming assistant message",
		Long:  `Stop a pending or streaming assistant message`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			chatId := common.String(cmd, "chat-id")
			messageId := common.String(cmd, "message-id")
			response, err := c.Chats.CancelMessage(ctx, org, chatId, messageId)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Message cancelled.")
			return nil
		},
	}
	cmd.Flags().String("message-id", "", "The message ID.")
	_ = cmd.MarkFlagRequired("message-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
