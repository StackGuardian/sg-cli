// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newCreateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create -- payload.json",
		Short: "Create a template or a new revision of it (payload must include TemplateType: IAC, IAC_GROUP or IAC_POLICY)",
		Long:  `Create a template or a new revision of it (payload must include TemplateType: IAC, IAC_GROUP or IAC_POLICY)`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			var req sggosdk.CreateTemplateRequest
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.Templates.CreateTemplateRevision(ctx, &sggosdk.CreateTemplateRevisionRequest{SgOrgid: org, Body: &req})
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Template created successfully.")
			return nil
		},
	}
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
