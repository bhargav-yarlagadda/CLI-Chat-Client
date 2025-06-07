package commands

import (
	"fmt"
	"os"
	"time"

	"cli-chat-client/api"
	"cli-chat-client/data"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

/* ───────────────────────────── CONSTANTS ───────────────────────────── */

const backendBaseURL = "http://localhost:8000" // <-- hard-coded for now

/* ───────────────────────────── TYPES ───────────────────────────────── */

type friendRequest struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (f friendRequest) Title() string       { return f.From }
func (f friendRequest) Description() string { return fmt.Sprintf("Status: %s | Created: %s", f.Status, f.CreatedAt.Format("2006-01-02 15:04:05")) }
func (f friendRequest) FilterValue() string { return f.From }

/* ───────────────────────────── MODEL ───────────────────────────────── */

type model struct {
	list     list.Model
	requests []friendRequest
	err      error
	quitting bool
}

func NewNotifyModel(requests []friendRequest) model {
	items := make([]list.Item, len(requests))
	for i, r := range requests {
		items[i] = r
	}

	const defaultWidth, defaultHeight = 50, 10
	l := list.New(items, list.NewDefaultDelegate(), defaultWidth, defaultHeight)
	l.Title = "Pending Friend Requests (↑/↓ navigate, ← reject, → accept, q quit)"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return model{list: l, requests: requests}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.quitting {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("left"))):
			i := m.list.Index()
			if m.requests[i].Status == "pending" {
				if err := sendResponse(m.requests[i].ID, false); err == nil {
					m.requests[i].Status = "rejected"
					m.list.SetItem(i, m.requests[i])
				} else {
					m.err = err
				}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("right"))):
			i := m.list.Index()
			if m.requests[i].Status == "pending" {
				if err := sendResponse(m.requests[i].ID, true); err == nil {
					m.requests[i].Status = "accepted"
					m.list.SetItem(i, m.requests[i])
				} else {
					m.err = err
				}
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-2)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}
	if len(m.requests) == 0 {
		return "No pending friend requests.\n"
	}
	if m.quitting {
		return "Exiting notification viewer.\n"
	}
	return m.list.View()
}

/* ─────────────────────────── HELPERS ───────────────────────────────── */

func sendResponse(requestID string, accept bool) error {
	if data.JWT_TOKEN == "" {
		return fmt.Errorf("please login first")
	}

	url := backendBaseURL + "/friend-request/respond"

	resp, err := resty.New().R().
		SetAuthToken(data.JWT_TOKEN).
		SetBody(map[string]interface{}{
			"request_id": requestID,
			"accept":     accept,
		}).
		Post(url)

	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("server returned %s", resp.Status())
	}
	return nil
}

func convertAPIRequests(apiReq []api.FriendRequest) []friendRequest {
	out := make([]friendRequest, len(apiReq))
	for i, r := range apiReq {
		out[i] = friendRequest{
			ID:        r.ID,
			From:      r.From,
			Status:    r.Status,
			CreatedAt: r.CreatedAt,
		}
	}
	return out
}

/* ─────────────────────────── ENTRY POINT ───────────────────────────── */

func Notify() {
	apiRequests, err := api.GetAllNotifications() // this still uses whatever URL logic lives in api package
	if err != nil {
		fmt.Println("Error fetching notifications:", err)
		return
	}
	if len(apiRequests) == 0 {
		fmt.Println("No pending friend requests.")
		return
	}

	m := NewNotifyModel(convertAPIRequests(apiRequests))

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running notification viewer:", err)
		os.Exit(1)
	}
}
