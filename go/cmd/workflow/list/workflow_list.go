package list

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/StackGuardian/sg-cli/cmd/output"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type RunOptions struct {
	OutputJson bool
	NoTUI      bool
}

func NewListCmd(c *client.Client) *cobra.Command {
	opts := &RunOptions{}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all workflows",
		Long:  "List all workflows in the specified workflow group. Opens an interactive browser by default.",
		Run: func(cmd *cobra.Command, args []string) {
			org := cmd.Parent().Flags().Lookup("org").Value.String()
			wfGrp := cmd.Parent().Flags().Lookup("workflow-group").Value.String()

			var response *sggosdk.WorkflowsListAll
			err := output.WithSpinner("Fetching workflows...", func() error {
				var apiErr error
				response, apiErr = c.Workflows.ListAllWorkflows(
					context.Background(),
					org,
					wfGrp,
					&sggosdk.ListAllWorkflowsRequest{},
				)
				return apiErr
			})
			if err != nil {
				output.Error(err.Error())
				os.Exit(1)
			}

			if opts.OutputJson {
				cmd.Println(response)
				return
			}

			if len(response.Msg) == 0 {
				output.Warning("No workflows found in this workflow group.")
				return
			}

			// Non-TTY or --no-tui: plain list output
			if opts.NoTUI {
				renderPlainList(response.Msg)
				return
			}

			// Check if stdout is a TTY
			if fi, _ := os.Stdout.Stat(); (fi.Mode()&os.ModeCharDevice) == 0 {
				renderPlainList(response.Msg)
				return
			}

			// Launch interactive TUI
			p := tea.NewProgram(
				newListModel(response.Msg, org, wfGrp),
				tea.WithAltScreen(),
			)
			if _, err := p.Run(); err != nil {
				output.Error("TUI error: " + err.Error())
				os.Exit(1)
			}
		},
	}

	listCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output API response as JSON.")
	listCmd.Flags().BoolVar(&opts.NoTUI, "no-tui", false, "Print a plain list instead of opening the interactive browser.")
	return listCmd
}

// ---------------------------------------------------------------------------
// Plain (non-TUI) output
// ---------------------------------------------------------------------------

func renderPlainList(workflows []*sggosdk.GeneratedWorkflowsListAllMsg) {
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED")).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E"))

	output.Separator()
	for _, wf := range workflows {
		fmt.Println(nameStyle.Render(wf.ResourceName))
		if wf.Description != "" {
			fmt.Println(descStyle.Render("  " + wf.Description))
		}
		if wf.LatestWfrunStatus != "" {
			fmt.Println(statusStyle.Render("  Status: " + wf.LatestWfrunStatus))
		}
		output.Separator()
	}
}

// ---------------------------------------------------------------------------
// Bubble Tea interactive list model
// ---------------------------------------------------------------------------

var (
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED")).Bold(true).Padding(0, 1)
	itemStyle     = lipgloss.NewStyle().Padding(0, 2)
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F9FAFB")).
			Background(lipgloss.Color("#7C3AED")).
			Bold(true).
			Padding(0, 2)
	descriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Padding(0, 4)
	helpBarStyle     = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6B7280")).
				BorderStyle(lipgloss.NormalBorder()).
				BorderTop(true).
				BorderForeground(lipgloss.Color("#374151")).
				Padding(0, 1)
	detailBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")).
			Padding(1, 2).
			Margin(1, 0)
	statusColors = map[string]lipgloss.Color{
		"COMPLETED": "#22C55E",
		"ERRORED":   "#EF4444",
		"RUNNING":   "#3B82F6",
		"PENDING":   "#F59E0B",
	}
)

func statusColor(s string) lipgloss.Color {
	if c, ok := statusColors[strings.ToUpper(s)]; ok {
		return c
	}
	return lipgloss.Color("#6B7280")
}

type listModel struct {
	workflows []*sggosdk.GeneratedWorkflowsListAllMsg
	cursor    int
	org       string
	wfGrp     string
	detail    bool
}

func newListModel(workflows []*sggosdk.GeneratedWorkflowsListAllMsg, org, wfGrp string) listModel {
	return listModel{
		workflows: workflows,
		org:       org,
		wfGrp:     wfGrp,
	}
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			if m.detail {
				m.detail = false
			} else {
				return m, tea.Quit
			}
		case "up", "k":
			if !m.detail && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if !m.detail && m.cursor < len(m.workflows)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.detail = !m.detail
		}
	}
	return m, nil
}

func (m listModel) View() string {
	if m.detail {
		return m.detailView()
	}
	return m.listView()
}

func (m listModel) listView() string {
	var b strings.Builder

	header := titleStyle.Render("Workflows") +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(
			fmt.Sprintf("  %s / %s  (%d)", m.org, m.wfGrp, len(m.workflows)),
		)
	b.WriteString(header + "\n\n")

	for i, wf := range m.workflows {
		statusBadge := ""
		if wf.LatestWfrunStatus != "" {
			statusBadge = " " + lipgloss.NewStyle().
				Foreground(statusColor(wf.LatestWfrunStatus)).
				Render("["+wf.LatestWfrunStatus+"]")
		}

		if i == m.cursor {
			b.WriteString(selectedStyle.Render("▶ "+wf.ResourceName) + statusBadge + "\n")
			if wf.Description != "" {
				b.WriteString(descriptionStyle.Render(wf.Description) + "\n")
			}
		} else {
			b.WriteString(itemStyle.Render("  "+wf.ResourceName) + statusBadge + "\n")
			if wf.Description != "" {
				b.WriteString(descriptionStyle.Render(wf.Description) + "\n")
			}
		}
	}

	help := helpBarStyle.Render("↑/k  ↓/j  navigate   enter  view details   q  quit")
	b.WriteString("\n" + help)

	return b.String()
}

func (m listModel) detailView() string {
	wf := m.workflows[m.cursor]

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Width(18)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#F9FAFB"))

	kv := func(label, value string) string {
		if value == "" {
			return ""
		}
		return labelStyle.Render(label+":") + " " + valueStyle.Render(value) + "\n"
	}

	var details strings.Builder
	details.WriteString(titleStyle.Render(wf.ResourceName) + "\n\n")
	details.WriteString(kv("Description", wf.Description))
	details.WriteString(kv("Type", wf.WfType))
	details.WriteString(kv("Status", wf.LatestWfrunStatus))
	details.WriteString(kv("Org", m.org))
	details.WriteString(kv("Workflow Group", m.wfGrp))
	details.WriteString(kv("Resource ID", wf.ResourceId))

	content := detailBoxStyle.Render(details.String())
	help := helpBarStyle.Render("esc / q  back to list")

	return content + "\n" + help
}
