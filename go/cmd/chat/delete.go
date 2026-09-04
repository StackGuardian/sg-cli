// Code generated for sg-cli; edit the generator spec or this file directly.

package chat

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Archive a chat",
		Long:  `Archive a chat`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			chatId := common.String(cmd, "chat-id")
			response, err := c.Chats.DeleteChat(ctx, org, chatId)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Chat archived successfully.")
			return nil
		},
	}
	cmd.Flags().String("chat-id", "", "The chat ID.")
	_ = cmd.MarkFlagRequired("chat-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
