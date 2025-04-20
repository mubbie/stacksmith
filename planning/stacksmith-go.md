# 🧑🏾‍🏭 Stacksmith Go - Features & Roadmap

## 🌳 Feature Status (Core Features)

| Feature | Description | Status |
|---------|-------------|--------|
| Main Menu UI | Interactive command selection | ✅ Implemented |
| Stack Command | Create a new branch | ✅ Implemented |
| Sync Command | Rebase multiple branches sequentially | ✅ Implemented |
| Fix-PR Command | Rebase one branch onto a new base | ✅ Implemented |
| Push Command | Smart push with upstream detection | ✅ Implemented |
| Graph Command | Show commit graph | ✅ Implemented |
| Error Handling | Improved error detection and recovery | ✅ Implemented |
| Post-Command Options | Option to return to menu or exit | ✅ Implemented |
| Terminal UI (TUI) | Full-screen branch visualization | 🚧 In Progress |
| UI Component Sharing | Shared UI components and styles | ✅ Implemented |

## 🧩 Feature Status (Extended Features)

| Feature | Description | Status |
|---------|-------------|--------|
| Configuration Management | User-configurable settings | ⏳ Planned |
| Testing Framework | Unit and integration tests | ⏳ Planned |
| Logging System | Structured logging | ⏳ Planned |
| Branch History Tracking | Track stack membership | ⏳ Planned |
| Shell Auto-completion | Tab completion for commands | ⏳ Planned |
| Cache System | Improved performance via caching | ⏳ Planned |
| Diff Viewer | Show what each PR contains | ⏳ Planned |
| Stack Status | Health + sync state of stack | ⏳ Planned |
| Stack Snapshot | Save + restore branch state | ⏳ Planned |
| Git API Integration | Auto-retarget PRs on GitHub/Azure | ⏳ Planned |
| Plugin System | Extensibility via plugins | ⏳ Planned |

## 🛠 Implementation Roadmap

### ✅ Phase 1: Foundation (Complete)
- Go module setup
- Cobra CLI scaffold
- Git command execution engine
- Basic command wiring (stack, push, etc.)
- Simple menu UI with Bubble Tea
- Shared UI components and styles
- Improved error handling

### 🚧 Phase 2: Terminal UI (TUI) (Priority)
- Interactive DAG viewer with branch visualization
- Scrollable branch list with actions
- Node interactions: checkout, rebase, push, etc.
- Stack visualization with parent-child relationships
- Status indicators for sync state, ahead/behind counts
- Keyboard shortcuts and footer help
- Branch diff viewer

### ⏳ Phase 3: Project Robustness (Next)
- Configuration management with user preferences
- Logging system with debug levels
- Command-line global flags
- Testing framework (unit + integration)
- Error recovery and guided fixes
- Shell auto-completion
- Git operation caching

### ⏳ Phase 4: Advanced Features (Future)
- Branch history tracking and stack membership
- Stack health status monitoring
- Stack snapshot and restore
- Git hosting platform API integration
- Plugin system for extensions
- Custom themes and styling

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
│   ├── stack.go               # Stack data structures and operations
│   ├── errors.go              # Custom error types
│   └── cache.go               # Git operation caching
│
├── config/                    # Configuration management
│   └── config.go              # User configuration
│
├── ui/
│   ├── styles/                # Shared UI styling
│   │   └── styles.go          # Centralized Lip Gloss styling
│   ├── simplemenu/            # Simple menu implementation
│   │   ├── menu.go            # Main menu UI
│   │   ├── baseprompt.go      # Base prompt structure
│   │   ├── list.go            # Reusable list component
│   │   ├── stack_prompt.go    # Stack command prompt
│   │   ├── sync_prompt.go     # Sync command prompt
│   │   └── fixpr_prompt.go    # Fix-PR command prompt
│   └── fulltui/               # Full-screen TUI
│       ├── model.go           # App state model
│       ├── update.go          # Key/message handlers
│       ├── view.go            # UI rendering logic
│       ├── styles.go          # Lip Gloss styling
│       ├── dag.go             # DAG visualization component
│       └── statusbar.go       # Status bar component
│
├── render/                    # Output styling utilities
│   └── printer.go             # Console output formatting
│
├── logger/                    # Logging system
│   └── logger.go              # Structured logging
│
├── plugins/                   # Plugin architecture
│   └── plugin.go              # Plugin interface
│
├── scripts/                   # Bash fallback tool
│   └── stacksmith-lite.sh     # Original shell script
│
├── test/                      # Testing utilities
│   └── git_fixture.go         # Git test fixtures
```