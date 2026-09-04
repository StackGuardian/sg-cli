// Code generated for sg-cli; edit the generator spec or this file directly.

package connector

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/connectors"
	"github.com/spf13/cobra"
)

func newReposCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "List the repositories reachable through a VCS connector",
		Long:  `List the repositories reachable through a VCS connector`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			connectorId := common.String(cmd, "connector-id")
			response, err := c.Connectors.ListRepositories(ctx, connectorId, org, &connectors.ListRepositoriesRequest{Login: common.OptString(cmd, "login"), InstallationId: common.OptString(cmd, "installation-id"), Page: common.OptString(cmd, "page"), Limit: common.OptInt(cmd, "limit"), SearchQuery: common.OptString(cmd, "search-query")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("connector-id", "", "The connector ID.")
	_ = cmd.MarkFlagRequired("connector-id")
	cmd.Flags().String("login", "", "GitHub account or organization login (GITHUB_COM).")
	cmd.Flags().String("installation-id", "", "GitHub App installation id (GITHUB_COM).")
	cmd.Flags().String("page", "", "Page number or provider cursor.")
	cmd.Flags().Int("limit", 0, "Page size.")
	cmd.Flags().String("search-query", "", "Filter by name or namespace.")
	return cmd
}
