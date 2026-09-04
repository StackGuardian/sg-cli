// Code generated for sg-cli; edit the generator spec or this file directly.

package benchmarkreport

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get the benchmark report",
		Long:  `Get the benchmark report`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.BenchmarkReports.GetBenchmarkReports(ctx, org, &sggosdk.GetBenchmarkReportsRequest{Date: common.OptString(cmd, "date"), Detailed: common.OptBool(cmd, "detailed"), FilterCsp: common.OptString(cmd, "filter-csp"), FilterAccountId: common.OptString(cmd, "filter-account-id"), FilterBenchmark: common.OptString(cmd, "filter-benchmark"), FilterControlId: common.OptString(cmd, "filter-control-id"), FilterControlTitle: common.OptString(cmd, "filter-control-title"), FilterRegion: common.OptString(cmd, "filter-region"), FilterSeverity: common.OptString(cmd, "filter-severity"), FilterStatus: common.OptString(cmd, "filter-status"), GroupBy: common.OptString(cmd, "group-by"), Limit: common.OptInt(cmd, "limit"), Page: common.OptInt(cmd, "page"), RequiredColumns: common.OptString(cmd, "required-columns")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("date", "", "Report date in DD_MM_YYYY format (defaults to today).")
	cmd.Flags().Bool("detailed", false, "Return detailed results instead of grouped data.")
	cmd.Flags().String("filter-csp", "", "Filter by cloud provider (AWS, AZURE, GCP).")
	cmd.Flags().String("filter-account-id", "", "Filter by account id.")
	cmd.Flags().String("filter-benchmark", "", "Filter by benchmark, e.g. cis_v200.")
	cmd.Flags().String("filter-control-id", "", "Filter by control id.")
	cmd.Flags().String("filter-control-title", "", "Filter by control title.")
	cmd.Flags().String("filter-region", "", "Filter by region.")
	cmd.Flags().String("filter-severity", "", "Filter by severity.")
	cmd.Flags().String("filter-status", "", "Filter by status (fails, passes, info, skips, error).")
	cmd.Flags().String("group-by", "", "Column(s) to group by (required when not --detailed).")
	cmd.Flags().Int("limit", 0, "Maximum number of records.")
	cmd.Flags().Int("page", 0, "Page number (requires --limit).")
	cmd.Flags().String("required-columns", "", "Columns to include in a detailed report (comma-separated).")
	return cmd
}
