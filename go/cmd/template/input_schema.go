// Code generated for sg-cli; edit the generator spec or this file directly.

package template

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/templates"
	"github.com/spf13/cobra"
)

func newInputSchemaCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "input-schema",
		Short: "Decode one of a template's input schemas (requires a user token)",
		Long:  `Decode one of a template's input schemas (requires a user token)`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			templateId := common.String(cmd, "template-id")
			templateType := common.String(cmd, "template-type")
			ownerOrg := common.String(cmd, "owner-org")
			schemaType := common.String(cmd, "schema-type")
			if err := common.CheckEnum("template-type", templateType, "IAC", "IAC_GROUP", "IAC_POLICY", "WORKFLOW_STEP"); err != nil {
				return err
			}
			if ownerOrg == "" {
				ownerOrg = org
			}
			response, err := c.Templates.GetTemplateInputSchema(ctx, ownerOrg, templateType, templateId, &templates.GetTemplateInputSchemaRequest{SchemaType: schemaType})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("template-id", "", "Template name, optionally with revision as NAME:REVISION.")
	_ = cmd.MarkFlagRequired("template-id")
	cmd.Flags().String("template-type", "", "Template type (IAC, IAC_GROUP, IAC_POLICY, WORKFLOW_STEP).")
	_ = cmd.MarkFlagRequired("template-type")
	cmd.Flags().String("owner-org", "", "Organization that owns the template (defaults to --org).")
	cmd.Flags().String("schema-type", "", "Schema to decode (FORM_JSONSCHEMA, RAW_JSON, NO_CODE_JSON, TIRITH_JSON).")
	_ = cmd.MarkFlagRequired("schema-type")
	return cmd
}
