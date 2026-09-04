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

func newMessagesCreateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create -- payload.json",
		Short: "Send a message (use --text or a payload)",
		Long:  `Send a message (use --text or a payload)`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			chatId := common.String(cmd, "chat-id")
			text := common.String(cmd, "text")
			var req chats.Message
			if len(args) > 0 {
				if err := common.ReadPayload(cmd, args, &req); err != nil {
					if errors.Is(err, common.ErrDryRun) {
						return nil
					}
					return err
				}
			}
			if text != "" {
				req.Content = []*chats.ContentBlock{{Type: "text", Text: &text}}
			}
			if v := common.OptString(cmd, "client-request-id"); v != nil {
				req.ClientRequestId = v
			}
			if len(req.Content) == 0 {
				return errors.New("provide --text or a payload with content blocks")
			}
			response, err := c.Chats.CreateMessage(ctx, org, chatId, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Message sent.")
			return nil
		},
	}
	cmd.Flags().String("text", "", "Text of the message (shortcut for a single text content block).")
	cmd.Flags().String("client-request-id", "", "Idempotency key.")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
