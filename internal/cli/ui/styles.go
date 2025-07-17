package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#1E88E5")).Bold(true).Padding(0, 1)

	ErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5252")).Bold(true)
	SuccessStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#69F0AE")).Bold(true)
	NormalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	SelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD600")).Bold(true)
	SubtitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500))")).Bold(true)
	WarningStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF9800")).Bold(true)

	AlertStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).
			Background(lipgloss.Color("#330000")).Bold(true).Padding(0, 1)

	ChatContainerStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#555555")).Padding(0, 1)

	ChatHeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#1E88E5")).Bold(true).Padding(0, 1)

	ChatInputStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	SystemMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Italic(true)
	OwnMessageStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	OtherMessageStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC"))
	NewMessageStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true)

	HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Italic(true)

	// Стили для магазина
	PromotionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true)

	SelectedTabStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#5A56E0")).Padding(0, 2)

	NormalTabStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A1A1A1")).Padding(0, 2)

	// Стили для элементов
	ItemNameStyle = lipgloss.NewStyle().Bold(true)

	ItemPriceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))

	ItemDescStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))

	ItemDetailStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
)
