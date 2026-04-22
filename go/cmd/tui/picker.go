// Package tui provides shared interactive terminal UI components.
package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Item is a selectable entry in the picker.
type Item struct {
	ID          string // the value returned on selection (e.g. ResourceName)
	Label       string // primary display text
	Description string // secondary display text (optional)
	Badge       string // short status badge (optional, e.g. "COMPLETED")
}

const pageSize = 8 // items visible per page (each item = 2 lines: label + desc)

var (
	pickerTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#88B7DA")).
				Bold(true).
				Padding(0, 1)

	pickerItemStyle = lipgloss.NewStyle().Padding(0, 2)

	pickerSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F9FAFB")).
				Background(lipgloss.Color("#88B7DA")).
				Bold(true).
				Padding(0, 2)

	pickerDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Padding(0, 4)

	pickerHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(lipgloss.Color("#374151")).
			Padding(0, 1)

	pickerPageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true)

	badgeColors = map[string]lipgloss.Color{
		"COMPLETED": "#22C55E",
		"ERRORED":   "#EF4444",
		"RUNNING":   "#3B82F6",
		"PENDING":   "#F59E0B",
	}
)

func badgeColor(s string) lipgloss.Color {
	if c, ok := badgeColors[strings.ToUpper(s)]; ok {
		return c
	}
	return lipgloss.Color("#6B7280")
}

// pickerModel is the Bubble Tea model for the picker.
type pickerModel struct {
	title    string
	subtitle string
	items    []Item
	cursor   int   // absolute index into items
	offset   int   // index of first visible item
	height   int   // terminal height (updated on WindowSizeMsg)
	selected string
	quitting bool
}

func (m pickerModel) pageSize() int {
	// Reserve lines: 3 header + 1 gap + 3 help + some margin
	reserved := 8
	h := m.height
	if h == 0 {
		h = 24 // safe fallback before first WindowSizeMsg
	}
	available := h - reserved
	if available < 4 {
		available = 4
	}
	// Each item takes up to 2 lines (label + description)
	items := available / 2
	if items < 2 {
		items = 2
	}
	return items
}

// NewPicker runs an interactive full-screen picker and returns the selected Item.ID,
// or an error if the user cancelled (ctrl+c / esc) or the list is empty.
func NewPicker(title, subtitle string, items []Item) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to select from")
	}

	// Clear the normal-screen buffer before entering alt-screen so that
	// when the picker exits and the terminal restores the prior screen,
	// stale output from previous actions is not visible.
	fmt.Print("\033[H\033[2J")

	m := pickerModel{title: title, subtitle: subtitle, items: items}
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return "", err
	}

	final := result.(pickerModel)
	if final.quitting || final.selected == "" {
		return "", fmt.Errorf("cancelled")
	}
	return final.selected, nil
}

