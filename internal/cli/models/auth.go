package models

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	authapi "lesta-start-battleship/cli/internal/api/auth"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"
)

type AuthModel struct {
	login       string
	password    string
	email       string
	activeField int // 0 - Логин, 1 - Пароль
	activeTab   int // 0 - Авторизация, 1 - Регистрация
	authMethod  int // 0 - Логин/Пароль, 1 - Google, 2 - Яндекс
	errorMsg    string
	Clients     *clientdeps.Client
}

func NewAuthModel(clients *clientdeps.Client) *AuthModel {
	return &AuthModel{
		login:       "",
		password:    "",
		email:       "",
		errorMsg:    "",
		activeField: 0,
		activeTab:   0,
		authMethod:  0,
		Clients:     clients,
	}
}

func (m *AuthModel) Init() tea.Cmd {
	return nil
}

func (m *AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m.handleEnter()

		case tea.KeyTab:
			m.errorMsg = ""
			m.activeTab = (m.activeTab + 1) % 2
			m.activeField = 0
			return m, nil

		case tea.KeyLeft:
			m.errorMsg = ""
			m.authMethod = 0
			if m.activeTab == 0 {
				m.authMethod = (m.authMethod - 1 + 3) % 3
			}
			return m, nil

		case tea.KeyRight:
			m.errorMsg = ""
			m.authMethod = 0
			if m.activeTab == 0 {
				m.authMethod = (m.authMethod + 1) % 3
			}
			return m, nil

		case tea.KeyDown:
			m.errorMsg = ""
			if m.authMethod == 0 {
				maxField := 1
				if m.activeTab == 1 {
					maxField = 2
				}
				if m.activeField < maxField {
					m.activeField++
				}
			}
			return m, nil

		case tea.KeyUp:
			m.errorMsg = ""
			if m.authMethod == 0 {
				if m.activeField > 0 {
					m.activeField--
				}
			}
			return m, nil

		case tea.KeyBackspace:
			m.errorMsg = ""
			if m.authMethod == 0 {
				if m.activeField == 0 && len(m.login) > 0 {
					m.login = m.login[:len(m.login)-1]
				} else if m.activeTab == 1 && m.activeField == 1 && len(m.email) > 0 {
					m.email = m.email[:len(m.email)-1]
				} else if ((m.activeTab == 0 && m.activeField == 1) || (m.activeTab == 1 && m.activeField == 2)) && len(m.password) > 0 {
					m.password = m.password[:len(m.password)-1]
				}
			}
			return m, nil

		case tea.KeyRunes:
			m.errorMsg = ""
			if m.authMethod == 0 {
				if m.activeField == 0 {
					m.login += string(msg.Runes)
				} else if m.activeTab == 1 && m.activeField == 1 {
					m.email += string(msg.Runes)
				} else if (m.activeTab == 0 && m.activeField == 1) || (m.activeTab == 1 && m.activeField == 2) {
					m.password += string(msg.Runes)
				}
			}
			return m, nil

		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *AuthModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Морской Бой"))
	sb.WriteString("\n\n")

	authTab := "Авторизация"
	regTab := "Регистрация"
	if m.activeTab == 0 {
		authTab = ui.SelectedStyle.Render(authTab)
		regTab = ui.NormalStyle.Render(regTab)
	} else {
		authTab = ui.NormalStyle.Render(authTab)
		regTab = ui.SelectedStyle.Render(regTab)
	}
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, authTab, " | ", regTab))
	sb.WriteString("\n\n")

	if m.activeTab == 0 {
		methods := []string{"Логин/Пароль", "Google", "Яндекс"}
		for i, method := range methods {
			if i == m.authMethod {
				sb.WriteString(ui.SelectedStyle.Render("[" + method + "]"))
			} else {
				sb.WriteString(ui.NormalStyle.Render(method + " "))
			}
		}
	} else {
		sb.WriteString(ui.SelectedStyle.Render("Логин/Пароль"))
	}
	sb.WriteString("\n\n")

	if m.authMethod == 0 {
		loginLabel := "Логин: "
		if m.activeField == 0 {
			sb.WriteString(ui.SelectedStyle.Render(loginLabel))
		} else {
			sb.WriteString(ui.NormalStyle.Render(loginLabel))
		}
		sb.WriteString(m.login)
		if m.activeField == 0 {
			sb.WriteString("_")
		}
		sb.WriteString("\n")

		if m.activeTab == 1 {
			emailLabel := "Email: "
			if m.activeField == 1 {
				sb.WriteString(ui.SelectedStyle.Render(emailLabel))
			} else {
				sb.WriteString(ui.NormalStyle.Render(emailLabel))
			}
			sb.WriteString(m.email)
			if m.activeField == 1 {
				sb.WriteString("_")
			}
			sb.WriteString("\n")
		}

		passwordLabel := "Пароль: "
		// Для авторизации поле password = 1, для регистрации = 2
		passwordActive := (m.activeTab == 0 && m.activeField == 1) || (m.activeTab == 1 && m.activeField == 2)
		if passwordActive {
			sb.WriteString(ui.SelectedStyle.Render(passwordLabel))
		} else {
			sb.WriteString(ui.NormalStyle.Render(passwordLabel))
		}
		sb.WriteString(strings.Repeat("*", len(m.password)))
		if passwordActive {
			sb.WriteString("_")
		}

	} else {
		provider := "Google"
		if m.authMethod == 2 {
			provider = "Яндекс"
		}
		sb.WriteString(fmt.Sprintf("Для входа через %s Нажмите Enter", provider))
	}
	sb.WriteString("\n")

	if m.errorMsg != "" {
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg + "\n"))
	}

	sb.WriteString(ui.HelpStyle.Render("\nTab - Авторизация/Регистрация, ←/→ - выбор метода, Enter - подтвердить"))
	sb.WriteString(ui.HelpStyle.Render("\n↑/↓ - Переключение полей, Esc/Ctrl+C - выход\n"))

	return sb.String()
}

