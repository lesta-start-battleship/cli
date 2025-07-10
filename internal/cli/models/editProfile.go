package models

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"lesta-start-battleship/cli/internal/api/auth"
	"lesta-start-battleship/cli/internal/cli/ui"
	"lesta-start-battleship/cli/internal/clientdeps"
	"strings"
)

type EditProfileModel struct {
	id            int
	username      string
	tempNick      string
	tempPass      string
	activeTab     int // 0 - ник, 1 - пароль, 2 - удалить аккаунт
	activeField   int // 0 - ввод, 1 - подтверждение
	errorMsg      string
	successMsg    string
	gold          int
	Clients       *clientdeps.Client
	confirmDelete bool // true если показываем подтверждение удаления
}

func NewEditProfileModel(id int, username string, gold int, clients *clientdeps.Client) *EditProfileModel {
	return &EditProfileModel{
		id:          id,
		username:    username,
		gold:        gold,
		tempNick:    username,
		activeTab:   0,
		activeField: 0,
		Clients:     clients,
	}
}

func (m *EditProfileModel) Init() tea.Cmd {
	return nil
}

func (m *EditProfileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.confirmDelete {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEsc:
				m.confirmDelete = false
				m.activeField = 0
				return m, nil
			case tea.KeyEnter:
				ctx := context.Background()
				err := m.Clients.AuthClient.DeleteUser(ctx, m.id)
				if err != nil {
					m.errorMsg = err.Error()
					m.confirmDelete = false
					m.activeField = 0
					return m, nil
				}
				return NewAuthModel(m.Clients), nil
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m.activeTab = (m.activeTab + 1) % 3
			m.activeField = 0
			return m, nil

		/*case tea.KeyDown:
			if m.activeTab != 2 && m.activeField < 1 {
				m.activeField++
			}
			return m, nil

		case tea.KeyUp:
			if m.activeTab != 2 && m.activeField > 0 {
				m.activeField--
			}
			return m, nil*/

		case tea.KeyEnter:
			if m.activeTab == 2 {
				m.confirmDelete = true
				return m, nil
			}
			if m.activeTab == 0 && m.activeField == 1 {
				if len(m.tempNick) < 3 {
					m.errorMsg = "Ник должен быть не менее 3 символов"
				} else {
					ctx := context.Background()
					profile, err := m.Clients.AuthClient.UpdateUser(ctx, m.id, auth.UpdateUserRequest{Username: m.tempNick})
					if err != nil {
						m.errorMsg = err.Error()
						return m, nil
					}
					m.username = profile.Username
					m.gold = profile.Currency.Gold
					m.errorMsg = ""
					m.tempNick = ""
					m.successMsg = "Ник успешно изменен!"
					return m, func() tea.Msg {
						return UsernameChangeMsg{NewUsername: m.username, Gold: m.gold}
					}
				}
			} else if m.activeTab == 1 && m.activeField == 1 {
				if len(m.tempPass) < 6 {
					m.errorMsg = "Пароль должен быть не менее 6 символов"
				} else {
					ctx := context.Background()
					profile, err := m.Clients.AuthClient.UpdateUser(ctx, m.id, auth.UpdateUserRequest{Password: m.tempPass})
					if err != nil {
						m.errorMsg = err.Error()
						return m, nil
					}
					m.username = profile.Username
					m.gold = profile.Currency.Gold
					m.errorMsg = ""
					m.successMsg = "Пароль успешно изменен!"
					m.tempPass = ""
					return m, func() tea.Msg {
						return UsernameChangeMsg{NewUsername: m.username, Gold: m.gold}
					}
				}
			} else {
				if m.activeTab != 2 && m.activeField < 1 {
					m.activeField = 1
				}
			}
			return m, nil

		case tea.KeyBackspace:
			if m.activeTab == 0 && m.activeField == 0 && len(m.tempNick) > 0 {
				m.tempNick = m.tempNick[:len(m.tempNick)-1]
			} else if m.activeTab == 1 && m.activeField == 0 && len(m.tempPass) > 0 {
				m.tempPass = m.tempPass[:len(m.tempPass)-1]
			}
			return m, nil

		case tea.KeyRunes:
			if m.activeTab == 0 && m.activeField == 0 {
				m.tempNick += string(msg.Runes)
			} else if m.activeTab == 1 && m.activeField == 0 {
				m.tempPass += string(msg.Runes)
			}
			return m, nil

		case tea.KeyEsc:
			return NewMainMenuModel(m.id, m.username, m.gold, m.Clients), nil
		}
	}
	return m, nil
}

func (m *EditProfileModel) View() string {
	var sb strings.Builder

	sb.WriteString(ui.TitleStyle.Render("Редактирование профиля"))
	sb.WriteString("\n\n")

	// Вкладки
	nickTab := "Ник"
	passTab := "Пароль"
	deleteTab := "Удалить аккаунт"
	if m.activeTab == 0 {
		nickTab = ui.SelectedStyle.Render(nickTab)
	} else {
		nickTab = ui.NormalStyle.Render(nickTab)
	}
	if m.activeTab == 1 {
		passTab = ui.SelectedStyle.Render(passTab)
	} else {
		passTab = ui.NormalStyle.Render(passTab)
	}
	if m.activeTab == 2 {
		deleteTab = ui.ErrorStyle.Render(deleteTab)
	} else {
		deleteTab = ui.NormalStyle.Render(deleteTab)
	}
	sb.WriteString(nickTab + " | " + passTab + " | " + deleteTab)
	sb.WriteString("\n\n")

	if m.confirmDelete {
		sb.WriteString(ui.ErrorStyle.Render("Вы точно хотите удалить аккаунт?"))
		sb.WriteString("\n")
		sb.WriteString(ui.NormalStyle.Render("Esc - отмена, Enter - удалить"))
		return sb.String()
	}

	if m.activeTab == 0 {
		sb.WriteString("Текущий ник: " + m.username + "\n\n")
		sb.WriteString("Новый ник:\n")
		if m.activeField == 0 {
			sb.WriteString(ui.SelectedStyle.Render("> " + m.tempNick + "_"))
		} else {
			sb.WriteString(" " + m.tempNick)
			sb.WriteString("\n\n")
			if m.activeField == 1 {
				sb.WriteString(ui.SuccessStyle.Render("Нажмите Enter для сохранения"))
			}
		}
	} else if m.activeTab == 1 {
		sb.WriteString("Новый пароль:\n")
		if m.activeField == 0 {
			sb.WriteString(ui.SelectedStyle.Render("> " + strings.Repeat("*", len(m.tempPass)) + "_"))
		} else {
			sb.WriteString(" " + strings.Repeat("*", len(m.tempPass)))
			sb.WriteString("\n\n")
			if m.activeField == 1 {
				sb.WriteString(ui.SuccessStyle.Render("Нажмите Enter для сохранения"))
			}
		}
	} else if m.activeTab == 2 {
		sb.WriteString("\n")
		sb.WriteString(ui.ErrorStyle.Render("Нажмите Enter для удаления аккаунта"))
	}

	if m.errorMsg != "" {
		sb.WriteString("\n\n")
		sb.WriteString(ui.ErrorStyle.Render(m.errorMsg))
	}
	if m.successMsg != "" {
		sb.WriteString("\n\n")
		sb.WriteString(ui.SuccessStyle.Render(m.successMsg))
	}

	sb.WriteString("\n\n")
	sb.WriteString(ui.HelpStyle.Render("Tab - переключение вкладок, Enter - подтвердить, Esc - назад"))

	return sb.String()
}
