// Code generated for sg-cli; edit the generator spec or this file directly.

package connector

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/connectors"
	"github.com/spf13/cobra"
)

func newListAccountsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-accounts -- payload.json",
		Short: "List the cloud accounts or subscriptions reachable with a connector's credentials",
		Long:  `List the cloud accounts or subscriptions reachable with a connector's credentials`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			var req connectors.ListAccountsRequest
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.Connectors.ListAccounts(ctx, org, &req)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	common.AddPayloadFlags(cmd)
	return cmd
}
