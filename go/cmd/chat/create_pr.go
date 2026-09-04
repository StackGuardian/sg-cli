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

func newCreatePrCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pr -- payload.json",
		Short: "Push the generated code and open a pull request",
		Long:  `Push the generated code and open a pull request`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			chatId := common.String(cmd, "chat-id")
			var req chats.CreateChatPrRequest
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.Chats.CreateChatPr(ctx, org, chatId, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Pull request request completed.")
			return nil
		},
	}
	cmd.Flags().String("chat-id", "", "The chat ID.")
	_ = cmd.MarkFlagRequired("chat-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
