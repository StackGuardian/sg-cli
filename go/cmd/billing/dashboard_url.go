// Code generated for sg-cli; edit the generator spec or this file directly.

package billing

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/billing"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDashboardUrlCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard-url",
		Short: "Get the URL of the embeddable billing dashboard",
		Long:  `Get the URL of the embeddable billing dashboard`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.Billing.GetDashboardUrl(ctx, org, &billing.DashboardRequest{Type: common.OptString(cmd, "type")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("type", "", "Dashboard type (invoices, usage).")
	return cmd
}
