# sptzxgo

sptzxgo is a modular and scalable WhatsApp bot built in Golang. It leverages the powerful `whatsmeow` library for WhatsApp Web integration and provides a flexible, event-driven architecture for handling commands and automating tasks efficiently.

## Features

### Core Features
- **WhatsApp Web Integration**:
  - Designed to interact with both group and personal chats.
  - Full duplex messaging and event handling.

- **Command Management System**:
  - Modular commands organized into categories (e.g., AI tools, games, group utilities).
  - Dynamic command registration and easy extensibility.

- **Database Session Management**:
  - Utilizes `SQLite3` for fast, lightweight storage of user sessions.

- **Environment-Driven Configuration**:
  - Easily configurable via `.env` files for dynamic runtime settings.

### Additional Functionality
- **Games and Trivia**:
  - Built-in mini-games like trivia, quizzes, and more.

- **Custom Anti-Spam Features**:
  - Rate limits (max messages per second / minute).
  - Ban users exceeding limits (duration customizable).

- **Advanced Bot Customization**
  - Timezone-aware settings.
  - Dynamic sticker production (custom sticker name/author).

## Prerequisites

To use this bot effectively, make sure you meet the requirements below:

- **Golang Version:** 1.21 or later
- **Database:** SQLite version 3
- **Other Utilities:**
  - A valid `.env` configuration file (see setup section).

## Installation and Setup

### Clone the Repository

```bash
git clone https://github.com/siputzx/sptzxgo.git
cd sptzxgo
```

### Install Go Dependencies

Run the following command to install dependencies:

```bash
go mod tidy
```

### Setup Environment Configuration

1. Copy the example `.env` file:
   ```bash
   cp .env.example .env
   ```

2. Open `.env` and customize the following fields:

| Variable           | Description                              | Required |
|--------------------|------------------------------------------|----------|
| `BOT_OWNERS`       | Comma-separated WhatsApp owners' numbers| Yes      |
| `BOT_PREFIX`       | Define available bot prefixes           | No       |
| `SESSION_DB`       | Path to your SQLite database            | No       |
| `LOGIN_METHOD`     | Either `qr` or `paircode`               | No       |
| `PAIRING_PHONE`    | Phone number for pairing                | Necessary if `paircode` |

Save the file after configuration changes.

### Run the Bot

Once everything is set up, execute:

```bash
go run main.go
```

## Project Structure

```plaintext
sptzxgo/
├── main.go (application entry point)
├── .env.example (runtime system environment definitions)
├── src/
│   ├── commands        (modular commands: AI, games, groups, etc.)
│   ├── core            (central bot systems, e.g. context storage)
│   ├── handler         (managing message + status events)
│   ├── middleware      (anti-spam)
│   ├── serialize       (helpers for WhatsApp message objects)
│   └── config          (environment + runtime loader)
```

## License

This project is licensed under the **Apache 2.0 License**. See the `LICENSE` file for full details.

---

### Contributing

Contributions, issues, and feature requests are welcome! Feel free to fork this repository and submit your improvements as pull requests.
