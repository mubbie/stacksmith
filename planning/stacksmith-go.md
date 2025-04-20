# 🧑🏾‍🏭 Stacksmith Go - Project Plan

## 🔥 Mission

Build a beautiful, fast, expressive, and fully native Go CLI tool that helps developers manage stacked Git branches using vanilla Git — crafted for both power users and first-timers, with rich UI powered by Bubble Tea, Bubbles, and Lip Gloss.

## 🌟 Problem Statement

Modern software development workflows often involve creating multiple branches "stacked" on top of each other. This is especially common in large codebases where changes need to be broken down into smaller, reviewable pieces. However, managing these stacked branches can be challenging:

- Rebasing a branch requires manually rebasing all dependent branches
- Moving a branch to a new base requires careful git commands
- Visualizing the relationships between branches isn't straightforward
- Keeping track of push/upstream status adds complexity

Shell scripts can help, but they often lack intuitive UIs and are difficult to extend. There's a need for a tool that combines:

1. The power and flexibility of Git
2. An intuitive, expressive interface
3. Support for both interactive and headless (script-friendly) usage
4. Native performance without external dependencies

## 🎯 Solution

Stacksmith Go is a command-line tool built in Go that manages stacked Git branches with an elegant UI. By leveraging Go's performance and the charm.sh libraries (Bubble Tea, Bubbles, and Lip Gloss), it provides:

1. A beautiful, space-efficient UI for interactive use
2. Command-specific prompts for guided operations
3. Traditional CLI capabilities for scripting and automation
4. A full-screen TUI mode for advanced visualization and management

## 🎨 UX Design Philosophy

| Behavior                      | Result                                 |
|-------------------------------|----------------------------------------|
| `stacksmith`                  | 🎛 Launches simple menu UI             |
| `stacksmith <command> <args>` | ⚙️ Runs the command headlessly         |
| `stacksmith <command>` (no args) | 💡 Launches command-specific UI prompt |
| `stacksmith tui`              | 🖥 Full-screen TUI (branch graph)      |

This design ensures:
- Smooth onboarding through interactive menus
- Scriptable CLI commands for advanced usage
- Rich visualization for complex branch relationships

## 🌳 Supported Features (Phase 1)

| Command | Description                      | UI Mode                    |
|---------|----------------------------------|----------------------------|
| stack   | Create a new branch atop another | Interactive prompt         |
| sync    | Rebase multiple branches sequentially | Interactive prompt     |
| fix-pr  | Rebase one branch onto a new base | Interactive prompt        |
| push    | Smart push with upstream detection | None (headless)          |
| graph   | Show commit graph (git log --graph) | None (headless)         |
| tui     | Full-screen DAG browser and navigator | Full Bubble Tea TUI   |
| (default) | Interactive menu to run any command | Simple menu UI        |

## 🧱 Project Structure

```
stacksmith/
├── cmd/                       # CLI commands via Cobra
│   ├── root.go                # stacksmith → launches simple menu
│   ├── stack.go               # stack command with interactive prompt
│   ├── sync.go                # sync command implementation
│   ├── fixpr.go               # fix-pr command implementation
│   ├── push.go                # push command implementation
│   ├── graph.go               # graph command implementation
│   └── tui.go                 # tui command (full-screen UI)
│
├── core/                      # Git + stack logic
│   ├── git.go                 # Git command executor
│   └── stack.go               # Optional stack helpers
│
├── ui/
│   ├── simplemenu/            # Simple menu implementation
│   │   ├── menu.go            # Main menu UI
│   │   └── stack_prompt.go    # Stack command prompt
│   └── fulltui/               # Full-screen TUI (Phase 3)
│       ├── model.go           # App state model
│       ├── update.go          # Key/message handlers
│       ├── view.go            # UI rendering logic
│       └── styles.go          # Lip Gloss styling
│
├── render/                    # Output styling utilities
│   └── printer.go             # Console output formatting
│
├── scripts/                   # Bash fallback tool
│   └── stacksmith-lite.sh     # Original shell script
│
├── planning/
│   └── stacksmith-go.md       # Notes & planning
│
├── go.mod                     # Go module dependencies
├── main.go                    # Entry point
└── README.md                  # Documentation
```

## 🛠 Implementation Phases

### ✅ Phase 1: Foundation (Complete)
- Go module setup
- Cobra CLI scaffold
- Git command execution engine
- Basic command wiring (stack, push, etc.)
- Simple menu UI with Bubble Tea

### 🚧 Phase 2: Command-Specific Interfaces (In Progress)
- Interactive prompts for stack, sync, fix-pr
- Simple menu UI for main command selection
- Command execution from menu selections
- Rich console output with emojis and color

### ⏳ Phase 3: Fullscreen TUI (Planned)
- Interactive DAG viewer
- Scrollable branch list
- Node actions: rebase, push, retarget
- Keybindings + footer instructions

### ⏳ Phase 4: Output & Polish (Planned)
- Enhanced printer utilities
- System feedback (✓, ✗, emoji lines)
- Configurable themes
- Error handling improvements

## 💻 Technical Approach

### UI Architecture
We're utilizing a hybrid approach:

1. **Simple Menu**: For the main command selection, we use a minimal Bubble Tea UI that:
   - Shows a list of available commands
   - Provides descriptions and emoji icons
   - Takes minimal screen space (not full-screen)
   - Returns to the terminal after selection

2. **Command Prompts**: For commands needing input, we show targeted prompts
   - Field-based input with validation
   - Contextual guidance
   - Keyboard navigation

3. **Full-screen TUI**: Reserved for the dedicated `tui` command
   - Branch visualization
   - Interactive navigation
   - Rich interaction model

### Git Integration
All Git operations are performed through a dedicated `GitExecutor` that:
- Executes Git commands safely
- Handles errors properly
- Provides high-level abstractions for common operations

## 🧪 Testing Strategy
- 🔁 Unit tests for core Git commands
- 🧪 Integration tests using temp git repos
- 📸 Golden snapshot tests for TUI
- ⚙️ CI with GitHub Actions (macOS + Linux)

## 🚀 Current Status

Phase 1 is complete, and Phase 2 is in progress:

- ✅ CLI foundation with Cobra
- ✅ Git execution engine
- ✅ Main menu UI with simplemenu
- ✅ Stack command with interactive prompt
- 🚧 Other command-specific prompts
- ⏳ Full-screen TUI implementation

## 📦 Future Ideas

| Feature   | Description                           |
|-----------|---------------------------------------|
| diff      | Show what each stacked PR contains    |
| status    | Health + sync state of the stack      |
| snapshot  | Save + restore branch stack state     |
| Git API   | Auto-retarget PRs on GitHub/Azure     |

## 🧑🏾‍🍳 Final Thoughts

Stacksmith Go is:
- 🌱 Lightweight when you need speed
- 🌲 Visual when you want structure
- 🧑🏾‍🏭 Artisan-crafted for joy and clarity
- ⚙️ Powered by native Go, no shell dependencies