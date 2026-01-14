package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Colors
var (
	Primary     = lipgloss.Color("#06B6D4") // Cyan
	Secondary   = lipgloss.Color("#8B5CF6") // Purple
	Success     = lipgloss.Color("#10B981") // Green
	Warning     = lipgloss.Color("#F59E0B") // Amber
	Danger      = lipgloss.Color("#EF4444") // Red
	Muted       = lipgloss.Color("#6B7280") // Gray
	Background  = lipgloss.Color("#1F2937") // Dark gray
	Foreground  = lipgloss.Color("#F9FAFB") // Light gray
	Border      = lipgloss.Color("#374151") // Medium gray
)

// Tab styles
var (
	ActiveTab = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(Primary).
			Padding(0, 2)

	InactiveTab = lipgloss.NewStyle().
			Foreground(Muted).
			Padding(0, 2)

	TabBar = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(Border).
		MarginBottom(1)
)

// List styles
var (
	SelectedItem = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	NormalItem = lipgloss.NewStyle().
			Foreground(Foreground)

	CompletedItem = lipgloss.NewStyle().
			Foreground(Success).
			Strikethrough(true)

	MutedText = lipgloss.NewStyle().
			Foreground(Muted)
)

// Form styles
var (
	FormLabel = lipgloss.NewStyle().
			Foreground(Muted).
			MarginRight(1)

	FormInput = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(0, 1)

	FormInputFocused = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Primary).
				Padding(0, 1)
)

// Status indicators
var (
	Checkbox = lipgloss.NewStyle().
			Foreground(Muted)

	CheckboxChecked = lipgloss.NewStyle().
				Foreground(Success)

	StreakBadge = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true)
)

// Layout
var (
	Container = lipgloss.NewStyle().
			PaddingTop(2).
			PaddingBottom(1).
			PaddingLeft(2).
			PaddingRight(2)

	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Foreground).
		MarginBottom(1)

	Subtitle = lipgloss.NewStyle().
			Foreground(Muted).
			MarginBottom(1)

	HelpText = lipgloss.NewStyle().
			Foreground(Muted).
			MarginTop(1)
)

// CategoryTag returns a styled category tag with the given color
func CategoryTag(name, color string) string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(color)).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 1).
		Render(name)
}

// TitledPanel renders a panel with the title inline with the top border
// Format: ╭─ Title ─────────────╮
func TitledPanel(title, content string, width, height int) string {
	borderColor := lipgloss.NewStyle().Foreground(Border)
	titleStyle := lipgloss.NewStyle().Foreground(Primary).Bold(true)

	// Build top border with title
	titleText := " " + title + " "
	titleLen := len(titleText)
	remainingWidth := width - 2 - titleLen - 1 // -2 for corners, -1 for initial dash
	if remainingWidth < 0 {
		remainingWidth = 0
	}

	topBorder := borderColor.Render("╭─") +
		titleStyle.Render(titleText) +
		borderColor.Render(repeat("─", remainingWidth)+"╮")

	// Content area with side borders
	contentStyle := lipgloss.NewStyle().
		Width(width - 6). // -6 for borders and padding
		Height(height - 3). // -3 for top/bottom borders and top padding
		PaddingLeft(2).
		PaddingRight(2).
		PaddingTop(1)

	contentLines := strings.Split(contentStyle.Render(content), "\n")
	var middle string
	for _, line := range contentLines {
		// Pad line to fill width
		lineWidth := lipgloss.Width(line)
		padding := width - 2 - lineWidth // -2 for side borders
		if padding < 0 {
			padding = 0
		}
		middle += borderColor.Render("│") + line + strings.Repeat(" ", padding) + borderColor.Render("│") + "\n"
	}

	// Bottom border
	bottomBorder := borderColor.Render("╰" + repeat("─", width-2) + "╯")

	return topBorder + "\n" + middle + bottomBorder
}

func repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
