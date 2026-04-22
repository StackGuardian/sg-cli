// Package interactive implements the full menu-driven TUI mode for sg-cli.
// Launch with: sg-cli interactive  (alias: sg-cli i)
package interactive

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/StackGuardian/sg-cli/cmd/gitscan/scan"
	wfcreate "github.com/StackGuardian/sg-cli/cmd/workflow/create"
	"github.com/StackGuardian/sg-cli/cmd/output"
	"github.com/StackGuardian/sg-cli/cmd/tui"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const dashboardURL = "https://app.stackguardian.io/orchestrator"

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// ---------------------------------------------------------------------------
// Result card styles
// ---------------------------------------------------------------------------

var (
	cardStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")).
			Padding(1, 3).
			MarginTop(1).
			MarginBottom(1)

	cardLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Width(14)

	cardValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F9FAFB")).
			Bold(true)

	cardSuccessTitle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#22C55E")).
				Bold(true)

	cardErrorTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true)

	cardURLStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60A5FA")).
			Underline(true)

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true)
)

func kv(label, value string) string {
	return cardLabelStyle.Render(label+":") + " " + cardValueStyle.Render(value)
}

// showResultCard renders a styled result box after an action completes.
func showResultCard(title, resource, org, wfGrp, runURL string, succeeded bool) {
	clearScreen()
	titleStr := cardSuccessTitle.Render("✓  " + title)
	if !succeeded {
		titleStr = cardErrorTitle.Render("✗  " + title)
	}

	var lines []string
	lines = append(lines, titleStr)
	lines = append(lines, "")
	lines = append(lines, kv("Resource", resource))
	lines = append(lines, kv("Org", org))
	lines = append(lines, kv("Wf Group", wfGrp))
	if runURL != "" {
		lines = append(lines, "")
		lines = append(lines, cardLabelStyle.Render("Dashboard:")+" "+cardURLStyle.Render(runURL))
	}

	fmt.Println(cardStyle.Render(strings.Join(lines, "\n")))
	fmt.Println(hintStyle.Render("  Press enter to return to menu"))
	fmt.Println()

	// Wait for enter
	_ = huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("").
			Affirmative("Back to menu").
			Negative("Exit").
			Value(new(bool)),
	)).Run()
}

// showErrorCard renders an error in a styled card and waits for the user to dismiss it.
func showErrorCard(msg string) {
	clearScreen()
	lines := []string{
		cardErrorTitle.Render("✗  Error"),
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FCA5A5")).Render(msg),
	}
	fmt.Println(cardStyle.Render(strings.Join(lines, "\n")))
	fmt.Println()
	_ = huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("").
			Affirmative("Back to menu").
			Negative("Exit").
			Value(new(bool)),
	)).Run()
}

// showDeletedCard renders a minimal card for delete confirmations.
func showDeletedCard(resource, resourceType string) {
	clearScreen()
	lines := []string{
		cardSuccessTitle.Render("✓  " + resourceType + " deleted"),
		"",
		kv("Resource", resource),
	}
	fmt.Println(cardStyle.Render(strings.Join(lines, "\n")))
	fmt.Println(hintStyle.Render("  Press enter to return to menu"))
	fmt.Println()
	_ = huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("").
			Affirmative("Back to menu").
			Negative("Exit").
			Value(new(bool)),
	)).Run()
}

// ---------------------------------------------------------------------------
// Command registration
// ---------------------------------------------------------------------------

func NewInteractiveCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "interactive",
		Aliases: []string{"i"},
		Short:   "Launch the interactive menu-driven TUI",
		Long:    "Navigate all sg-cli commands through a guided interactive menu — no flags required.",
		Run: func(cmd *cobra.Command, args []string) {
			run(c)
		},
	}
	return cmd
}

// ---------------------------------------------------------------------------
// Session state
// ---------------------------------------------------------------------------

type session struct {
	client *client.Client
	org    string
	wfGrp  string
}

// ---------------------------------------------------------------------------
// Entry point
// ---------------------------------------------------------------------------

