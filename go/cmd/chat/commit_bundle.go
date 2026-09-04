// Code generated for sg-cli; edit the generator spec or this file directly.

package chat

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/chats"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newCommitBundleCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit-bundle",
		Short: "Commit an uploaded git bundle",
		Long:  `Commit an uploaded git bundle`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			chatId := common.String(cmd, "chat-id")
			expectedVersion, _ := cmd.Flags().GetInt("expected-version")
			response, err := c.Chats.CommitChatBundle(ctx, org, chatId, &chats.CommitChatBundleRequest{ExpectedVersion: expectedVersion})
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Bundle committed.")
			return nil
		},
	}
	cmd.Flags().String("chat-id", "", "The chat ID.")
	_ = cmd.MarkFlagRequired("chat-id")
	cmd.Flags().Int("expected-version", 0, "The BundleVersion last seen by the caller.")
	_ = cmd.MarkFlagRequired("expected-version")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
