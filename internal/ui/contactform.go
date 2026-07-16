package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"mayur-athavale-tui/internal/ui/theme"
)

const (
	telegramSendMessageURL = "https://api.telegram.org/bot%s/sendMessage"

	fieldName = iota
	fieldEmail
	fieldCompany
	fieldMessage
	fieldSubmit
	fieldCount
)

type submitStatus int

const (
	statusIdle submitStatus = iota
	statusSending
	statusSuccess
	statusError
)

// submitResultMsg carries a contact-form submission's outcome back into the
// Bubble Tea event loop. Delivered regardless of which tab is active when it
// arrives, so switching away mid-submit doesn't lose the result.
type submitResultMsg struct {
	ok      bool
	message string
}

// contactForm is the Contact tab's own small Bubble Tea sub-model. It only
// receives key messages while editing is true (entered with Enter, left with
// Esc) -- otherwise the global tab-navigation keys (digits, tab, h/j/k/l)
// behave exactly as they do on every other tab, so viewing this tab without
// intending to type never accidentally starts eating keystrokes.
type contactForm struct {
	botToken string
	chatID   string

	editing bool
	focus   int

	name    textinput.Model
	email   textinput.Model
	company textinput.Model
	message textarea.Model

	status    submitStatus
	statusMsg string
}

func newContactForm(botToken, chatID string) contactForm {
	name := textinput.New()
	name.Placeholder = "Jane Recruiter"
	name.CharLimit = 100

	email := textinput.New()
	email.Placeholder = "jane@company.com"
	email.CharLimit = 200

	company := textinput.New()
	company.Placeholder = "(optional)"
	company.CharLimit = 100

	message := textarea.New()
	message.Placeholder = "What's this about?"
	message.CharLimit = 2000
	message.SetHeight(5)
	message.ShowLineNumbers = false

	return contactForm{
		botToken: botToken,
		chatID:   chatID,
		name:     name,
		email:    email,
		company:  company,
		message:  message,
	}
}

func (f *contactForm) startEditing() {
	f.editing = true
	f.focus = fieldName
	f.status = statusIdle
	f.statusMsg = ""
	f.name.Focus()
}

func (f *contactForm) stopEditing() {
	f.editing = false
	f.name.Blur()
	f.email.Blur()
	f.company.Blur()
	f.message.Blur()
}

func (f *contactForm) setFocus(i int) {
	f.name.Blur()
	f.email.Blur()
	f.company.Blur()
	f.message.Blur()

	f.focus = (i + fieldCount) % fieldCount
	switch f.focus {
	case fieldName:
		f.name.Focus()
	case fieldEmail:
		f.email.Focus()
	case fieldCompany:
		f.company.Focus()
	case fieldMessage:
		f.message.Focus()
	}
}

func (f contactForm) update(msg tea.Msg) (contactForm, tea.Cmd) {
	switch msg := msg.(type) {
	case submitResultMsg:
		if msg.ok {
			f.status = statusSuccess
		} else {
			f.status = statusError
		}
		f.statusMsg = msg.message
		return f, nil

	case tea.KeyMsg:
		if !f.editing {
			if msg.String() == "enter" {
				f.startEditing()
			}
			return f, nil
		}

		switch msg.String() {
		case "esc":
			f.stopEditing()
			return f, nil

		case "tab", "down":
			if f.focus == fieldMessage && msg.String() == "down" {
				break // let the textarea handle cursor movement inside itself
			}
			f.setFocus(f.focus + 1)
			return f, nil

		case "shift+tab", "up":
			if f.focus == fieldMessage && msg.String() == "up" {
				break
			}
			f.setFocus(f.focus - 1)
			return f, nil

		case "enter":
			if f.focus == fieldSubmit {
				f.status = statusSending
				f.statusMsg = ""
				return f, f.submit()
			}
			if f.focus != fieldMessage {
				f.setFocus(f.focus + 1)
				return f, nil
			}
			// fieldMessage: fall through so the textarea inserts a newline
		}
	}

	var cmd tea.Cmd
	switch f.focus {
	case fieldName:
		f.name, cmd = f.name.Update(msg)
	case fieldEmail:
		f.email, cmd = f.email.Update(msg)
	case fieldCompany:
		f.company, cmd = f.company.Update(msg)
	case fieldMessage:
		f.message, cmd = f.message.Update(msg)
	}
	return f, cmd
}

