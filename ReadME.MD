# Go SSE EventStream Live Demos

This project provides experimental demonstrations of real-time events using **Server-Sent Events (SSE)** in Go. It features two different implementations to showcase the versatility of streaming data to the browser over a single HTTP connection.

## Features

- **Real-time Server-to-Client Communication**: Utilizes Server-Sent Events (SSE) to push data from a Go backend to a web client.
- **Graceful Shutdown**: The Go HTTP server intercepts interrupt signals (like `Ctrl+C`) and performs a graceful shutdown.
- **Two Unique Demos**:
  - 🔄 **Original Auto-Play Stream**: A classic raw-text stream that types out predefined "Alert" and "Panic" messages.
  - 🤖 **Interactive AI Assistant Search**: A dynamic JSON-based stream where the server "generates" live responses based on your search queries.
- **Modern UI Design**: Clean, responsive frontend interfaces equipped with light/dark theme toggles.

## How It Works

### Server-Side (Go)

The server exposes multiple endpoints:

- `/` & `/events`: Serves the original `typer.html` file and handles the raw text stream. It uses `http.Flusher` to push chunks of text with a 500ms delay.
- `/search-app` & `/search-events`: Serves the new `search.html` file and handles interactive search queries. It parses the `?q=` parameter, simulates a generated response, and streams properly structured JSON objects (e.g., `data: {"word": "..."}`).

### Client-Side (HTML/JS)

- **Connection Management**: An `EventSource` object is instantiated to establish a persistent connection. In the interactive demo, older connections are manually closed via `sse.close()` before starting a new search to prevent overlapping streams.
- **Live Rendering**: The `onmessage` event listener captures each incoming data chunk (either raw text or parsed JSON) and appends it directly before the blinking cursor in the DOM, creating a live typing effect.
- **Error Handling**: The `onerror` event listener is implemented to gracefully catch disconnects and prevent the browser's default behavior of an infinite automatic reconnection loop once the stream is explicitly finished.

## Prerequisites

- [Go](https://golang.org/dl/) 1.23 or higher

## Getting Started

1. **Clone or navigate to the repository:**

   ```bash
   cd go-sse-eventstream-live
   ```

2. **Run the Go application:**

   ```bash
   go run main.go
   ```

   *Optional: You can specify a custom port using the `-addr` flag (default is `:8080`).*

   ```bash
   go run main.go -addr=":9000"
   ```

3. **Open the web applications:**

   Open your browser and navigate to the following URLs:
   - 👉 **[http://localhost:8080](http://localhost:8080)** (Original Auto-Play Demo)
   - 👉 **[http://localhost:8080/search-app](http://localhost:8080/search-app)** (Interactive AI Search Demo)

## What is SSE (Server-Sent Events)?

Server-Sent Events (SSE) is a technology that allows real-time communication from the server to the client via a persistent, unidirectional connection.

Unlike WebSockets (which are bidirectional) or traditional HTTP polling (where the client repeatedly asks for data), SSE allows the server to efficiently push messages to the client as long as the connection remains open.

### Why use SSE?

- It's a standard HTTP protocol.
- Built-in reconnection mechanism (though explicitly disabled in these demos to show stream completion).
- Lightweight and incredibly simple to implement for one-way (server-to-client) real-time data flows like live logs, status updates, or AI text generation.
