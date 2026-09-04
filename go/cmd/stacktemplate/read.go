// Code generated for sg-cli; edit the generator spec or this file directly.

package stacktemplate

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get a stack template",
		Long:  `Get a stack template`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			templateId := common.String(cmd, "template-id")
			response, err := c.StackTemplates.ReadStackTemplate(ctx, org, templateId)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("template-id", "", "The template ID.")
	_ = cmd.MarkFlagRequired("template-id")
	return cmd
}
