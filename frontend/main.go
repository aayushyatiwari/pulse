package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const sockPath = "/tmp/pulsed.sock"

// ── styles ────────────────────────────────────────────────────────────────────

var (
	header = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("13")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("8")).
		Width(80)

	inputBar = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("13")).
		Padding(0, 1)

	namePrompt = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("13"))

	// message part styles
	nameStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))  // bright blue
	ipStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))              // dark gray
	tsStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))             // amber
	msgStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))             // white

	// name (ip) at time: "message"
	msgRe = regexp.MustCompile(`^(.+) \(([^)]+)\) at (\d{2}:\d{2}:\d{2}): (.*)$`)
)

// ── tea messages ──────────────────────────────────────────────────────────────

type (
	connectedMsg struct {
		conn    net.Conn
		scanner *bufio.Scanner
	}
	lineMsg string
	errMsg  error
)

// ── model ─────────────────────────────────────────────────────────────────────

type state int

const (
	stateNaming state = iota
	stateChatting
)

type model struct {
	state   state
	name    string
	conn    net.Conn
	scanner *bufio.Scanner
	vp      viewport.Model
	input   textinput.Model
	lines   []string
	ready   bool
	width   int
	height  int
}

func newModel() model {
	ti := textinput.New()
	ti.Focus()

	name := loadName()
	if name == "" {
		ti.Placeholder = "your name…"
		return model{state: stateNaming, input: ti}
	}

	ti.Placeholder = "message…"
	return model{state: stateChatting, name: name, input: ti}
}

// ── init ──────────────────────────────────────────────────────────────────────

// colorLine parses "name (ip) at HH:MM:SS: message" and applies per-part colors.
// Falls back to plain text if the format doesn't match.
func colorLine(line string) string {
	m := msgRe.FindStringSubmatch(line)
	if m == nil {
		return line
	}
	return nameStyle.Render(m[1]) +
		" (" + ipStyle.Render(m[2]) + ")" +
		" at " + tsStyle.Render(m[3]) + ": " +
		msgStyle.Render(m[4])
}

func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{textinput.Blink}
	if m.state == stateChatting {
		cmds = append(cmds, dialCmd())
	}
	return tea.Batch(cmds...)
}

// ── update ────────────────────────────────────────────────────────────────────

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		header = header.Width(m.width)
		vpH := m.height - 4 // header(2) + input(2)
		if !m.ready {
			m.vp = viewport.New(m.width, vpH)
			m.ready = true
		} else {
			m.vp.Width, m.vp.Height = m.width, vpH
		}

	case connectedMsg:
		m.conn, m.scanner = msg.conn, msg.scanner
		cmds = append(cmds, readCmd(m.scanner))

	case lineMsg:
		line := string(msg)
		if !strings.HasPrefix(line, "Enter your name:") {
			m.lines = append(m.lines, colorLine(line))
			m.vp.SetContent(strings.Join(m.lines, "\n"))
			m.vp.GotoBottom()
		}
		cmds = append(cmds, readCmd(m.scanner))

	case errMsg:
		m.lines = append(m.lines, fmt.Sprintf("⚠  %v", msg))
		if m.ready {
			m.vp.SetContent(strings.Join(m.lines, "\n"))
			m.vp.GotoBottom()
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEnter:
			text := strings.TrimSpace(m.input.Value())
			if text == "" {
				return m, nil
			}
			m.input.Reset()

			if m.state == stateNaming {
				m.name = text
				saveName(text)
				m.state = stateChatting
				m.input.Placeholder = "message…"
				return m, tea.Batch(dialCmd(), textinput.Blink)
			}

			fmt.Fprintf(m.conn, "PULSE:%s|%s\n", m.name, text)
			return m, nil
		}
	}

	var c1, c2 tea.Cmd
	m.input, c1 = m.input.Update(msg)
	m.vp, c2 = m.vp.Update(msg)
	return m, tea.Batch(append(cmds, c1, c2)...)
}

// ── view ──────────────────────────────────────────────────────────────────────

func (m model) View() string {
	if m.state == stateNaming {
		return "\n\n  " + namePrompt.Render("pulse 💬") +
			"\n\n  " + m.input.View() +
			"\n\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("press enter to join")
	}
	if !m.ready {
		return "\n  connecting…"
	}
	return fmt.Sprintf(
		"%s\n%s\n%s",
		header.Render("  pulse 💬"),
		m.vp.View(),
		inputBar.Width(m.width-4).Render(m.input.View()),
	)
}

// ── commands ──────────────────────────────────────────────────────────────────

func dialCmd() tea.Cmd {
	return func() tea.Msg {
		conn, err := net.Dial("unix", sockPath)
		if err != nil {
			return errMsg(err)
		}
		return connectedMsg{conn: conn, scanner: bufio.NewScanner(conn)}
	}
}

func readCmd(sc *bufio.Scanner) tea.Cmd {
	return func() tea.Msg {
		if sc.Scan() {
			return lineMsg(sc.Text())
		}
		return errMsg(fmt.Errorf("disconnected from daemon"))
	}
}

// ── config ────────────────────────────────────────────────────────────────────

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "pulse", "name.txt")
}

func loadName() string {
	b, err := os.ReadFile(configPath())
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func saveName(name string) {
	p := configPath()
	os.MkdirAll(filepath.Dir(p), 0755)
	os.WriteFile(p, []byte(name), 0644)
}

// ── main ──────────────────────────────────────────────────────────────────────

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