func (m *AuthModel) handleEnter() (tea.Model, tea.Cmd) {
	// Логин/пароль
	if m.authMethod == 0 {
		maxField := 1
		if m.activeTab == 1 {
			maxField = 2
		}
		if m.activeField < maxField {
			m.activeField++
			return m, nil
		}

		if m.activeTab == 0 && m.activeField == 1 { // Авторизация (2 поля: логин, пароль)
			if m.login == "" || m.password == "" {
				m.errorMsg = "Введите логин и пароль"
				return m, nil
			}
			// логика авторизации
			ctx := context.Background()
			_, profile, err := m.Clients.AuthClient.Login(ctx, authapi.LoginRequest{Username: m.login, Password: m.password})
			if err != nil {
				m.errorMsg = fmt.Sprintf("Ошибка авторизации: %v", err)
				return m, nil
			}
			return m, func() tea.Msg {
				return AuthSuccessMsg{
					Username: profile.Username,
					Gold:     profile.Currency.Gold,
					ID:       profile.ID,
					//Username: m.login,
					//Gold:     1000,
				}
			}
		} else if m.activeTab == 1 && m.activeField == 2 { // Регистрация (3 поля: логин, email, пароль)
			if len(m.login) < 3 || len(m.password) < 6 {
				m.errorMsg = "Логин минимум 3 символа, пароль - 6"
				return m, nil
			}
			ctx := context.Background()
			_, profile, err := m.Clients.AuthClient.Register(ctx, authapi.UserRegRequest{
				Username: m.login,
				Password: m.password,
				Email:    m.email,
			})
			if err != nil {
				m.errorMsg = fmt.Sprintf("%v", err)
				return m, nil
			}
			return m, func() tea.Msg {
				return AuthSuccessMsg{
					ID:       profile.ID,
					Username: profile.Username,
					Gold:     profile.Currency.Gold,
				}
			}
		}
	}

	// OAuth авторизация
	provider := []string{"", "google", "yandex"}[m.authMethod]
	ctx := context.Background()
	deviceAuth, err := m.Clients.AuthClient.InitOAuthDeviceFlow(ctx, provider)
	if err != nil {
		m.errorMsg = fmt.Sprintf("Ошибка авторизации через %s: %v", provider, err)
		return m, nil
	}

	return NewOAuthModel(m, provider, m.Clients, deviceAuth.VerificationURL, deviceAuth.DeviceCode,
		deviceAuth.UserCode, deviceAuth.Interval, deviceAuth.ExpiresIn), nil
	//return NewOAuthModel(m, provider, m.Clients, "deviceAuth.VerificationURL", "ABCD-1234"), nil
}
