# Histograph

A terminal-based browser history visualizer for Chrome and Firefox, built in Go using Bubble Tea.

## Features
- Interactive TUI for visualizing recent browser history
- Supports Chrome and Firefox on Linux, macOS, and Windows
- Multiple views: Overview, Timeline, Top Sites, Details
- Auto-detects browser history paths, with environment variable overrides
- User-friendly error handling and cross-platform support

## Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/akshatsrivastava11/Histograph.git
   cd Histograph
   ```
2. **Build the project:**
   ```sh
   go build -o histograph ./cmd/Histograph
   ```

## Usage

Run the application:
```sh
./histograph
```

- Select your browser (Chrome or Firefox) from the menu.
- Interact with the TUI using the following keys:
  - `1`/`2`/`3`/`4`: Switch between Overview, Timeline, Top Sites, Details
  - `↑`/`↓`: Navigate entries
  - `q`: Quit

## Configuration

By default, Histograph auto-detects browser history file locations. You can override these with environment variables:

- `CHROME_HISTORY_PATH`: Path to Chrome's `History` SQLite file
- `FIREFOX_HISTORY_PATH`: Path to Firefox's `places.sqlite` file

Example:
```sh
CHROME_HISTORY_PATH=/custom/path/History ./histograph
```

## Cross-Platform Support
- **Linux:**
  - Chrome: `~/.config/google-chrome/Default/History`
  - Firefox: `~/.mozilla/firefox/<profile>/places.sqlite`
- **macOS:**
  - Chrome: `~/Library/Application Support/Google/Chrome/Default/History`
  - Firefox: `~/Library/Application Support/Firefox/Profiles/<profile>/places.sqlite`
- **Windows:**
  - Chrome: `%USERPROFILE%\AppData\Local\Google\Chrome\User Data\Default\History`
  - Firefox: `%USERPROFILE%\AppData\Roaming\Mozilla\Firefox\Profiles\<profile>\places.sqlite`

## Development & Testing
- Run tests:
  ```sh
  go test ./tests/...
  ```
- Code is organized in `internals/` by feature (parse, render, types).

## Contributing
Pull requests and issues are welcome! Please:
- Write clear commit messages
- Add tests for new features
- Document your code

## License
MIT 