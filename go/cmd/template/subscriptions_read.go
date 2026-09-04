// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newSubscriptionsReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get the subscriptions of a type",
		Long:  `Get the subscriptions of a type`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			subscriptionType := common.String(cmd, "subscription-type")
			if err := common.CheckEnum("subscription-type", subscriptionType, "IACGroupSubscriptions", "IACSubscriptions", "PolicySubscriptions", "WfStepSubscriptions"); err != nil {
				return err
			}
			response, err := c.Templates.ReadSubscription(ctx, org, &sggosdk.ReadSubscriptionRequest{SubscriptionType: sggosdk.ReadSubscriptionRequestSubscriptionType(subscriptionType)})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("subscription-type", "", "Subscription type (IACGroupSubscriptions, IACSubscriptions, PolicySubscriptions, WfStepSubscriptions).")
	_ = cmd.MarkFlagRequired("subscription-type")
	return cmd
}
