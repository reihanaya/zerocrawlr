Here is the exact content for your `README.md` file. You can copy everything below and paste it directly into a file named `README.md` in your project folder:

```markdown
# zerocrawlr

A lightweight, stealthy, terminal-based web crawler written in Go. 
`zerocrawlr` utilizes a headless Chromium browser to bypass basic anti-bot protections, extract all visible links from a target webpage, and export them into a clean JSON file.

![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)

## ✨ Features

*   **Beautiful TUI:** Clean, interactive terminal interface built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lipgloss](https://github.com/charmbracelet/lipgloss), and [Bubbles](https://github.com/charmbracelet/bubbles).
*   **Browser Emulation:** Uses [chromedp](https://github.com/chromedp/chromedp) to drive a real headless browser, perfectly executing JavaScript and rendering dynamic content (React, Vue, Angular).
*   **Anti-Bot Evasion:** Bypasses basic bot detection mechanisms by:
    *   Rotating modern, realistic User-Agents.
    *   Stripping the `AutomationControlled` flag from Chromium.
    *   Ignoring SSL certificate mismatches.
*   **Smart Scraping:** Automatically scrolls to the bottom of the page before scraping, triggering lazy-loaded content.
*   **Clean Export:** Outputs neatly formatted JSON files (e.g., `2023-10-27_example_com.json`).

## 🛠️ Prerequisites

Because `zerocrawlr` relies on Chrome DevTools Protocol, you **must have Google Chrome or Chromium installed** on your system. 

*   **Go:** v1.21 or higher recommended.
*   **Chrome/Chromium:** Installed and available in your system's `PATH`.

## 📦 Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/reihanaya/zerocrawlr.git
cd zerocrawlr
go mod tidy
go build -o zerocrawlr
```

## 🚀 Usage

Run the compiled binary:

```bash
./zerocrawlr
```

1. Type the target URL (e.g., `https://example.com`).
2. Press `Enter` to start the crawler.
3. Wait for the headless browser to emulate, scroll, and extract the data.
4. Press `Esc` at any time to safely quit.

## 📄 Output Example

If you target `https://example.com`, `zerocrawlr` will generate a file named something like `2023-10-27_example_com.json`:

```json
[
  {
    "title": "More information...",
    "url": "https://www.iana.org/domains/example"
  },
  {
    "title": "Example Domain",
    "url": "https://www.iana.org/domains/reserved"
  }
]
```

## ⚙️ How it Works

1.  **Initialization:** Sets up a headless Chrome instance with specific flags (`--no-sandbox`, `--disable-gpu`, `--ignore-certificate-errors`).
2.  **Stealth:** Selects a random User-Agent from a pool of up-to-date browsers (Chrome, Firefox, Safari) and disables the `AutomationControlled` feature.
3.  **Navigation:** Navigates to the target URL and waits for the `body` tag to become visible.
4.  **Interaction:** Executes a JavaScript snippet to scroll to the absolute bottom of the page and pauses for 3 seconds. This allows time for AJAX calls and lazy-loading scripts to populate the DOM.
5.  **Extraction:** Queries the DOM for all `<a>` tags with `href` attributes and valid inner text, mapping them to a structured object.
6.  **Export:** Marshals the data into formatted JSON and writes it to a timestamped file.

## ⚠️ Disclaimer

This tool is intended for educational purposes and legitimate web scraping (like extracting public links). Always respect a website's `robots.txt` file and Terms of Service. The author is not responsible for any misuse of this software.

## 📜 License

This project is licensed under the MIT License - see the `LICENSE` file for details.
```