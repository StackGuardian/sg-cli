// Code generated for sg-cli; edit the generator spec or this file directly.

package secret

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a secret",
		Long:  `Delete a secret`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			secretId := common.String(cmd, "secret-id")
			response, err := c.Secrets.DeleteSecret(ctx, org, secretId)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Secret deleted successfully.")
			return nil
		},
	}
	cmd.Flags().String("secret-id", "", "The secret name.")
	_ = cmd.MarkFlagRequired("secret-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
