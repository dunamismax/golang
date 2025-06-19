<h1 align="center">socket_chat</h1>

<p align="center">
  socket_chat is a lightweight, concurrent, and scalable TCP-based chat server and client, built with pure Go.
  <br />
  It demonstrates idiomatic, concurrent Go using goroutines and channels, without external frameworks.
</p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Language-Go-blue.svg" alt="Go"></a>
  <a href="https://github.com/dunamismax/golang/blob/main/LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
  <a href="https://github.com/dunamismax/golang/pulls"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square" alt="PRs Welcome"></a>
  <a href="https://github.com/dunamismax/golang/stargazers"><img src="https://img.shields.io/github/stars/dunamismax/golang?style=social" alt="GitHub Stars"></a>
</p>

---

## ✨ Guiding Philosophy

socket_chat is built on a few core principles:

- **Simplicity Over Complexity**: We use Go's standard library (`net`, `log/slog`) as the foundation. The architecture avoids unnecessary abstractions and third-party dependencies, focusing on clear, readable code. The Go proverb "Clear is better than clever" is our guiding light.
- **Concurrency by Design**: The application is architected around Go's native concurrency primitives. A central "hub" uses channels for safe, concurrent state management, and each client connection is handled in its own goroutine to maximize I/O throughput.
- **Idiomatic Go**: The application is designed to feel natural to Go developers. We leverage Go's strengths—simplicity, concurrency, and powerful tooling—to create a clean and productive development experience.
- **Structured & Observable**: From the ground up, the application uses structured JSON logging via the `slog` package, ensuring that the server and client are observable and production-ready.

---

## 🚀 Getting Started

You will need the Go toolchain installed (version 1.22+ recommended). All commands must be run from the `socket_chat` directory, which is the root of the Go module.

First, navigate to the project's root directory:

```sh
cd socket_chat
```

### 1. Run the Server

Open a terminal in the `socket_chat` directory and run the server. It will listen for connections on port `8080`.

```sh
# From the project root (socket_chat/)
go run ./cmd/server/main.go
```

The server will log that it has started.

```json
{
  "time": "2025-06-17T...",
  "level": "INFO",
  "msg": "Starting chat server",
  "address": "localhost:8080"
}
```

### 2. Run the Client(s)

Open one or more new terminal windows. **In each new terminal, navigate to the `socket_chat` directory first.**

```sh
# Ensure you are in the correct directory
cd socket_chat

# Run the client
go run ./cmd/client/main.go
```

First, the client will prompt you for a nickname. After you enter one, you can begin sending messages. Messages sent from one client will be broadcast to all other connected clients.

**Terminal 1 (Client A):**

```sh
$ go run ./cmd/client/main.go
Enter your nickname: Alice
You can now start sending messages.
Hello everyone!
[Bob]: Hi Alice!
```

**Terminal 2 (Client B):**

```sh
$ go run ./cmd/client/main.go
Enter your nickname: Bob
You can now start sending messages.
[Alice]: Hello everyone!
Hi Alice!
```

---

## 🏗️ Project Structure

socket_chat is organized using the Standard Go Project Layout for clarity and maintainability.

```sh
socket_chat/
├── cmd/
│   ├── server/             # Entry point for the chat server.
│   └── client/             # Entry point for the command-line client.
├── internal/
│   └── server/             # Core application code: hub, client mgmt, server logic.
├── .golangci.yml           # Linter configuration.
├── go.mod                  # Go module definition.
└── go.sum                  # Dependency checksums.
```

---

## 🤝 Contribute

**This project is built by the community, for the community. We need your help!**

Whether you're a seasoned Go developer or just starting, there are many ways to contribute:

- **Report Bugs:** Find something broken? [Open an issue](https://github.com/dunamismax/golang/issues) and provide as much detail as possible.
- **Suggest Features:** Have a great idea for a new feature (e.g., private messages, multiple rooms)? [Start a discussion](https://github.com/dunamismax/golang/discussions) or open a feature request issue.
- **Write Code:** Grab an open issue, fix a bug, or implement a new system. [Submit a Pull Request](https://github.com/dunamismax/golang/pulls) and we'll review it together.
- **Improve Documentation:** Great documentation is as important as great code. Help us make our guides and examples clearer and more comprehensive.

If this project excites you, please **give it a star!** ⭐ It helps us gain visibility and attract more talented contributors like you.

### Connect

Connect with the author, **dunamismax**, on:

- **Twitter:** [@dunamismax](https://twitter.com/dunamismax)
- **Bluesky:** [@dunamismax.bsky.social](https://bsky.app/profile/dunamismax.bsky.social)
- **Reddit:** [u/dunamismax](https://www.reddit.com/user/dunamismax)
- **Discord:** `dunamismax`
- **Signal:** `dunamismax.66`

## 📜 License

This project is licensed under the **MIT License**. See the `LICENSE` file for details.
