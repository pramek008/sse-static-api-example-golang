# High-Efficiency Go Streaming API (LLM Simulation)

A lightweight, high-performance REST API written in **Go (Golang)** designed to simulate Large Language Model (LLM) streaming responses.

This project is optimized for **low resource consumption** and **high concurrency**. It uses the Go standard library (`net/http`) to minimize dependencies and binary size, making it ideal for load testing, benchmarking, or acting as a mock server for frontend development without the overhead of a real LLM.

## 🚀 Key Features

- **⚡ Ultra-Lightweight:** Docker image size is **< 15MB** (via Multi-stage build).
- **🛡️ Resource Efficient:** Configured to run smoothly on **0.1 CPU** and **20MB RAM**.
- **🧵 High Concurrency:** Uses Go Goroutines to handle thousands of simultaneous connections with minimal overhead.
- **📡 Streaming Support:**
  - **SSE (Server-Sent Events):** Standard web streaming (simulating ChatGPT/Claude).
  - **NDJSON (Newline Delimited JSON):** Line-by-line JSON streaming (simulating Ollama).
- **🛠️ Zero Dependencies:** Built using only Go's standard library. No frameworks (Gin/Fiber) needed.

---

## 🛠️ Prerequisites

- Docker
- Docker Compose

---

## 🏃 Quick Start

1.  **Clone or Create Files**
    Ensure you have `main.go`, `Dockerfile`, and `docker-compose.yml` in the same directory.

2.  **Run with Docker Compose**
    This command builds the tiny Alpine image and starts the container with strict resource limits.

    ```bash
    docker compose up --build
    ```

3.  **Verify**
    The server will start at `http://localhost:3000`.

---

## 🔗 API Endpoints

### 1. SSE Stream (LLM Simulation)

Simulates a token-by-token text generation using Server-Sent Events.

- **URL:** `GET /stream-sse`
- **Format:** `text/event-stream`
- **Test:**
  ```bash
  curl -N http://localhost:3000/stream-sse
  ```

### 2. NDJSON Stream (Ollama Style)

Simulates streaming using Newline Delimited JSON.

- **URL:** `GET /stream-ndjson`
- **Format:** `application/x-ndjson`
- **Test:**
  ```bash
  curl -N http://localhost:3000/stream-ndjson
  ```

### 3. Infinite Loop Stream

Streams text endlessly. Useful for testing connection stability or client-side memory handling.

- **URL:** `GET /stream-loop`
- **Format:** `text/event-stream`

### 4. Static JSON

Returns a standard, complete JSON response (non-streaming).

- **URL:** `GET /api/data`
- **Format:** `application/json`

### 5. Health Check

- **URL:** `GET /health`

---

## 🏗️ Architecture & Optimizations

Why is this version more efficient than Node.js?

### 1. Memory Management

- **Global Pre-processing:** The text data is split into arrays only once during application startup (`init()` function).
- **Zero Allocations per Request:** Unlike the Node.js version which might split strings on every request, this Go version reuses the same memory pointers for the text data, significantly reducing Garbage Collection (GC) pressure.

### 2. Concurrency Model

- **Goroutines:** Each incoming request spawns a Goroutine, which costs only ~2KB of stack memory. This allows the server to handle high-throughput scenarios without blocking the main thread.
- **Context Awareness:** If a client disconnects, `r.Context().Done()` detects it immediately, killing the Goroutine and freeing resources instantly.

### 3. Docker Optimization

- **Multi-Stage Build:** We compile the binary in a heavy `golang` image, then copy _only_ the resulting binary to an empty `alpine` image.
- **Hard Limits:** The `docker-compose.yml` strictly enforces:
  ```yaml
  limits:
    cpus: "0.1" # 10% of a single core
    memory: 20M # Max 20MB RAM
  ```

---

## 🧪 Load Testing

Since this service is highly optimized, you can stress test it using tools like **k6**, **wrk**, or **Apache Bench** to see how well it handles traffic under the 20MB RAM limit.

**Example using `wrk`:**

```bash
# 12 threads, 400 connections, for 30 seconds
wrk -t12 -c400 -d30s http://localhost:3000/api/data
```