func run(c *client.Client) {
	fmt.Println(output.Banner("2.1.0"))
	fmt.Println()
	output.Info("Interactive mode — press Esc or q at any menu to go back / exit.")
	fmt.Println()

	s := &session{client: c}

	for {
		subtitle := ""
		if s.org != "" && s.wfGrp != "" {
			subtitle = contextLabel(s.org, s.wfGrp)
		}

		action, err := tui.NewPicker("What would you like to manage?", subtitle, []tui.Item{
			{ID: "workflow",  Label: "workflow",       Description: "Create, run, and manage workflows"},
			{ID: "stack",     Label: "stack",          Description: "Create, run, and manage stacks"},
			{ID: "artifacts", Label: "artifacts",      Description: "List workflow artifacts"},
			{ID: "git-scan",  Label: "git-scan",       Description: "Scan GitHub / GitLab for Terraform repos"},
			{ID: "switch",    Label: "switch context", Description: "Change org / workflow group"},
			{ID: "exit",      Label: "exit",           Description: "Quit interactive mode"},
		})
		if err != nil || action == "exit" {
			fmt.Println()
			output.Info("Goodbye.")
			return
		}

		if action == "switch" {
			s.clearContext()
			continue
		}

		switch action {
		case "workflow":
			s.workflowMenu()
		case "stack":
			s.stackMenu()
		case "artifacts":
			s.artifactsMenu()
		case "git-scan":
			s.gitScanMenu()
		}
	}
}

// ensureContext prompts for --org and --workflow-group if not already set.
func (s *session) ensureContext() error {
	if s.org != "" && s.wfGrp != "" {
		return nil
	}

	clearScreen()
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED")).Bold(true).Render("Organisation & Workflow Group"))
	fmt.Println()

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Organisation").
				Description("Your StackGuardian organisation name").
				Validate(func(v string) error {
					if strings.TrimSpace(v) == "" {
						return fmt.Errorf("organisation is required")
					}
					return nil
				}).
				Value(&s.org),

			huh.NewInput().
				Title("Workflow Group").
				Description("The workflow group to operate on").
				Validate(func(v string) error {
					if strings.TrimSpace(v) == "" {
						return fmt.Errorf("workflow group is required")
					}
					return nil
				}).
				Value(&s.wfGrp),
		),
	).Run()
}

func (s *session) clearContext() {
	s.org = ""
	s.wfGrp = ""
}

// ---------------------------------------------------------------------------
// Workflow menu
// ---------------------------------------------------------------------------

func (s *session) workflowMenu() {
	for {
		if err := s.ensureContext(); err != nil {
			return
		}
		action, err := tui.NewPicker("Workflow", contextLabel(s.org, s.wfGrp), []tui.Item{
			{ID: "list",    Label: "list",    Description: "Browse all workflows"},
			{ID: "apply",   Label: "apply",   Description: "Trigger an Apply run"},
			{ID: "destroy", Label: "destroy", Description: "Trigger a Destroy run"},
			{ID: "delete",  Label: "delete",  Description: "Permanently delete a workflow"},
			{ID: "read",    Label: "read",    Description: "View workflow details"},
			{ID: "create",  Label: "create",  Description: "Create a workflow from a JSON payload file"},
			{ID: "back",    Label: "← back",  Description: "Return to main menu"},
		})
		if err != nil || action == "back" {
			return
		}

		switch action {
		case "list":
			s.workflowList()
		case "apply":
			s.workflowAction("apply")
		case "destroy":
			s.workflowAction("destroy")
		case "delete":
			s.workflowDelete()
		case "read":
			s.workflowRead()
		case "create":
			s.workflowCreate()
		}
	}
}