func (m pickerModel) Init() tea.Cmd { return nil }

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		// Re-clamp offset after resize
		m.offset = clampOffset(m.cursor, m.offset, m.pageSize())

	case tea.KeyMsg:
		ps := m.pageSize()
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.offset {
					m.offset = m.cursor
				}
			}

		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
				if m.cursor >= m.offset+ps {
					m.offset = m.cursor - ps + 1
				}
			}

		case "pgup", "left", "h":
			m.cursor -= ps
			if m.cursor < 0 {
				m.cursor = 0
			}
			m.offset = clampOffset(m.cursor, m.offset, ps)

		case "pgdn", "right", "l":
			m.cursor += ps
			if m.cursor >= len(m.items) {
				m.cursor = len(m.items) - 1
			}
			m.offset = clampOffset(m.cursor, m.offset, ps)

		case "home", "g":
			m.cursor = 0
			m.offset = 0

		case "end", "G":
			m.cursor = len(m.items) - 1
			m.offset = clampOffset(m.cursor, m.offset, ps)

		case "enter", " ":
			m.selected = m.items[m.cursor].ID
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m pickerModel) View() string {
	var b strings.Builder
	ps := m.pageSize()

	// Brand line
	brand := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#88B7DA")).
		Bold(true).
		Render("StackGuardian")
	brandLine := brand + lipgloss.NewStyle().Foreground(lipgloss.Color("#374151")).Render("  ·  sg-cli")
	b.WriteString(brandLine + "\n\n")

	// Header
	header := pickerTitleStyle.Render(m.title)
	if m.subtitle != "" {
		header += lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render("  " + m.subtitle)
	}

	// Page indicator (only shown when list overflows)
	pageInfo := ""
	if len(m.items) > ps {
		currentPage := m.offset/ps + 1
		totalPages := (len(m.items) + ps - 1) / ps
		pageInfo = pickerPageStyle.Render(fmt.Sprintf("  page %d/%d", currentPage, totalPages))
	}

	b.WriteString(header + pageInfo + "\n\n")

	// Visible window
	end := m.offset + ps
	if end > len(m.items) {
		end = len(m.items)
	}

	// Determine if any visible item has a badge — use table layout only then
	hasBadges := false
	labelWidth := 4
	badgeWidth := 0
	for i := m.offset; i < end; i++ {
		if m.items[i].Badge != "" {
			hasBadges = true
		}
		if l := lipgloss.Width(m.items[i].Label); l > labelWidth {
			labelWidth = l
		}
		if bw := lipgloss.Width(m.items[i].Badge); bw > badgeWidth {
			badgeWidth = bw
		}
	}
	if badgeWidth > 0 {
		badgeWidth += 2 // brackets
	}

	for i := m.offset; i < end; i++ {
		item := m.items[i]

		if hasBadges {
			// Table layout: name | badge | description all on one line
			badgeStr := strings.Repeat(" ", badgeWidth)
			if item.Badge != "" {
				badgeStr = lipgloss.NewStyle().
					Foreground(badgeColor(item.Badge)).
					Width(badgeWidth).
					Render("[" + item.Badge + "]")
			}

			nameWidth := labelWidth + 4 // padding for cursor prefix + spacing
			if i == m.cursor {
				b.WriteString(pickerSelectedStyle.Width(nameWidth).Render("▶ " + item.Label))
			} else {
				b.WriteString(pickerItemStyle.Width(nameWidth).Render("  " + item.Label))
			}
			b.WriteString("  " + badgeStr)
			if item.Description != "" {
				b.WriteString("  " + pickerDescStyle.Render(item.Description))
			}
			b.WriteString("\n")
		} else {
			// Simple layout: name on one line, description indented below
			if i == m.cursor {
				b.WriteString(pickerSelectedStyle.Render("▶ " + item.Label) + "\n")
			} else {
				b.WriteString(pickerItemStyle.Render("  " + item.Label) + "\n")
			}
			if item.Description != "" {
				b.WriteString(pickerDescStyle.Render(item.Description) + "\n")
			}
		}
	}

	// Scroll hints
	scrollHints := ""
	if m.offset > 0 {
		scrollHints += pickerPageStyle.Render("  ↑ more above")
	}
	if end < len(m.items) {
		if scrollHints != "" {
			scrollHints += "   "
		}
		scrollHints += pickerPageStyle.Render("  ↓ more below")
	}
	if scrollHints != "" {
		b.WriteString("\n" + scrollHints)
	}

	// Help bar
	sep := lipgloss.NewStyle().Foreground(lipgloss.Color("#374151")).Render("  │  ")
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#88B7DA")).Bold(true)
	actStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))

	key := func(k, act string) string {
		return keyStyle.Render(k) + " " + actStyle.Render(act)
	}

	parts := []string{
		key("↑/k  ↓/j", "navigate"),
	}
	if len(m.items) > ps {
		parts = append(parts, key("PgUp/PgDn", "page"), key("g/G", "top/bottom"))
	}
	parts = append(parts, key("enter", "select"), key("esc", "cancel"))

	helpText := strings.Join(parts, sep)
	b.WriteString("\n" + pickerHelpStyle.Render(helpText))

	return b.String()
}

// clampOffset keeps cursor visible within the current page window.
func clampOffset(cursor, offset, ps int) int {
	if cursor < offset {
		return cursor
	}
	if cursor >= offset+ps {
		return cursor - ps + 1
	}
	return offset
}
