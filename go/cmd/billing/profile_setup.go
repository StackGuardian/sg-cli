// Code generated for sg-cli; edit the generator spec or this file directly.

package billing

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/billing"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newProfileSetupCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup -- payload.json",
		Short: "Create the billing profile (once per organization)",
		Long:  `Create the billing profile (once per organization)`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			var req billing.BillingProfileRequest
			if len(args) > 0 {
				if err := common.ReadPayload(cmd, args, &req); err != nil {
					if errors.Is(err, common.ErrDryRun) {
						return nil
					}
					return err
				}
			}
			response, err := c.Billing.SetupBillingProfile(ctx, org, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Billing profile created.")
			return nil
		},
	}
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
