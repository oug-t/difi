<h1 align="center">difi</h1>
<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Bubble_Tea-E2386F?style=for-the-badge&logo=tea&logoColor=white" />
  <img src="https://img.shields.io/github/license/oug-t/difi?style=for-the-badge&color=2e3440" />
</p>

<p align="center">
  <strong>A calm, focused way to review Git diffs.</strong><br />
  Review code with clarity. Polish before you push.
</p>

<p align="center">
  <img src="https://via.placeholder.com/800x450.png?text=Showcase+Your+UI+Here" alt="difi demo" width="100%" />
</p>

## Why difi?

- ⚡️ **Blazing Fast** — Built in Go. Starts instantly.
- 🎨 **Semantic UI** — Split-pane layout with syntax highlighting and Nerd Font icons.
- 🧠 **Context Aware** — Opens your editor (nvim/vim) at the exact line you are reviewing.
- ⌨️ **Vim Native** — Navigate with `h j k l`. Zero mouse required.

## Installation

### Homebrew (macOS & Linux)

```bash
brew tap oug-t/difi
brew install difi
```

## Installation

### Go Install

```bash
go install github.com/oug-t/difi/cmd/difi@latest
```

### Manual (Linux / Windows)

- Download the binary from Releases and add it to your $PATH.

## Workflow

- Run difi in any Git repository.
- By default, it compares your current branch against main.

```bash
cd my-project
difi
```

## Controls

| Key           | Action                                       |
| ------------- | -------------------------------------------- |
| `Tab`         | Toggle focus between File Tree and Diff View |
| `j / k`       | Move cursor down / up                        |
| `h / l`       | Focus Left (Tree) / Focus Right (Diff)       |
| `e` / `Enter` | Edit file (opens editor at selected line)    |
| `?`           | Toggle help drawer                           |
| `q`           | Quit                                         |

## Contributing

- PRs are welcome!
- We use Bubble Tea for the TUI.

```bash
git clone https://github.com/oug-t/difi
cd difi
go run cmd/difi/main.go
```

---

<p align="center"> Made with ❤️ by <a href="https://github.com/oug-t">oug-t</a> </p>
