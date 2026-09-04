// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/templates"
	"github.com/spf13/cobra"
)

func newSubscriptionsUpdateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update -- payload.json",
		Short: "Subscribe to or unsubscribe from templates",
		Long:  `Subscribe to or unsubscribe from templates`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			subscription := common.String(cmd, "subscription")
			action := common.String(cmd, "action")
			var req templates.UpdateSubscriptionRequest
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			if err := common.CheckEnum("action", action, "subscribe", "unsubscribe"); err != nil {
				return err
			}
			response, err := c.Templates.UpdateSubscription(ctx, org, subscription, action, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Subscription updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("subscription", "default", "The subscription name.")
	cmd.Flags().String("action", "", "subscribe or unsubscribe.")
	_ = cmd.MarkFlagRequired("action")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
