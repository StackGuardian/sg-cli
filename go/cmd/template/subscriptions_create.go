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

func newSubscriptionsCreateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create -- payload.json",
		Short: "Create a subscription set (all four subscription maps are required)",
		Long:  `Create a subscription set (all four subscription maps are required)`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			var req templates.CreateSubscriptionRequest
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.Templates.CreateSubscription(ctx, org, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Subscription created successfully.")
			return nil
		},
	}
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
