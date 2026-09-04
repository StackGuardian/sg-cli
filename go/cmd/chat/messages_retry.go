// Code generated for sg-cli; edit the generator spec or this file directly.

package chat

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newMessagesRetryCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retry",
		Short: "Retry an errored or stopped assistant message",
		Long:  `Retry an errored or stopped assistant message`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			chatId := common.String(cmd, "chat-id")
			messageId := common.String(cmd, "message-id")
			response, err := c.Chats.RetryMessage(ctx, org, chatId, messageId)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Message retried.")
			return nil
		},
	}
	cmd.Flags().String("message-id", "", "The message ID.")
	_ = cmd.MarkFlagRequired("message-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
