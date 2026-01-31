<a id="readme-top"></a>
<h1 align="center"><code>difi</code></h1>
<p align="center"><em>Review and refine Git diffs before you push</em></p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Bubble_Tea-E2386F?style=for-the-badge&logo=tea&logoColor=white" />
  <img src="https://img.shields.io/github/license/oug-t/difi?style=for-the-badge&color=2e3440" />
</p>

<img width="2560" height="1440" alt="image" src="https://github.com/user-attachments/assets/fbea297a-b99d-4e98-b369-2925a7651a13" />

## Why difi?

- ⚡️ **Instant startup** — Built in Go, no background daemon.
- 🎨 **Structured review** — Tree view + side-by-side diffs.
- 🧠 **Editor-aware** — Jump to the exact line in `nvim` / `vim`.
- ⌨️ **Keyboard-first** — Designed for `h j k l`, no mouse.

## Why not `git diff`?

- `git diff` is powerful, but it’s optimized for output — not review.
- difi is designed for the *moment before you push or open a PR*:

## Installation

### Homebrew (macOS & Linux)

```bash
brew tap oug-t/difi
brew install difi
```

### Go Install

```bash
go install github.com/oug-t/difi/cmd/difi@latest
```

### Manual (Linux / Windows)

- Download the binary from Releases and add it to your `$PATH`.

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
<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Contributing

```bash
git clone https://github.com/oug-t/difi
cd difi
go run cmd/difi/main.go
```

Contributions are especially welcome in:
- diff rendering edge cases
- UI polish and accessibility
- Windows support
<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Star History

<a href="https://star-history.com/#oug-t/difi&Date">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=oug-t/difi&type=Date&theme=dark" />
      <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=oug-t/difi&type=Date" />
      <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=oug-t/difi&type=Date" />
    </picture>
  </a>
</div>
<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

<p align="center"> Made with ❤️ by <a href="https://github.com/oug-t">oug-t</a> </p>






