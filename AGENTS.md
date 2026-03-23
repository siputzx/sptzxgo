# AGENTS.md

This document provides essential information for agentic coding assistants working on the **sptzx** project.

## 🛠 Build, Lint, and Test Commands

The project is written in Go (Golang) and uses standard Go tooling.

### Build and Run
- **Install Dependencies:** `go mod tidy`
- **Build:** `go build .`
- **Run:** `go run main.go`

### Linting
The project follows standard Go formatting and linting:
- **Format Code:** `go fmt ./...`
- **Static Analysis:** `go vet ./...`

### Testing
- **Run All Tests:** `go test ./...`
- **Run Tests in a Package:** `go test -v ./src/path/to/package`
- **Run a Specific Test:** `go test -v -run TestFunctionName ./src/path/to/package`
- **Run with Coverage:** `go test -cover ./...`

---

## 🎨 Code Style Guidelines

### Formatting and Indentation
- **Indentation:** ALWAYS use **tabs** for indentation (Standard Go style).
- **Line Length:** Avoid excessively long lines; prefer wrapping at ~100-120 characters.
- **Braces:** Use Egyptian brackets (opening brace on the same line).

### Naming Conventions
- **Packages:** Use lowercase, single-word names (e.g., `core`, `config`, `serialize`).
- **Exported Identifiers:** Use `PascalCase` for exported functions, structs, and fields.
- **Internal Identifiers:** Use `camelCase` for non-exported variables, functions, and struct fields.
- **Interfaces:** Usually end in `-er` (e.g., `Handler`, `Provider`).

### Imports
Group imports into three blocks separated by a newline:
1. Standard library imports.
2. Third-party library imports (e.g., `go.mau.fi/whatsmeow`).
3. Local project imports (prefixed with `sptzx/`).

Example:
```go
import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow/types"

	"sptzx/src/core"
)
```

### Error Handling
- **Check Every Error:** Never ignore errors. Use `if err != nil { return err }`.
- **Contextual Errors:** Use `fmt.Errorf("failed to do X: %w", err)` to provide context when returning errors.
- **Panic:** Avoid using `panic()` unless it's a truly unrecoverable initialization error in `main.go`.

### Concurrency
- Use `context.Context` for cancellation and timeouts.
- Be careful with shared state; use `sync.Mutex` or channels to ensure thread safety.

---

## 🏗 Project Architecture

- **`main.go`**: The entry point. Initializes the database, WhatsApp client, and event handlers.
- **`src/core/`**: Core logic including the `Bot` struct, command registry, and the `Ptz` (context) struct.
- **`src/commands/`**: Contains command handlers categorized by folder.
- **`src/handler/`**: Logic for processing incoming events from WhatsApp.
- **`src/serialize/`**: Utility functions for sending messages and formatting data.
- **`src/config/`**: Configuration management via environment variables.

---

## ⚡ Adding New Commands

Commands are registered automatically via `init()` functions.

### Command Structure
Create a new file in `src/commands/<category>/<command_name>.go`:

```go
package category

import "sptzx/src/core"

func init() {
	core.Use(&core.Command{
		Name:        "example",
		Aliases:     []string{"ex", "test"},
		Description: "An example command",
		Usage:       "example <args>",
		Category:    "general", // Should match the folder name
		Handler: func(ptz *core.Ptz) error {
			if ptz.RawArgs == "" {
				return ptz.ReplyText("Please provide an argument.")
			}
			return ptz.ReplyText("You said: " + ptz.RawArgs)
		},
	})
}
```

### Useful `Ptz` Methods
- `ptz.ReplyText(text)` - Send a text reply to the current chat.
- `ptz.ReplyImage(data, mime, caption)` - Send an image reply.
- `ptz.React(emoji)` - Add a reaction to the message.
- `ptz.IsOwner()` - Returns true if the sender is a bot owner.
- `ptz.IsGroup` - Returns true if the message is from a group.
- `ptz.Args` - `[]string` of arguments.
- `ptz.RawArgs` - `string` of everything after the command.

---

## 🔒 Security and Best Practices
- **API Keys:** Never hardcode API keys or secrets. Use environment variables and load them via `src/config`.
- **Validation:** Always validate user input (args) before processing, especially for commands that interact with external APIs or the filesystem.
- **Performance:** Avoid heavy processing in the main event loop; use goroutines for long-running tasks if necessary.

---

## 🤖 Rule Files
No specific `.cursorrules` or `.github/copilot-instructions.md` were found in this repository. Follow the guidelines in this `AGENTS.md` file as the primary source of truth for code style and project conventions.
