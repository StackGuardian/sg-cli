package gitscan

import (
	"github.com/StackGuardian/sg-cli/cmd/gitscan/scan"
	"github.com/spf13/cobra"
)

func NewGitScanCmd() *cobra.Command {
	gitScanCmd := &cobra.Command{
		Use:   "git-scan",
		Short: "Scan GitHub or GitLab organizations for Terraform repositories",
		Long:  "Scan GitHub or GitLab organizations for Terraform repositories and generate an sg-payload.json for bulk workflow creation.",
	}

	gitScanCmd.AddCommand(scan.NewScanCmd())

	return gitScanCmd
}
