package initCli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"lesta-start-battleship/cli/internal/cli/models"
	"lesta-start-battleship/cli/internal/clientdeps"
	guildStorage "lesta-start-battleship/cli/storage/guild"
)

type CLI struct {
	currentScreen tea.Model
	chatComponent *models.ChatComponent
	clients       *clientdeps.Client
	gold          int
	userID        int
	username      string
}

func NewCLI(clients *clientdeps.Client) *CLI {
	return &CLI{
		currentScreen: models.NewAuthModel(clients),
		chatComponent: models.NewChatComponent("", 0),
		clients:       clients,
	}
}

func (a *CLI) Init() tea.Cmd {
	return nil
}

func (a *CLI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEsc && a.chatComponent.IsVisible() && a.chatComponent.Focused {
			a.chatComponent.Close()
			return a, nil
		}
	}

	switch msg := msg.(type) {
	case models.ChatKeyHandledMsg:
		return a, nil

	case models.AuthSuccessMsg:
		a.userID = msg.ID
		a.gold = msg.Gold
		a.username = msg.Username
		a.currentScreen = models.NewMainMenuModel(a.userID, a.username, a.gold, a.clients)
		a.chatComponent = models.NewChatComponent(a.username, 1)
		return a, nil

	case models.LogoutMsg:
		a.userID = 0
		a.gold = 0
		a.username = ""
		guildStorage.CleanStorage()
		a.currentScreen = models.NewAuthModel(a.clients)
		a.chatComponent.Close()
		a.chatComponent = models.NewChatComponent("", 0)
		return a, nil

	case models.UsernameChangeMsg:
		a.username = msg.NewUsername
		a.gold = msg.Gold
		a.chatComponent.Username = msg.NewUsername
		return a, nil

	case models.OpenChatMsg:
		// Если GuildID передан, используем его для инициализации чата гильдии
		guildID := 0
		if msg.GuildID != 0 {
			guildID = msg.GuildID
		}
		a.chatComponent = models.NewChatComponent(a.username, guildID)
		a.chatComponent.Toggle()
		if a.chatComponent.IsVisible() {
			return a, a.chatComponent.Init()
		}
		return a, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}

		if msg.Type == tea.KeyCtrlG && a.chatComponent.IsVisible() {
			a.chatComponent.Focused = !a.chatComponent.Focused
			return a, nil
		}
	}

	var chatCmd tea.Cmd
	if a.chatComponent.IsVisible() {
		var updatedChat tea.Model
		updatedChat, chatCmd = a.chatComponent.Update(msg)
		if updatedChat != nil {
			a.chatComponent = updatedChat.(*models.ChatComponent)
		}

		if a.chatComponent.Focused {
			return a, chatCmd
		}
	}

	var mainCmd tea.Cmd
	a.currentScreen, mainCmd = a.currentScreen.Update(msg)

	return a, tea.Batch(mainCmd, chatCmd)
}

func (a *CLI) View() string {
	mainView := a.currentScreen.View()

	if a.chatComponent.IsVisible() {
		chatView := a.chatComponent.View()
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width(100).Render(mainView),
			chatView,
		)
	}

	return mainView
}