func (s *session) workflowList() {
	workflows, err := s.fetchWorkflows()
	if err != nil {
		showErrorCard(err.Error())
		return
	}
	if len(workflows) == 0 {
		output.Warning("No workflows found.")
		return
	}

	// Build index for quick lookup
	byName := make(map[string]*sggosdk.GeneratedWorkflowsListAllMsg, len(workflows))
	for _, wf := range workflows {
		byName[wf.ResourceName] = wf
	}

	// Keep showing the list until the user escapes
	for {
		items := make([]tui.Item, len(workflows))
		for i, wf := range workflows {
			items[i] = tui.Item{
				ID:          wf.ResourceName,
				Label:       wf.ResourceName,
				Description: wf.Description,
				Badge:       wf.LatestWfrunStatus,
			}
		}

		selected, err := tui.NewPicker("Workflows", contextLabel(s.org, s.wfGrp), items)
		if err != nil {
			// esc / cancelled — go back
			return
		}

		wf := byName[selected]
		desc := wf.Description
		if desc == "" {
			desc = "(no description)"
		}

		clearScreen()
		lines := []string{
			cardSuccessTitle.Render(wf.ResourceName),
			"",
			kv("Type",       wf.WfType),
			kv("Status",     wf.LatestWfrunStatus),
			kv("Org",        s.org),
			kv("Wf Group",   s.wfGrp),
			kv("Resource ID", wf.ResourceId),
			"",
			cardLabelStyle.Render("Description:") + " " + lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")).Render(desc),
		}
		fmt.Println(cardStyle.Render(strings.Join(lines, "\n")))
		fmt.Println()

		goBack := true
		_ = huh.NewForm(huh.NewGroup(
			huh.NewConfirm().
				Title("").
				Affirmative("Back to list").
				Negative("Back to menu").
				Value(&goBack),
		)).Run()

		if !goBack {
			return
		}
	}
}

func (s *session) workflowAction(action string) {
	title := strings.Title(action) + " Workflow"
	wfId, err := s.pickWorkflow(title)
	if err != nil {
		if err.Error() != "cancelled" {
			showErrorCard(err.Error())
		}
		return
	}

	// Confirm before executing
	confirmed := false
	_ = huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title(fmt.Sprintf("Run %s on  %s?", strings.ToUpper(action), wfId)).
			Description(fmt.Sprintf("Org: %s   Wf Group: %s", s.org, s.wfGrp)).
			Affirmative("Yes, run it").
			Negative("Cancel").
			Value(&confirmed),
	)).Run()

	if !confirmed {
		output.Warning("Cancelled.")
		return
	}

	var apiErr error
	switch action {
	case "apply":
		apiErr = output.WithSpinner("Running apply on "+wfId+"...", func() error {
			_, e := s.client.WorkflowRuns.CreateWorkflowRun(
				context.Background(), s.org, wfId, s.wfGrp,
				&sggosdk.WorkflowRun{TerraformAction: &sggosdk.TerraformAction{Action: sggosdk.ActionEnumApply.Ptr()}},
			)
			return e
		})
	case "destroy":
		apiErr = output.WithSpinner("Running destroy on "+wfId+"...", func() error {
			_, e := s.client.WorkflowRuns.CreateWorkflowRun(
				context.Background(), s.org, wfId, s.wfGrp,
				&sggosdk.WorkflowRun{TerraformAction: &sggosdk.TerraformAction{Action: sggosdk.ActionEnumDestroy.Ptr()}},
			)
			return e
		})
	}

	if apiErr != nil {
		showErrorCard(apiErr.Error())
		return
	}

	runURL := dashboardURL + "/orgs/" + s.org + "/wfgrps/" + s.wfGrp + "/wfs/" + wfId + "?tab=runs"
	showResultCard("Workflow "+action+" triggered", wfId, s.org, s.wfGrp, runURL, true)
}

func (s *session) workflowDelete() {
	wfId, err := s.pickWorkflow("Delete Workflow")
	if err != nil {
		if err.Error() != "cancelled" {
			showErrorCard(err.Error())
		}
		return
	}

	confirmed := false
	_ = huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Delete  " + wfId + "?").
			Description("This action is permanent and cannot be undone.").
			Affirmative("Yes, delete it").
			Negative("Cancel").
			Value(&confirmed),
	)).Run()

	if !confirmed {
		output.Warning("Deletion cancelled.")
		return
	}

	err = output.WithSpinner("Deleting "+wfId+"...", func() error {
		_, e := s.client.Workflows.DeleteWorkflow(context.Background(), s.org, wfId, s.wfGrp)
		return e
	})
	if err != nil {
		showErrorCard(err.Error())
		return
	}

	showDeletedCard(wfId, "Workflow")
}

