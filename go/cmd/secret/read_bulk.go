// Code generated for sg-cli; edit the generator spec or this file directly.

package secret

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/secrets"
	"github.com/spf13/cobra"
)

func newReadBulkCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read-bulk",
		Short: "Read several secrets by name",
		Long:  `Read several secrets by name`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			names := common.OptStringSlice(cmd, "names")
			response, err := c.Secrets.ReadSecretsInBulk(ctx, org, &secrets.ReadSecretsInBulkRequest{ResourceNames: names})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().StringSlice("names", nil, "Secret names (comma-separated).")
	_ = cmd.MarkFlagRequired("names")
	return cmd
}
