package output

import "github.com/charmbracelet/lipgloss"

var (
	// Brand colors
	colorPrimary = lipgloss.Color("#88B7DA") // StackGuardian steel blue
	colorSuccess = lipgloss.Color("#22C55E")
	colorError   = lipgloss.Color("#EF4444")
	colorWarning = lipgloss.Color("#F59E0B")
	colorInfo    = lipgloss.Color("#88B7DA")
	colorMuted   = lipgloss.Color("#6B7280")
	colorURL     = lipgloss.Color("#88B7DA")
	colorWhite   = lipgloss.Color("#F9FAFB")

	// Base text styles
	Bold   = lipgloss.NewStyle().Bold(true)
	Muted  = lipgloss.NewStyle().Foreground(colorMuted)
	Italic = lipgloss.NewStyle().Italic(true)

	// Semantic styles
	SuccessStyle = lipgloss.NewStyle().Foreground(colorSuccess).Bold(true)
	ErrorStyle   = lipgloss.NewStyle().Foreground(colorError).Bold(true)
	WarningStyle = lipgloss.NewStyle().Foreground(colorWarning).Bold(true)
	InfoStyle    = lipgloss.NewStyle().Foreground(colorInfo)
	URLStyle     = lipgloss.NewStyle().Foreground(colorURL).Underline(true)
	PrimaryStyle = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)

	// Prefix badges
	successBadge = lipgloss.NewStyle().
			Foreground(colorWhite).
			Background(colorSuccess).
			Bold(true).
			Padding(0, 1).
			Render(" ✓ ")

	errorBadge = lipgloss.NewStyle().
			Foreground(colorWhite).
			Background(colorError).
			Bold(true).
			Padding(0, 1).
			Render(" ✗ ")

	warningBadge = lipgloss.NewStyle().
			Foreground(colorWhite).
			Background(colorWarning).
			Bold(true).
			Padding(0, 1).
			Render(" ! ")

	infoBadge = lipgloss.NewStyle().
			Foreground(colorWhite).
			Background(colorInfo).
			Bold(true).
			Padding(0, 1).
			Render(" i ")

	// Banner box
	BannerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 2)

	// Section header
	SectionStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(colorPrimary)

	// Key-value label
	LabelStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Width(18)

	ValueStyle = lipgloss.NewStyle().
			Foreground(colorWhite)

	// Dashboard URL hint
	HintStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true)
)
