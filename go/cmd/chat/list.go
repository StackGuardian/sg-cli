// Code generated for sg-cli; edit the generator spec or this file directly.

package chat

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/chats"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newListCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List chats",
		Long:  `List chats`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.Chats.ListAllChats(ctx, org, &chats.ListAllChatsRequest{Limit: common.OptInt(cmd, "limit"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key"), SearchQuery: common.OptString(cmd, "search-query"), RefreshPrStatus: common.OptBool(cmd, "refresh-pr-status"), ChatType: common.OptString(cmd, "chat-type")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().Int("limit", 0, "Limit the number of results returned.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	cmd.Flags().String("search-query", "", "Comma-separated terms matched against name, description and tags.")
	cmd.Flags().Bool("refresh-pr-status", false, "Refresh pull request statuses from the VCS provider.")
	cmd.Flags().String("chat-type", "", "Filter by chat type (infra2code, kyro).")
	return cmd
}
