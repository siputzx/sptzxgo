# sptzx

WhatsApp bot written in Go, built on [whatsmeow](https://github.com/tulir/whatsmeow).

---

## Requirements

- Go 1.23+
- FFmpeg
- SQLite3

## Setup

```bash
git clone https://github.com/siputzx/sptzx.git
cd sptzx
cp .env.example .env
```

Edit `.env`, then:

```bash
go mod tidy
go run main.go
```

## Configuration

| Variable | Required | Default | Description |
|---|---|---|---|
| `BOT_OWNERS` | ✅ | — | Owner numbers, comma-separated (e.g. `628xxx`) |
| `BOT_PREFIX` | ❌ | `!,.,/` | Command prefixes, comma-separated |
| `SESSION_DB` | ❌ | `file:session.db` | SQLite session path |
| `LOGIN_METHOD` | ❌ | `qr` | `qr` or `paircode` |
| `PAIRING_PHONE` | ❌ | — | Required if `LOGIN_METHOD=paircode` |
| `LOG_LEVEL` | ❌ | `INFO` | `DEBUG` `INFO` `WARN` `ERROR` |
| `SIPUTZX_ENABLED` | ❌ | `false` | Enable siputzx API |
| `SIPUTZX_BASE_URL` | ❌ | `https://api.siputzx.my.id` | External API base URL |

## Project Structure

```
sptzx/
├── main.go
├── .env.example
├── go.mod
└── src/
    ├── api/
    ├── config/
    ├── core/
    ├── handler/
    ├── middleware/
    ├── serialize/
    └── commands/
        ├── ai/
        ├── downloader/
        ├── games/
        ├── general/
        ├── group/
        ├── info/
        ├── maker/
        ├── owner/
        ├── primbon/
        ├── random/
        ├── search/
        ├── stalk/
        ├── sticker/
        └── tools/
```

## Adding a Command

Create a new file in the appropriate category folder:

```go
package general

import "sptzx/src/core"

func init() {
    core.Use(&core.Command{
        Name:        "hello",
        Aliases:     []string{"hi"},
        Description: "Greet the user",
        Usage:       "hello",
        Category:    "general",
        Handler: func(ptz *core.Ptz) error {
            return ptz.ReplyText("Hello, " + ptz.GetPushName() + "!")
        },
    })
}
```

### Command Fields

| Field | Type | Description |
|---|---|---|
| `Name` | `string` | Primary command name |
| `Aliases` | `[]string` | Alternative names |
| `Description` | `string` | Short description |
| `Usage` | `string` | Usage example |
| `Category` | `string` | Category shown in menu |
| `OwnerOnly` | `bool` | Restrict to owner |
| `GroupOnly` | `bool` | Group chat only |
| `AdminOnly` | `bool` | Group admin only |
| `BotAdmin` | `bool` | Bot must be admin |
| `Handler` | `func(*Ptz) error` | Command handler |

### Ptz Methods

```go
ptz.ReplyText(text)
ptz.ReplyImage(data, mime, caption)
ptz.ReplyVideo(data, mime, caption)
ptz.ReplyAudio(data, mime)
ptz.ReplySticker(data, mime, animated)
ptz.ReplyDocument(data, mime, filename, caption)
ptz.React(emoji)
ptz.Unreact()
ptz.IsOwner()
ptz.IsAdmin()
ptz.IsBotAdmin()
ptz.GetPushName()
ptz.Args        // []string
ptz.RawArgs     // string
ptz.IsGroup     // bool
ptz.Chat        // JID
ptz.Sender      // JID
ptz.Bot.Config
ptz.Bot.Client
```

## License

MIT