func (f contactForm) submit() tea.Cmd {
	if strings.TrimSpace(f.name.Value()) == "" || strings.TrimSpace(f.email.Value()) == "" || strings.TrimSpace(f.message.Value()) == "" {
		return func() tea.Msg {
			return submitResultMsg{ok: false, message: "Name, email, and message are required."}
		}
	}
	if f.botToken == "" || f.chatID == "" {
		return func() tea.Msg {
			return submitResultMsg{ok: false, message: "Form isn't configured yet (missing TG_BOT_TOKEN/TG_CONTACT_CHAT_ID)."}
		}
	}

	company := strings.TrimSpace(f.company.Value())
	companyLine := ""
	if company != "" {
		companyLine = fmt.Sprintf("Company: %s\n", company)
	}
	text := fmt.Sprintf(
		"New message from the SSH portfolio\n\nName: %s\nEmail: %s\n%s\n%s",
		f.name.Value(), f.email.Value(), companyLine, f.message.Value(),
	)

	payload := map[string]string{
		"chat_id": f.chatID,
		"text":    text,
	}
	url := fmt.Sprintf(telegramSendMessageURL, f.botToken)

	return func() tea.Msg {
		body, err := json.Marshal(payload)
		if err != nil {
			return submitResultMsg{ok: false, message: "Internal error building the request."}
		}

		client := &http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return submitResultMsg{ok: false, message: "Internal error creating the request."}
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return submitResultMsg{ok: false, message: "Network error -- try again, or email me directly."}
		}
		defer resp.Body.Close()

		var result struct {
			OK          bool   `json:"ok"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return submitResultMsg{ok: false, message: "Unexpected response -- try again, or email me directly."}
		}
		if !result.OK {
			return submitResultMsg{ok: false, message: "Telegram API error: " + result.Description}
		}
		return submitResultMsg{ok: true, message: "Sent! I'll get back to you soon."}
	}
}

func (f contactForm) view(s theme.Styles) string {
	var b strings.Builder

	if !f.editing {
		b.WriteString(s.Muted.Render("Press Enter to fill out a message."))
		return b.String()
	}

	labelStyle := func(active bool) lipgloss.Style {
		if active {
			return s.Accent
		}
		return s.Subtitle
	}

	row := func(label string, active bool, value string) string {
		return fmt.Sprintf("%s\n%s\n", labelStyle(active).Render(label), value)
	}

	b.WriteString(row("Name", f.focus == fieldName, f.name.View()))
	b.WriteString("\n")
	b.WriteString(row("Email", f.focus == fieldEmail, f.email.View()))
	b.WriteString("\n")
	b.WriteString(row("Company (optional)", f.focus == fieldCompany, f.company.View()))
	b.WriteString("\n")
	b.WriteString(row("Message", f.focus == fieldMessage, f.message.View()))
	b.WriteString("\n")

	submitLabel := "[ Send Message ]"
	if f.status == statusSending {
		submitLabel = "[ Sending... ]"
	}
	if f.focus == fieldSubmit {
		b.WriteString(lipgloss.NewStyle().Foreground(theme.ColorBright).Background(theme.ColorPrimary).Padding(0, 1).Render(submitLabel))
	} else {
		b.WriteString(s.Subtitle.Render(submitLabel))
	}
	b.WriteString("\n\n")

	switch f.status {
	case statusSuccess:
		b.WriteString(s.Success.Render(f.statusMsg))
	case statusError:
		b.WriteString(lipgloss.NewStyle().Foreground(theme.ColorAccent).Render(f.statusMsg))
	}

	b.WriteString("\n")
	b.WriteString(s.Muted.Render("tab/shift+tab: next/prev field · enter: submit (on button) · esc: cancel"))

	return b.String()
}
