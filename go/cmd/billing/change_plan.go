// Code generated for sg-cli; edit the generator spec or this file directly.

package billing

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/billing"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newChangePlanCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-plan",
		Short: "Move the organization to another billing package",
		Long:  `Move the organization to another billing package`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.Billing.ChangePlan(ctx, org, &billing.ChangePlanRequest{TargetPackageId: common.OptString(cmd, "target-package-id"), TargetPackage: common.OptString(cmd, "target-package")})
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Plan change submitted.")
			return nil
		},
	}
	cmd.Flags().String("target-package-id", "", "Id of the target package.")
	cmd.Flags().String("target-package", "", "Name of the target package, e.g. \"SG | Standard\".")
	common.AddOutputJSONFlag(cmd)
	cmd.MarkFlagsOneRequired("target-package-id", "target-package")
	cmd.MarkFlagsMutuallyExclusive("target-package-id", "target-package")
	return cmd
}
