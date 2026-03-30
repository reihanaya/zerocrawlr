package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chromedp/chromedp"
)

// Styling using Lipgloss
var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99")).MarginBottom(1)
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
)

// List of modern, realistic User-Agents to bypass basic anti-bot systems
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
}

// Helper function to pick a random User-Agent
func getRandomUserAgent() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return userAgents[r.Intn(len(userAgents))]
}

// ScrapedData stores the crawl results
type ScrapedData struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Message sent when the crawler finishes
type crawlFinishedMsg struct {
	err   error
	count int
	file  string
}

// Model represents the application state
type model struct {
	urlInput textinput.Model
	spinner  spinner.Model
	crawling bool
	done     bool
	err      error
	count    int
	filename string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "https://example.com"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50
	ti.PromptStyle = focusedStyle
	ti.TextStyle = focusedStyle
	ti.Cursor.Style = focusedStyle

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return model{
		urlInput: ti,
		spinner:  s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case crawlFinishedMsg:
		m.crawling = false
		m.done = true
		m.err = msg.err
		m.count = msg.count
		m.filename = msg.file
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			if !m.crawling && !m.done {
				m.crawling = true
				targetURL := m.urlInput.Value()
				if targetURL == "" {
					targetURL = "https://example.com"
				}

				return m, tea.Batch(m.spinner.Tick, startCrawling(targetURL))
			}
		}
	}

	if !m.crawling && !m.done {
		var cmd tea.Cmd
		m.urlInput, cmd = m.urlInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.done {
		if m.err != nil {
			return errorStyle.Render(fmt.Sprintf("\n[x] Error: %v\n", m.err))
		}
		return successStyle.Render(fmt.Sprintf("\n[v] Done! %d items exported to: %s\n", m.count, m.filename))
	}

	if m.crawling {
		return fmt.Sprintf("\n %s Emulating browser to crawl target, please wait...\n", m.spinner.View())
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("=== zerocrawlr ===") + "\n")
	b.WriteString("Enter Target URL:\n")
	b.WriteString(m.urlInput.View() + "\n\n")
	b.WriteString(blurredStyle.Render("Press Enter to start • Esc to quit"))

	return b.String()
}

func generateFilename(target string) string {
	dateStr := time.Now().Format("2006-01-02")
	cleanURL := strings.ReplaceAll(target, "https://", "")
	cleanURL = strings.ReplaceAll(cleanURL, "http://", "")
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	cleanURL = reg.ReplaceAllString(cleanURL, "_")
	cleanURL = strings.Trim(cleanURL, "_")
	return fmt.Sprintf("%s_%s.json", dateStr, cleanURL)
}

// startCrawling uses chromedp to emulate a real browser with a random User-Agent
func startCrawling(target string) tea.Cmd {
	return func() tea.Msg {
		if !strings.HasPrefix(target, "http") {
			target = "https://" + target
		}

		randomUA := getRandomUserAgent()

		// Set up context for chromedp with Anti-Bot & SSL evasion flags
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("ignore-certificate-errors", true),                // Bypass SSL mismatch errors
			chromedp.Flag("disable-blink-features", "AutomationControlled"), // Bypass basic bot detection
			chromedp.UserAgent(randomUA),
		)

		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		ctx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()

		// 60 seconds timeout
		ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		var rawLinks []map[string]string

		// The "Aggressive" Browser Actions
		err := chromedp.Run(ctx,
			chromedp.Navigate(target),
			chromedp.WaitVisible(`body`, chromedp.ByQuery),
			chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight);`, nil),
			chromedp.Sleep(3*time.Second),
			chromedp.Evaluate(`
				Array.from(document.querySelectorAll('a[href]')).map(a => ({
					title: a.innerText.trim(),
					url: a.href
				})).filter(item => item.url !== "" && item.title !== "")
			`, &rawLinks),
		)

		if err != nil {
			return crawlFinishedMsg{err: fmt.Errorf("browser emulation failed: %v", err)}
		}

		var results []ScrapedData
		for _, item := range rawLinks {
			results = append(results, ScrapedData{
				Title: item["title"],
				URL:   item["url"],
			})
		}

		outputFile := generateFilename(target)

		if len(results) > 0 {
			file, err := os.Create(outputFile)
			if err != nil {
				return crawlFinishedMsg{err: err}
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(results); err != nil {
				return crawlFinishedMsg{err: err}
			}
		}

		return crawlFinishedMsg{count: len(results), file: outputFile}
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Failed to run TUI: %v\n", err)
		os.Exit(1)
	}
}
