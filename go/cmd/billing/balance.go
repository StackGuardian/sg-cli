// Code generated for sg-cli; edit the generator spec or this file directly.

package billing

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newBalanceCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance",
		Short: "Show the credit balances",
		Long:  `Show the credit balances`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.Billing.ListBalance(ctx, org)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	return cmd
}
