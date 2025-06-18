# **A Master Prompt Sequence for Building GoChat**

**Author:** dunamismax
**Version:** 1.0.0
**Date:** 2025-06-17

---

## **Phase 1: Project Initialization & Configuration**

**Objective:** Atomically generate the complete project skeleton, dependency definitions, and toolchain configuration.

---

### **Prompt 1: Generate Directory Structure and Foundational Files**

**Directive:**
Generate the foundational structure for the `socket_chat` project. Create the standard Go project layout and the `.gitignore` and `LICENSE` files.

**1. Directory Structure:**

```sh
socket_chat/
├── .gitignore
├── LICENSE
├── cmd/
│   ├── client/
│   └── server/
├── internal/
│   └── server/
└── pkg/
```

**2. File: `.gitignore`**

```gitignore
# Go
#
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Go workspace file
go.work

# Environment variables
.env

# Build cache
/build/
```

**3. File: `LICENSE`**

```text
MIT License

Copyright (c) 2025 dunamismax

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

### **Prompt 2: Initialize Go Module and Linter Configuration**

**Directive:**
Initialize the Go module and create the linter configuration. These steps establish the project's identity and define its code quality standards.

**1. Initialize Go Module:**
Execute the following command from within the `socket_chat` directory:

```sh
go mod init github.com/dunamismax/golang/socket_chat
```

This will generate the `go.mod` file.

**2. Create Linter Configuration File: `.golangci.yml`**

```yaml
# Copyright (c) 2025-present dunamismax. All rights reserved.
#
# filename: .golangci.yml
# author: dunamismax
# version: 1.0.0
# date: 17-06-2025
# github: <https://github.com/dunamismax>
# description: Configuration for the golangci-lint linter.

run:
  timeout: 5m
  # The default concurrency (GOMAXPROCS) is used if not specified.
  # We can set it to 4 for consistency.
  concurrency: 4

linters-settings:
  govet:
    check-shadowing: true
  gocyclo:
    min-complexity: 15
  maligned:
    suggest-new: true
  goconst:
    min-len: 2
    min-occurrences: 3

linters:
  disable-all: true
  enable:
    - gofmt
    - goimports
    - revive
    - govet
    - staticcheck
    - unused
    - errcheck
    - goconst
    - gocyclo
    - ineffassign
    - typecheck
    - wastedassign

issues:
  exclude-rules:
    # Exclude complaining about `context.Context` from function signatures.
    - path: _test\.go
      linters:
        - gocyclo
        - errcheck
```

---

## **Phase 2: Core Concurrent Server Logic**

**Objective:** Implement the server's core components, focusing on concurrent, channel-based communication.

---

### **Prompt 3: Generate the Server-Side Client Representation**

**Directive:**
Create the file `internal/server/client.go`. This module defines the `Client` struct, which encapsulates a single user's connection. It includes the concurrent `readPump` and `writePump` goroutines that are fundamental to preventing race conditions and blocking I/O on the server. Adhere to all file header conventions.

---

### **Prompt 4: Generate the Concurrent Hub**

**Directive:**
Create the file `internal/server/hub.go`. This module defines the central `Hub`, the heart of the chat server's concurrency model. The hub runs in a single goroutine and uses channels to safely manage client registration, unregistration, and message broadcasting, eliminating the need for mutexes.

---

### **Prompt 5: Generate the Main Server**

**Directive:**
Create the file `internal/server/server.go`. This module defines the primary `Server` struct. Its responsibility is to listen for incoming TCP connections, create a `Client` instance for each connection, and register it with the `Hub`. This abstracts away the network listener logic from the hub's state management logic.

---

## **Phase 3: Application Entrypoints**

**Objective:** Construct the executable `main` packages for both the server and the client.

---

### **Prompt 6: Generate the Server Entrypoint**

**Directive:**
Create the file `cmd/server/main.go`. This is the entrypoint for the chat server application. It initializes structured logging (`slog`), sets up graceful shutdown to handle `SIGINT` and `SIGTERM`, and starts the server.

---

### **Prompt 7: Generate the Client Entrypoint**

**Directive:**
Create the file `cmd/client/main.go`. This is the entrypoint for the command-line client. The application connects to the server and launches two goroutines: one for reading messages from the server and one for reading user input from `stdin` and sending it to the server. This ensures the UI remains responsive while waiting for network I/O.

---

## **Phase 4: Final Verification**

**Objective:** Ensure the generated codebase complies with all Go standards for formatting, dependency management, and quality.

---

### **Prompt 8: Execute Quality Assurance Protocol**

**Directive:**
The final action is to verify project integrity. The following commands must be executed from the `socket_chat` project root. The project is complete only when all commands pass without error.

**Execution Commands:**

```sh
# 1. Tidy dependencies. This ensures go.mod and go.sum are accurate.
go mod tidy

# 2. Format all Go source files. Code that is not gofmt-ed is considered broken.
gofmt -w .

# 3. Run the linter to check for style issues and potential bugs.
# Ensure golangci-lint is installed: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run

# 4. Build both binaries to ensure compilation.
go build -o build/server ./cmd/server
go build -o build/client ./cmd/client

# 5. [Optional] Run the server to confirm it starts without crashing.
./build/server
```

This sequence is now complete and optimized for Go. Awaiting next directive.