func (s *session) workflowRead() {
	wfId, err := s.pickWorkflow("Read Workflow")
	if err != nil {
		if err.Error() != "cancelled" {
			showErrorCard(err.Error())
		}
		return
	}

	var response *sggosdk.WorkflowGetResponse
	err = output.WithSpinner("Fetching "+wfId+"...", func() error {
		response, err = s.client.Workflows.ReadWorkflow(context.Background(), s.org, wfId, s.wfGrp)
		return err
	})
	if err != nil {
		showErrorCard(err.Error())
		return
	}

	showWorkflowDetail(response, s.org, s.wfGrp)
}

func showWorkflowDetail(response *sggosdk.WorkflowGetResponse, org, wfGrp string) {
	clearScreen()
	wf := response.Msg
	if wf == nil {
		showErrorCard("empty response from API")
		return
	}

	// tags
	tags := "—"
	if len(wf.Tags) > 0 {
		tags = strings.Join(wf.Tags, ", ")
	}

	// status color
	status := wf.LatestWfrunStatus
	if status == "" {
		status = "—"
	}
	statusStyled := lipgloss.NewStyle().
		Foreground(func() lipgloss.Color {
			switch strings.ToUpper(status) {
			case "COMPLETED":
				return lipgloss.Color("#22C55E")
			case "ERRORED":
				return lipgloss.Color("#EF4444")
			case "RUNNING":
				return lipgloss.Color("#3B82F6")
			default:
				return lipgloss.Color("#F59E0B")
			}
		}()).
		Bold(true).
		Render(status)

	divider := lipgloss.NewStyle().Foreground(lipgloss.Color("#374151")).Render(strings.Repeat("─", 44))

	lines := []string{
		cardSuccessTitle.Render("  " + wf.ResourceName),
		"",
		kv("Status",   statusStyled),
		kv("Type",     wf.WfType),
		kv("Org",      org),
		kv("Wf Group", wfGrp),
	}

	if wf.Description != "" {
		lines = append(lines, kv("Description", wf.Description))
	}

	if tags != "—" {
		lines = append(lines, divider, kv("Tags", tags))
	}

	fmt.Println(cardStyle.Render(strings.Join(lines, "\n")))
	fmt.Println()

	// "View raw JSON" or "Back to menu"
	viewJSON := false
	_ = huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("").
			Affirmative("Back to menu").
			Negative("View raw JSON").
			Value(&viewJSON),
	)).Run()

	if !viewJSON {
		raw, err := json.MarshalIndent(wf, "", "  ")
		if err == nil {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")).Render(string(raw)))
			fmt.Println()
			_ = huh.NewForm(huh.NewGroup(
				huh.NewConfirm().Title("").Affirmative("Back to menu").Negative("").Value(new(bool)),
			)).Run()
		}
	}
}

func (s *session) workflowCreate() {
	// Step 1: pick mode
	mode, err := tui.NewPicker("Create Workflow", contextLabel(s.org, s.wfGrp), []tui.Item{
		{ID: "single", Label: "single", Description: "Create one workflow from a JSON payload file"},
		{ID: "bulk",   Label: "bulk",   Description: "Create many workflows from a JSON array file (e.g. from git-scan)"},
		{ID: "back",   Label: "← back", Description: "Return to workflow menu"},
	})
	if err != nil || mode == "back" {
		return
	}

	// Step 2: collect options
	var payloadPath string
	var patchPayload string
	var runAfter bool
	var dryRun bool

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Payload file path").
				Description("Path to your workflow JSON payload file").
				Validate(func(v string) error {
					v = strings.TrimSpace(v)
					if v == "" {
						return fmt.Errorf("payload path is required")
					}
					if _, statErr := os.Stat(v); statErr != nil {
						return fmt.Errorf("file not found: %s", v)
					}
					return nil
				}).
				Value(&payloadPath),

			huh.NewInput().
				Title("Patch payload (optional)").
				Description("JSON string to merge into the payload (leave blank to skip)").
				Value(&patchPayload),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Trigger a run immediately after creation?").
				Value(&runAfter),

			huh.NewConfirm().
				Title("Dry run? (preview only, nothing will be created)").
				Value(&dryRun),
		),
	).Run()
	if err != nil {
		return
	}

	// Step 3: run
	clearScreen()
	fmt.Println()
	opts := &wfcreate.RunOptions{
		Org:          s.org,
		WfgGrp:       s.wfGrp,
		Payload:      payloadPath,
		PatchPayload: patchPayload,
		Run:          runAfter,
		DryRun:       dryRun,
		Bulk:         mode == "bulk",
	}
	wfcreate.RunCreate(s.client, opts)
	fmt.Println()

	_ = huh.NewForm(huh.NewGroup(
		huh.NewConfirm().Title("").Affirmative("Back to menu").Negative("Exit").Value(new(bool)),
	)).Run()
}

