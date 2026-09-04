// Code generated for sg-cli; edit the generator spec or this file directly.

package billing

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/billing"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newInvoicesCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoices",
		Short: "List invoices",
		Long:  `List invoices`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.Billing.ListInvoices(ctx, org, &billing.ListInvoicesRequest{Limit: common.OptInt(cmd, "limit"), NextPage: common.OptString(cmd, "next-page")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().Int("limit", 0, "Page size (1-100).")
	cmd.Flags().String("next-page", "", "Cursor from the previous response.")
	return cmd
}
