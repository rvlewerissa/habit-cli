package habits

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/vittolewerissa/hbt/internal/shared/model"
	"github.com/vittolewerissa/hbt/internal/shared/ui"
)

// getFilteredEmojis returns emojis filtered by search term
func (m *FormModel) getFilteredEmojis() []string {
	search := strings.ToLower(m.emojiSearch.Value())
	if search == "" {
		return model.CommonEmojis
	}

	var filtered []string
	for _, emoji := range model.CommonEmojis {
		keywords := model.EmojiKeywords[emoji]
		if strings.Contains(strings.ToLower(keywords), search) {
			filtered = append(filtered, emoji)
		}
	}
	return filtered
}

// ensureVisibleEmoji ensures the selected emoji is within the visible scroll area
func (m *FormModel) ensureVisibleEmoji() {
	const emojisPerRow = 8
	const maxVisibleRows = 8

	// If (none) is selected, scroll to top
	if m.emojiIndex == -1 {
		m.scrollOffset = 0
		return
	}

	filtered := m.getFilteredEmojis()
	if len(filtered) == 0 {
		return
	}

	// Calculate current row of selected emoji
	selectedRow := m.emojiIndex / emojisPerRow

	// Adjust scroll offset to keep selected emoji visible
	if selectedRow < m.scrollOffset {
		m.scrollOffset = selectedRow
	} else if selectedRow >= m.scrollOffset+maxVisibleRows {
		m.scrollOffset = selectedRow - maxVisibleRows + 1
	}

	// Ensure scroll offset is valid
	totalRows := (len(filtered) + emojisPerRow - 1) / emojisPerRow
	maxScroll := totalRows - maxVisibleRows
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scrollOffset > maxScroll {
		m.scrollOffset = maxScroll
	}
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}

// renderEmojiSelector renders the emoji field display
func (m *FormModel) renderEmojiSelector(focused bool) string {
	var display string
	if m.selectedEmoji == "" {
		display = "(none)"
		if !focused {
			return ui.MutedText.Render(display)
		}
	} else {
		display = m.selectedEmoji
		if !focused {
			return display
		}
	}

	// When focused, show with brackets
	return ui.SelectedItem.Render("[" + display + "]")
}

// renderEmojiModalBox renders the emoji picker modal content
func (m *FormModel) renderEmojiModalBox() string {
	var s string

	// Modal title
	s += ui.Title.Render("Pick an Emoji") + "\n\n"

	// Search input
	s += "Search: " + m.emojiSearch.View() + "\n\n"

	// Show (none) option
	noneText := "(none)"
	if m.emojiIndex == -1 {
		s += "[" + noneText + "]" + "\n"
	} else {
		s += " " + noneText + " " + "\n"
	}
	s += "\n"

	// Show filtered emoji grid - 8 per row, max 8 rows visible
	const emojisPerRow = 8
	const maxVisibleRows = 8
	filtered := m.getFilteredEmojis()

	if len(filtered) == 0 {
		s += ui.MutedText.Render("No emojis found") + "\n"
	} else {
		totalRows := (len(filtered) + emojisPerRow - 1) / emojisPerRow
		startRow := m.scrollOffset
		endRow := startRow + maxVisibleRows
		if endRow > totalRows {
			endRow = totalRows
		}

		// Show scroll indicator at top if not at beginning
		if startRow > 0 {
			s += ui.MutedText.Render("        ▲ more above ▲") + "\n"
		}

		// Render visible rows only
		for row := startRow; row < endRow; row++ {
			var rowEmojis []string
			for col := 0; col < emojisPerRow; col++ {
				idx := row*emojisPerRow + col
				if idx >= len(filtered) {
					break
				}
				emoji := filtered[idx]
				display := emoji
				if idx == m.emojiIndex {
					display = "[" + emoji + "]"
				} else {
					display = " " + emoji + " "
				}
				rowEmojis = append(rowEmojis, display)
			}
			s += lipgloss.JoinHorizontal(lipgloss.Left, rowEmojis...) + "\n"
		}

		// Show scroll indicator at bottom if there's more content
		if endRow < totalRows {
			s += ui.MutedText.Render("        ▼ more below ▼") + "\n"
		}
	}
	s += "\n"

	s += ui.MutedText.Render("↑↓←→: navigate  pgup/pgdn: scroll  enter: select  esc: cancel")

	// Add modal box styling with dark background
	modalWidth := 54
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.Primary).
		Background(lipgloss.Color("#1F2937")).
		Padding(1, 2).
		Width(modalWidth)

	return modalStyle.Render(s)
}
