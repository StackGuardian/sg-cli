// Code generated for sg-cli; edit the generator spec or this file directly.

package billing

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewBillingCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "billing",
		Short: "Inspect and manage the organization's billing",
		Long:  `Inspect and manage the organization's billing`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  balance                Show the credit balances
  dashboard-url          Get the URL of the embeddable billing dashboard
  invoices               List invoices
  details                Show the current package and billing information
  change-plan            Move the organization to another billing package
  profile                Manage the Stripe billing profile
  payment-methods        Manage the cards on the billing profile`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newBalanceCmd(c),
		newDashboardUrlCmd(c),
		newInvoicesCmd(c),
		newDetailsCmd(c),
		newChangePlanCmd(c),
		newProfileCmd(c),
		newPaymentMethodsCmd(c),
	)
	return cmd
}
func newProfileCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage the Stripe billing profile",
		Long:  `Manage the Stripe billing profile`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  setup                  Create the billing profile (once per organization)
  update                 Update the billing profile's email, name or address`)
		},
	}
	cmd.AddCommand(
		newProfileSetupCmd(c),
		newProfileUpdateCmd(c),
	)
	return cmd
}
func newPaymentMethodsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "payment-methods",
		Short: "Manage the cards on the billing profile",
		Long:  `Manage the cards on the billing profile`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  list                   List payment methods
  setup-intent           Create a Stripe SetupIntent for adding a card
  set-default            Make a card the default payment method
  detach                 Remove a card`)
		},
	}
	cmd.AddCommand(
		newPaymentMethodsListCmd(c),
		newPaymentMethodsSetupIntentCmd(c),
		newPaymentMethodsSetDefaultCmd(c),
		newPaymentMethodsDetachCmd(c),
	)
	return cmd
}
