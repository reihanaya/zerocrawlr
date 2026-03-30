```markdown
# zerocrawlr

**zerocrawlr** is a powerful, interactive Command-Line Interface (CLI) web crawler built with Go. 
It goes beyond simple HTML parsing by utilizing headless browser emulation to execute JavaScript, trigger lazy-loaded elements, and bypass basic anti-bot protections.

## ✨ Features

- **Interactive TUI**: Built with Bubble Tea for a smooth, terminal-based user interface.
- **Headless Browser Emulation**: Uses `chromedp` to render Single Page Applications (SPAs) and dynamic content.
- **Aggressive Scraping**: Automatically scrolls the page to trigger lazy-loaded elements before scraping.
- **Evasion Tactics**: Includes User-Agent randomization, SSL error bypassing, and bot-detection evasion flags.
- **Auto-Export**: Automatically formats the output into a clean JSON file named by date and target URL (e.g., `2026-03-30_example_com.json`).

## 🚀 Installation

Ensure you have [Go](https://golang.org/dl/) installed on your machine.

1. Clone the repository:
   ```bash
   git clone [https://github.com/yourusername/zerocrawlr.git](https://github.com/yourusername/zerocrawlr.git)
   cd zerocrawlr
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the executable:
   ```bash
   go build -o zerocrawlr main.go
   ```

## 🛠️ Usage

Run the compiled binary from your terminal:

```bash
./zerocrawlr
```

1. The TUI will prompt you to enter a **Target URL**.
2. Press `Enter`.
3. Wait for the browser emulation to finish processing the page.
4. Check your current directory for the generated JSON file containing all the extracted links and titles.

## 📦 Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Bubbles](https://github.com/charmbracelet/bubbles) - For the TUI.
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - For terminal styling.
- [Chromedp](https://github.com/chromedp/chromedp) - For headless browser automation.

## ⚠️ Disclaimer

This tool is created for educational purposes. Please respect the `robots.txt` of the websites you crawl and ensure you are not violating their Terms of Service.
```