// ---------------------------------------------------------------------------
// Stack menu
// ---------------------------------------------------------------------------

func (s *session) stackMenu() {
	for {
		if err := s.ensureContext(); err != nil {
			return
		}
		action, err := tui.NewPicker("Stack", contextLabel(s.org, s.wfGrp), []tui.Item{
			{ID: "apply",   Label: "apply",   Description: "Trigger an Apply run"},
			{ID: "destroy", Label: "destroy", Description: "Trigger a Destroy run"},
			{ID: "delete",  Label: "delete",  Description: "Permanently delete a stack"},
			{ID: "outputs", Label: "outputs", Description: "View stack outputs"},
			{ID: "create",  Label: "create",  Description: "Create a stack from a JSON payload file"},
			{ID: "back",    Label: "← back",  Description: "Return to main menu"},
		})
		if err != nil || action == "back" {
			return
		}

		switch action {
		case "apply":
			s.stackAction("apply")
		case "destroy":
			s.stackAction("destroy")
		case "delete":
			s.stackDelete()
		case "outputs":
			s.stackOutputs()
		case "create":
			s.stackCreate()
		}
	}
}

func (s *session) stackAction(action string) {
	title := strings.Title(action) + " Stack"
	stackId, err := s.pickStack(title)
	if err != nil {
		if err.Error() != "cancelled" {
			showErrorCard(err.Error())
		}
		return
	}

	confirmed := false
	_ = huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title(fmt.Sprintf("Run %s on  %s?", strings.ToUpper(action), stackId)).
			Description(fmt.Sprintf("Org: %s   Wf Group: %s", s.org, s.wfGrp)).
			Affirmative("Yes, run it").
			Negative("Cancel").
			Value(&confirmed),
	)).Run()

	if !confirmed {
		output.Warning("Cancelled.")
		return
	}

	var actionEnum sggosdk.ActionEnum
	if action == "apply" {
		actionEnum = sggosdk.ActionEnumApply
	} else {
		actionEnum = sggosdk.ActionEnumDestroy
	}

	apiErr := output.WithSpinner("Running "+action+" on "+stackId+"...", func() error {
		_, e := s.client.StackRuns.CreateStackRun(
			context.Background(), s.org, stackId, s.wfGrp,
			&sggosdk.StackAction{ActionType: string(actionEnum)},
		)
		return e
	})
	if apiErr != nil {
		showErrorCard(apiErr.Error())
		return
	}

	runURL := dashboardURL + "/orgs/" + s.org + "/wfgrps/" + s.wfGrp + "/stacks/" + stackId + "?tab=runs"
	showResultCard("Stack "+action+" triggered", stackId, s.org, s.wfGrp, runURL, true)
}

func (s *session) stackDelete() {
	stackId, err := s.pickStack("Delete Stack")
	if err != nil {
		if err.Error() != "cancelled" {
			showErrorCard(err.Error())
		}
		return
	}

	confirmed := false
	_ = huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Delete  " + stackId + "?").
			Description("This action is permanent and cannot be undone.").
			Affirmative("Yes, delete it").
			Negative("Cancel").
			Value(&confirmed),
	)).Run()

	if !confirmed {
		output.Warning("Deletion cancelled.")
		return
	}

	err = output.WithSpinner("Deleting stack "+stackId+"...", func() error {
		_, e := s.client.Stacks.DeleteStack(context.Background(), s.org, stackId, s.wfGrp)
		return e
	})
	if err != nil {
		showErrorCard(err.Error())
		return
	}

	showDeletedCard(stackId, "Stack")
}

