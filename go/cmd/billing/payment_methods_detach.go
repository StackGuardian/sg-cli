// Code generated for sg-cli; edit the generator spec or this file directly.

package billing

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/billing"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newPaymentMethodsDetachCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach",
		Short: "Remove a card",
		Long:  `Remove a card`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			paymentMethodId := common.String(cmd, "payment-method-id")
			response, err := c.Billing.DetachCard(ctx, org, &billing.PaymentMethodRequest{PaymentMethodId: paymentMethodId})
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Card detached.")
			return nil
		},
	}
	cmd.Flags().String("payment-method-id", "", "Stripe payment method id.")
	_ = cmd.MarkFlagRequired("payment-method-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
