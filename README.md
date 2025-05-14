# go-tcp-server

A lightweight Go-based TCP server that tracks transactions per second (TPS) and responds to every request with a 4-byte payload representing the integer **1800** (used as a simple health check).

---

## Features

* **Simple TCP listener**: Binds to a configurable host and port.
* **Fixed response**: Every incoming request gets back the byte sequence `[0x18, 0x00, 0x00, 0x00]` (decimal 1800).
* **Configurable via YAML**: Define `host`, `port`, and `statsInterval` in `config.yaml`.
* **Hourly stats reporting**: Logs total connections, TPS, and uptime at a configurable interval.
* **Graceful shutdown**: Handles `SIGINT`/`SIGTERM` to close listeners and stop cleanly.

---

## Prerequisites

* Go **1.20** or later
* `gopkg.in/yaml.v3` (automatically fetched via Go modules)

---

## Project Layout

```
go-tcp-server/
├── cmd/
│   └── tcpserver/
│       └── main.go          # Application entrypoint
├── config.yaml              # YAML config for host, port, statsInterval
├── internal/
│   ├── config/
│   │   └── config.go        # Reads & parses config.yaml
│   └── server/
│       └── server.go        # Core TCP server logic
├── go.mod                   # Module definition & dependencies
└── go.sum                   # Locked dependency checksums
```

---

## Configuration (`config.yaml`)

```yaml
host: "0.0.0.0"       # Interface to bind (e.g. 0.0.0.0 for all interfaces)
port: 8080              # TCP port to listen on
statsInterval: "1h"    # Go duration string for stats logging (e.g. "30m", "1h")
```

---

## Quick Start

1. **Install dependencies & tidy**

   ```bash
   go mod tidy
   ```

2. **Run the server**

   ```bash
   go run ./cmd/tcpserver
   ```

   The server reads `config.yaml` by default and starts listening on the configured host and port.

3. **Build a standalone binary** (optional)

   ```bash
   mkdir -p bin
   go build -o bin/tcpserver ./cmd/tcpserver
   ./bin/tcpserver
   ```

---

## Testing the Server

In a separate terminal, send a random payload and inspect the 4-byte response:

```bash
# Generate a random 16-character alphanumeric string, send to the server, and hex-dump the reply:
echo -n "$(head -c16 /dev/urandom | base64 | tr -dc 'A-Za-z0-9' | head -c16)" \
  | nc localhost 8080 \
  | xxd
```

**Expected output:**

```
00000000: 1800 0000                      ....
```

This confirms the TCP server is running and responding with the constant `1800` payload.

---

## License

MIT License. See [LICENSE](LICENSE) for details.