func (s *session) stackOutputs() {
	stackId, err := s.pickStack("Stack Outputs")
	if err != nil {
		if err.Error() != "cancelled" {
			showErrorCard(err.Error())
		}
		return
	}

	var response interface{}
	err = output.WithSpinner("Fetching outputs for "+stackId+"...", func() error {
		response, err = s.client.Stacks.ReadStackOutputs(context.Background(), s.org, stackId, s.wfGrp)
		return err
	})
	if err != nil {
		showErrorCard(err.Error())
		return
	}

	clearScreen()
	lines := []string{
		cardSuccessTitle.Render("✓  " + stackId + "  —  Outputs"),
		"",
		kv("Org",      s.org),
		kv("Wf Group", s.wfGrp),
	}
	fmt.Println(cardStyle.Render(strings.Join(lines, "\n")))
	fmt.Println(output.Muted.Render("─── Outputs ───────────────────────────────────"))
	fmt.Println(response)
	fmt.Println()

	_ = huh.NewForm(huh.NewGroup(
		huh.NewConfirm().Title("").Affirmative("Back to menu").Negative("Exit").Value(new(bool)),
	)).Run()
}

func (s *session) stackCreate() {
	var payloadPath string

	err := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Payload file path").
			Description("Path to your stack JSON payload file").
			Validate(func(v string) error {
				v = strings.TrimSpace(v)
				if v == "" {
					return fmt.Errorf("payload path is required")
				}
				if _, err := os.Stat(v); err != nil {
					return fmt.Errorf("file not found: %s", v)
				}
				return nil
			}).
			Value(&payloadPath),
	)).Run()
	if err != nil {
		return
	}

	hint := fmt.Sprintf("sg-cli stack create --org %s --workflow-group %s %s", s.org, s.wfGrp, payloadPath)
	lines := []string{
		cardSuccessTitle.Render("✓  Ready to create"),
		"",
		kv("Payload",  payloadPath),
		kv("Org",      s.org),
		kv("Wf Group", s.wfGrp),
		"",
		cardLabelStyle.Render("Command:"),
		"  " + lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED")).Render(hint),
	}
	fmt.Println(cardStyle.Render(strings.Join(lines, "\n")))
	fmt.Println(hintStyle.Render("  Copy and run the command above to execute."))
	fmt.Println()

	_ = huh.NewForm(huh.NewGroup(
		huh.NewConfirm().Title("").Affirmative("Back to menu").Negative("Exit").Value(new(bool)),
	)).Run()
}

// ---------------------------------------------------------------------------
// Artifacts
// ---------------------------------------------------------------------------

func (s *session) artifactsMenu() {
	if err := s.ensureContext(); err != nil {
		return
	}
	wfId, err := s.pickWorkflow("Artifacts — Select Workflow")
	if err != nil {
		if err.Error() != "cancelled" {
			showErrorCard(err.Error())
		}
		return
	}

	var response interface{}
	err = output.WithSpinner("Fetching artifacts for "+wfId+"...", func() error {
		response, err = s.client.Workflows.ListAllWorkflowArtifacts(
			context.Background(), s.org, wfId, s.wfGrp,
		)
		return err
	})
	if err != nil {
		if strings.Contains(err.Error(), "the server responded with nothing") {
			output.Warning("No artifacts found for " + wfId)
		} else {
			showErrorCard(err.Error())
		}
		return
	}

	clearScreen()
	lines := []string{
		cardSuccessTitle.Render("✓  " + wfId + "  —  Artifacts"),
		"",
		kv("Org",      s.org),
		kv("Wf Group", s.wfGrp),
	}
	fmt.Println(cardStyle.Render(strings.Join(lines, "\n")))
	fmt.Println(output.Muted.Render("─── Artifacts ─────────────────────────────────"))
	fmt.Println(response)
	fmt.Println()

	_ = huh.NewForm(huh.NewGroup(
		huh.NewConfirm().Title("").Affirmative("Back to menu").Negative("Exit").Value(new(bool)),
	)).Run()
}

// ---------------------------------------------------------------------------
// Git-scan
// ---------------------------------------------------------------------------

