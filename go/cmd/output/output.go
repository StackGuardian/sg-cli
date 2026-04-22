package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Success prints a green success message with a checkmark badge.
func Success(msg string) {
	fmt.Println(successBadge + " " + SuccessStyle.Render(msg))
}

// Error prints a red error message with an X badge to stderr.
func Error(msg string) {
	fmt.Fprintln(os.Stderr, errorBadge+" "+ErrorStyle.Render(msg))
}

// Warning prints a yellow warning message with a ! badge.
func Warning(msg string) {
	fmt.Println(warningBadge + " " + WarningStyle.Render(msg))
}

// Info prints a blue info message with an i badge.
func Info(msg string) {
	fmt.Println(infoBadge + " " + InfoStyle.Render(msg))
}

// URL prints a styled, underlined URL with a hint label.
func URL(label, url string) {
	fmt.Println(HintStyle.Render(label) + " " + URLStyle.Render(url))
}

// KV prints a key-value pair with aligned label styling.
func KV(label, value string) {
	fmt.Println(LabelStyle.Render(label+":") + " " + ValueStyle.Render(value))
}

// Section prints a styled section header.
func Section(title string) {
	fmt.Println(SectionStyle.Render(title))
	fmt.Println()
}

// Banner renders the root CLI banner.
func Banner(version string) string {
	title := PrimaryStyle.Render("sg-cli")
	ver := Muted.Render("v" + version)
	tagline := Muted.Render("Manage resources on the StackGuardian platform")

	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center, title, "  ", ver),
		tagline,
	)
	return BannerStyle.Render(content)
}

// ---------------------------------------------------------------------------
// Spinner — wraps a blocking function call with an animated spinner
// ---------------------------------------------------------------------------

type spinnerModel struct {
	spinner  spinner.Model
	label    string
	done     bool
	quitting bool
}

type doneMsg struct{}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case doneMsg:
		m.done = true
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m spinnerModel) View() string {
	if m.done || m.quitting {
		return ""
	}
	return lipgloss.NewStyle().Foreground(colorPrimary).Render(m.spinner.View()) +
		" " + Muted.Render(m.label)
}

// WithSpinner runs fn in a goroutine while showing a spinner with label.
// Returns any error fn produced. Safe to call in non-TTY (silently no-ops spinner).
func WithSpinner(label string, fn func() error) error {
	// Check if stdout is a TTY; skip spinner in CI/pipe environments
	if fi, _ := os.Stdout.Stat(); (fi.Mode() & os.ModeCharDevice) == 0 {
		return fn()
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(colorPrimary)

	m := spinnerModel{spinner: s, label: label}

	var fnErr error
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))

	go func() {
		fnErr = fn()
		p.Send(doneMsg{})
	}()

	if _, err := p.Run(); err != nil {
		return err
	}
	return fnErr
}

// ---------------------------------------------------------------------------
// Helpers used across commands
// ---------------------------------------------------------------------------

// Separator prints a dim horizontal rule.
func Separator() {
	width := 50
	fmt.Println(Muted.Render(strings.Repeat("─", width)))
}

// Newline prints a blank line.
func Newline() {
	fmt.Println()
}
