// Code generated for sg-cli; edit the generator spec or this file directly.

package policy

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get details of a policy",
		Long:  `Get details of a policy`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			policyId := common.String(cmd, "policy-id")
			response, err := c.Policies.ReadPolicy(ctx, org, policyId)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("policy-id", "", "The policy ID.")
	_ = cmd.MarkFlagRequired("policy-id")
	return cmd
}