func (s *session) gitScanMenu() {
	var provider string
	var token string
	var target string
	var orgOrUser string
	var maxReposStr string
	var outputPath string
	var managedState bool

	target = "org"
	maxReposStr = "0"
	outputPath = "sg-payload.json"

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("VCS Provider").
				Options(
					huh.NewOption("GitHub", "github"),
					huh.NewOption("GitLab", "gitlab"),
				).
				Value(&provider),

			huh.NewInput().
				Title("Access Token").
				Description("GitHub PAT (ghp_...) or GitLab PAT (glpat-...)").
				Password(true).
				Validate(func(v string) error {
					if strings.TrimSpace(v) == "" {
						return fmt.Errorf("token is required")
					}
					return nil
				}).
				Value(&token),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Scan target").
				Options(
					huh.NewOption("Organisation / Group", "org"),
					huh.NewOption("User", "user"),
				).
				Value(&target),

			huh.NewInput().
				Title("Organisation / User name").
				Description("Leave blank to scan all repos accessible to the token").
				Value(&orgOrUser),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Output file").
				Value(&outputPath),

			huh.NewInput().
				Title("Max repositories").
				Description("0 = no limit").
				Value(&maxReposStr),

			huh.NewConfirm().
				Title("Enable SG-managed Terraform state?").
				Value(&managedState),
		),
	).Run()
	if err != nil {
		return
	}

	maxRepos := 0
	if n, e := strconv.Atoi(strings.TrimSpace(maxReposStr)); e == nil {
		maxRepos = n
	}

	opts := &scan.RunOptions{
		Provider:     provider,
		Token:        token,
		ManagedState: managedState,
		Output:       outputPath,
		MaxRepos:     maxRepos,
	}
	if target == "org" {
		opts.Org = orgOrUser
	} else {
		opts.User = orgOrUser
	}

	clearScreen()
	fmt.Println()
	scan.RunScan(opts)
	fmt.Println()

	_ = huh.NewForm(huh.NewGroup(
		huh.NewConfirm().Title("").Affirmative("Back to menu").Negative("Exit").Value(new(bool)),
	)).Run()
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

func (s *session) fetchWorkflows() ([]*sggosdk.GeneratedWorkflowsListAllMsg, error) {
	var response *sggosdk.WorkflowsListAll
	err := output.WithSpinner("Fetching workflows...", func() error {
		var e error
		response, e = s.client.Workflows.ListAllWorkflows(
			context.Background(), s.org, s.wfGrp, &sggosdk.ListAllWorkflowsRequest{},
		)
		return e
	})
	if err != nil {
		return nil, err
	}
	return response.Msg, nil
}

func (s *session) pickWorkflow(title string) (string, error) {
	workflows, err := s.fetchWorkflows()
	if err != nil {
		return "", err
	}
	if len(workflows) == 0 {
		return "", fmt.Errorf("no workflows found in %s / %s", s.org, s.wfGrp)
	}
	items := make([]tui.Item, len(workflows))
	for i, wf := range workflows {
		items[i] = tui.Item{
			ID:          wf.ResourceName,
			Label:       wf.ResourceName,
			Description: wf.Description,
			Badge:       wf.LatestWfrunStatus,
		}
	}
	return tui.NewPicker(title, contextLabel(s.org, s.wfGrp), items)
}

func (s *session) pickStack(title string) (string, error) {
	var response *sggosdk.GeneratedStackListAllResponse
	err := output.WithSpinner("Fetching stacks...", func() error {
		var e error
		response, e = s.client.Stacks.ListAllStacks(
			context.Background(), s.org, s.wfGrp, &sggosdk.ListAllStacksRequest{},
		)
		return e
	})
	if err != nil {
		return "", err
	}
	if len(response.Msg) == 0 {
		return "", fmt.Errorf("no stacks found in %s / %s", s.org, s.wfGrp)
	}
	items := make([]tui.Item, len(response.Msg))
	for i, st := range response.Msg {
		items[i] = tui.Item{
			ID:          st.ResourceName,
			Label:       st.ResourceName,
			Description: st.Description,
			Badge:       st.LatestWfStatus,
		}
	}
	return tui.NewPicker(title, contextLabel(s.org, s.wfGrp), items)
}

func contextLabel(org, wfGrp string) string {
	if org == "" {
		return ""
	}
	return org + " / " + wfGrp
}

