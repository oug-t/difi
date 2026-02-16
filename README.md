<a id="readme-top"></a>

<h1 align="center"><code>difi</code></h1>
<p align="center"><em>Review and refine Git diffs before you push</em></p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Bubble_Tea-E2386F?style=for-the-badge&logo=tea&logoColor=white" />
  <img src="https://img.shields.io/github/license/oug-t/difi?style=for-the-badge&color=2e3440" />
</p>

<p align="center">
  <img src= "https://github.com/user-attachments/assets/3695cfd2-148c-463d-9630-547d152adde0" alt="difi_demo" />
</p>

## Why difi?

**git diff** shows changes. **difi** helps you _review_ them.

- ⚡️ **Instant** — Built in Go. Launches immediately with no daemon or indexing.
- 🎨 **Structured** — A clean file tree and focused diffs for fast mental parsing.
- 🧠 **Editor-Aware** — Jump straight to the exact line in `nvim`/`vim` to fix issues.
- ⌨️ **Keyboard-First** — Navigate everything with `h j k l`. No mouse required.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Installation

#### Homebrew (macOS & Linux)

```bash
brew install difi
```

#### Go Install

```bash
go install github.com/oug-t/difi/cmd/difi@latest
```

#### AUR (Arch Linux)

**Binary (pre-built):**

```bash
pikaur -S difi-bin
```

**Build from source:**

```bash
pikaur -S difi
```

#### Manual (Linux / Windows)

- Download the binary from Releases and add it to your `$PATH`.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Workflow

- Run difi in any Git repository against main:

```bash
cd my-project
difi
```

**Piping & Alternative VCS**

- You can also pass raw diffs directly into `difi` via standard input. This is perfect for patch files or other version control systems like Jujutsu:

```bash
# Review a saved patch file
cat changes.patch | difi

# Review changes in Jujutsu (jj)
jj diff --git | difi

# Pipe standard git diff output
git diff | difi
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

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

## Integrations

#### vim-fugitive

- **The "Unix philosophy" approach:** Uses the industry-standard Git wrapper to provide a robust, side-by-side editing experience.
- **Side-by-Side Editing:** Instantly opens a vertical split (:Gvdiffsplit!) against the index.
- **Merge Conflicts:** Automatically detects conflicts and opens a 3-way merge view for resolution.
- **Config**: Add the line below to if using **lazy.nvim**.

```lua
{
  "tpope/vim-fugitive",
  cmd = { "Gvdiffsplit", "Git" }, -- Add this line
}
```

<p align="left"> 
  <a href="https://github.com/tpope/vim-fugitive.git">
    <img src="https://img.shields.io/badge/Supports-vim--fugitive-4d4d4d?style=for-the-badge&logo=vim&logoColor=white" alt="Supports vim-fugitive" />
  </a>
</p>

#### difi.nvim

Get the ultimate review experience with **[difi.nvim](https://github.com/oug-t/difi.nvim)**.

- **Auto-Open:** Instantly jumps to the file and line when you press `e` in the CLI.
- **Visual Diff:** Renders diffs inline with familiar green/red highlights—just like reviewing a PR on GitHub.
- **Interactive Review:** Restore a "deleted" line by simply removing the `-` marker. Discard an added line by deleting it entirely.
- **Context Aware:** Automatically syncs with your `difi` session target.

<p align="left">
  <a href="https://github.com/oug-t/difi.nvim">
    <img src="https://img.shields.io/badge/Get_difi.nvim-57A143?style=for-the-badge&logo=neovim&logoColor=white" alt="Get difi.nvim" />
  </a>
</p>

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Git Integration

To use `difi` as a native git command (e.g., `git difi`), add it as an alias in your global git config:

```bash
git config --global alias.difi '!difi'
```

Now you can run it directly from git:

```bash
git difi
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Contributing

```bash
git clone https://github.com/oug-t/difi
cd difi
go run cmd/difi/main.go
```

Contributions are especially welcome in:

- diff.nvim rendering edge cases
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
