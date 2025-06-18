# Project Ideas

## **Category 1: Foundational CLI & Networking Tools**

These projects solidify core Go skills in concurrency, I/O, and data manipulation. They are direct parallels to the `socket_chat` project's philosophy.

1. **Concurrent Web Scraper / Crawler**

   - **Concept:** A command-line tool that takes a starting URL and recursively scrapes all links on the same domain up to a certain depth. It should process multiple pages concurrently to maximize speed.
   - **Learning Value:**
     - **Concurrency:** Heavily utilize goroutines and channels to manage a pool of workers fetching URLs. This prevents the application from blocking on network I/O.
     - **Networking:** Use the `net/http` package to make robust HTTP requests.
     - **Data Parsing:** Use the `goquery` library to parse HTML and extract links, which is a practical application of external libraries.
     - **State Management:** Use maps and channels to keep track of visited URLs to avoid redundant fetches and getting stuck in loops.

2. **`gogrep`: A Concurrent `grep` Clone**

   - **Concept:** A CLI tool that searches for a pattern in files. It should be able to search multiple files concurrently and be significantly faster than the standard `grep` for a large number of files.
   - **Learning Value:**
     - **File I/O:** Master concurrent file reading using `os` and `bufio` within goroutines.
     - **Concurrency Control:** Use a `sync.WaitGroup` to manage the lifecycle of file-reading goroutines and a channel to aggregate results.
     - **CLI Arguments:** Use the `flag` package to parse command-line flags for patterns, file paths, and options (e.g., case-insensitivity).

3. **Simple Reverse Proxy / Load Balancer**
   - **Concept:** A server that listens on a port and forwards incoming HTTP requests to one or more backend servers in a round-robin fashion.
   - **Learning Value:**
     - **Advanced `net/http`:** Go beyond a basic server and dive into the `net/http/httputil` package, specifically `NewSingleHostReverseProxy`.
     - **Concurrency:** Each incoming request is inherently handled in its own goroutine by the Go HTTP server, providing a natural extension of the `socket_chat` model.
     - **State Management:** Maintain a thread-safe list of backend servers and a counter for the round-robin logic.

## **Category 2: Intermediate Web & Distributed Systems**

These projects introduce database interaction and more complex service-oriented architectures, adhering to the established "ORM-free" and "standard-library-first" mandates.

1. **URL Shortener Service**

   - **Concept:** A web service with a RESTful API that accepts a long URL and returns a short, unique URL. When a user visits the short URL, the service redirects them to the original long URL.
   - **Learning Value:**
     - **Web Service:** Build a clean API using `net/http` and the `chi` router for routing.
     - **Database Interaction:** Use the `database/sql` package with `pgx` (for PostgreSQL) or `go-sqlite3` (for SQLite) to store the URL mappings. This reinforces the "raw SQL" and "Repository Pattern" laws.
     - **API Design:** Practice designing clear and simple RESTful endpoints (`POST /shorten`, `GET /{short_code}`).

2. **Distributed Key-Value Store**

   - **Concept:** A simplified, in-memory key-value database that can be run on multiple nodes. A write to one node should replicate to others in the cluster.
   - **Learning Value:**
     - **Distributed Systems:** Tackle the fundamental challenges of state replication and node communication.
     - **Concurrency:** Use concurrent maps (`sync.Map`) for safe in-memory storage.
     - **Networking:** Design a simple TCP or HTTP-based protocol for nodes to communicate write operations to each other.

3. **Real-Time Analytics Engine**
   - **Concept:** A service that ingests streaming data (e.g., website clicks, application events) and provides a real-time view of aggregated metrics via a simple API or WebSocket connection.
   - **Learning Value:**
     - **Data Pipelines:** Use channels to create a pipeline: one goroutine ingests data, another aggregates it (e.g., counts events per second), and another serves it.
     - **Time-Series Data:** Work with time-based data and learn techniques for windowing and aggregation.
     - **WebSockets:** For a more advanced version, push live updates to a web dashboard using a WebSocket library, which complements the TCP knowledge from `socket_chat`.

## **Category 3: Advanced Systems Programming**

These projects require a deeper understanding of Go and the underlying operating system, directly aligning with the "Build Your Own X" philosophy.

1. **`gogit`: Build Your Own Git**

   - **Concept:** Re-implement a small subset of Git commands from scratch. Focus on understanding the core objects: blobs, trees, and commits.
   - **Learning Value:**
     - **File System Mastery:** Deeply interact with the file system to create `.git` directories and objects.
     - **Hashing & a-:** Use the `crypto/sha1` package to hash file contents and create Git's content-addressed storage.
     - **Data Structures:** Design structs to represent Git's internal objects and understand how they link together.

2. **`gocontainer`: Build Your Own Docker**
   - **Concept:** Use Linux namespaces and control groups (cgroups) to build a very basic container runtime that can isolate a process's filesystem and PID.
   - **Learning Value:**
     - **Low-Level OS Interaction:** Use the `syscall` package to interact directly with Linux kernel features. This is a significant step up in systems programming.
     - **Process Management:** Use `os/exec` to create and manage child processes that will run inside the "container."
     - **Root Filesystems:** Learn how to use `chroot` to give a process its own isolated root directory.

Select a project that aligns with your current learning objectives. All of these will provide a substantial and rewarding challenge.
