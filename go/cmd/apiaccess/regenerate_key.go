// Code generated for sg-cli; edit the generator spec or this file directly.

package apiaccess

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newRegenerateKeyCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "regenerate-key",
		Short: "Regenerate the key of an API access",
		Long:  `Regenerate the key of an API access`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			accessId := common.String(cmd, "access-id")
			noExpiration := common.Bool(cmd, "no-expiration")
			var body sggosdk.ApiAccessRegenerate
			if v := common.OptInt(cmd, "expires-at"); v != nil {
				body.ExpiresAt = sggosdk.Optional(*v)
			}
			if noExpiration {
				body.ExpiresAt = sggosdk.Null[int]()
			}
			response, err := c.AccessManagement.RegenerateApiKey(ctx, accessId, org, &body)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "API key regenerated.")
			return nil
		},
	}
	cmd.Flags().String("access-id", "", "The API access ID.")
	_ = cmd.MarkFlagRequired("access-id")
	cmd.Flags().Int("expires-at", 0, "New expiration as a Unix timestamp in milliseconds.")
	cmd.Flags().Bool("no-expiration", false, "Remove the expiration.")
	common.AddOutputJSONFlag(cmd)
	cmd.MarkFlagsMutuallyExclusive("expires-at", "no-expiration")
	return cmd
}
