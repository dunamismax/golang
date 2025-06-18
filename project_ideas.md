# Project Ideas (backend / networking focus)

## Category 1: Real-Time Communication & Streaming Systems

3. **gRPC Bidirectional Streaming Service**
   - **Concept:** A service where both the client and server can send a stream of messages to each other over a single, long-lived connection. A good example would be a "price checker" where the client sends a stream of product IDs and the server streams back price updates for those products as they happen.
   - **Why it's a great project:** It introduces you to gRPC, a high-performance RPC framework that is a cornerstone of modern microservices architecture. It's a significant step beyond REST/JSON.
   - **Core Go Concepts:**
     - **Protocol Buffers (Protobuf):** Defining your API contract in `.proto` files, which is a language-agnostic way to define data structures and services.
     - **gRPC:** Implementing both server and client stubs from the generated Protobuf code.
     - **Advanced Concurrency:** Managing the dual streams, handling errors, and using context for cancellation and timeouts are more complex than with a simple request-response.

## Category 2: Distributed Systems & State Management

These projects force you to solve problems of state, consensus, and communication between multiple server nodes.

3. **Peer-to-Peer (P2P) File Sharing Application**
   - **Concept:** A command-line application where a user can "share" a file. Other instances of the application on the network can then discover and download parts of that file from multiple peers simultaneously, similar to BitTorrent.
   - **Why it's a great project:** It combines low-level networking, concurrency, and distributed discovery into a single, challenging application.
   - **Core Go Concepts:**
     - **Custom Protocol:** Design a protocol for peers to request file chunks and send them.
     - **Peer Discovery:** Use UDP multicasting (via `net.ListenMulticastUDP`) for peers on the same local network to find each other without a central server.
     - **Concurrent I/O:** Simultaneously download different chunks of a file from multiple peers and write them to the correct offset in the destination file on disk.

## Category 3: Re-implementing Foundational Technologies

These projects provide the deepest learning by forcing you to build simplified versions of tools you use every day.

1. **A Simple Reverse Proxy (`mini-nginx`)**

   - **Concept:** An HTTP server that receives a request and forwards it to a configured backend service. It should be able to handle basic load balancing (e.g., round-robin) between multiple backends.
   - **Why it's a great project:** Reverse proxies are a fundamental building block of web architecture. Building one teaches you the intricacies of the HTTP protocol and request/response manipulation.
   - **Core Go Concepts:**
     - **`net/http/httputil`:** The `ReverseProxy` struct is the heart of this project. You'll learn to customize its `Director` and `Transport` to modify requests and handle responses.
     - **Concurrency:** The Go `net/http` server handles each request in a goroutine automatically, making this concurrently scalable by default.
     - **Configuration Management:** Use Viper to load backend server configurations from a file or environment variables.

2. **A Network Packet Analyzer (`gopcap`)**
   - **Concept:** A tool that listens on a network interface and decodes network packets, printing a summary of the protocol (e.g., "TCP packet from IP:Port to IP:Port").
   - **Why it's a great project:** It takes you to the lowest level of the network stack you can practically access from user space. It demystifies what is actually happening "on the wire."
   - **Core Go Concepts:**
     - **Raw Sockets / Packet Capture:** Use a library like `gopacket` which provides Go bindings for `libpcap` to capture raw network data.
     - **Protocol Decoding:** Write parsers for Ethernet, IP, TCP, and UDP headers to extract meaningful information from the raw byte slices.
     - **CLI Design:** Present the captured data in a clear, useful format, similar to `tcpdump`.
