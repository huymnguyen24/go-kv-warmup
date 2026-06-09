# go-kv-warmup

A small, concurrent-safe, in-memory key/value store exposed over HTTP — built with the Go standard library only.

This is a **Phase 1 Go warmup project**: the smallest program that exercises every idiom you need before writing real Go services (Kubernetes controllers, in particular). It touches HTTP routing, a mutex-guarded shared map, a background goroutine, and graceful shutdown via signal-aware context. See [phase-1-go-warmup-guided.md](phase-1-go-warmup-guided.md) for the full guided exercise this code was written against.

## Features

- **In-memory KV store** — string keys → string values, guarded by a `sync.RWMutex` for safe concurrent access.
- **REST-ish HTTP API** — `PUT` / `GET` / `DELETE` a key, and `GET` the list of all keys.
- **Background reporter** — a goroutine logs the current key count every 10 seconds.
- **Graceful shutdown** — `Ctrl+C` (SIGINT) or SIGTERM drains in-flight requests with a 5-second timeout before exiting.
- **Standard library only** — no third-party router or dependencies.

## Requirements

- Go 1.23.4 or newer

## Running

```sh
go run .
```

The server listens on `:8080`:

```
listening on :8080
store has 0 keys
```

Press `Ctrl+C` to shut down cleanly.

## API

The store lives at `/kv`. Keys are taken from the URL path.

| Method   | Path         | Body        | Description                          | Responses |
|----------|--------------|-------------|--------------------------------------|-----------|
| `PUT`    | `/kv/{key}`  | value (raw) | Create or update a key               | `201 Created` (new key), `200 OK` (updated) |
| `GET`    | `/kv/{key}`  | —           | Read a value                         | `200 OK` + value, `404 Not Found` |
| `DELETE` | `/kv/{key}`  | —           | Delete a key                         | `204 No Content`, `404 Not Found` |
| `GET`    | `/kv`        | —           | List all keys (JSON array)           | `200 OK` + `["a","b",...]` |

The request body for `PUT` is stored verbatim as the value (any bytes). List responses are JSON; value responses are the raw stored bytes.

### Examples

```sh
# Create a key
curl -i -X PUT --data 'world' http://localhost:8080/kv/hello
# -> 201 Created

# Update the same key
curl -i -X PUT --data 'there' http://localhost:8080/kv/hello
# -> 200 OK

# Read it back
curl http://localhost:8080/kv/hello
# -> there

# List all keys
curl http://localhost:8080/kv
# -> ["hello"]

# Delete it
curl -i -X DELETE http://localhost:8080/kv/hello
# -> 204 No Content

# Reading a missing key
curl -i http://localhost:8080/kv/hello
# -> 404 Not Found
```

## Project layout

| File          | Responsibility                                              |
|---------------|-------------------------------------------------------------|
| `main.go`     | Server lifecycle: routing, background reporter, graceful shutdown |
| `store.go`    | `Store` — the mutex-guarded map and its `Get`/`Set`/`Delete`/`Keys`/`Count` methods |
| `handlers.go` | `API` — HTTP handlers that translate requests into store calls |

## License

[MIT](LICENSE) © 2026 Huy Nguyen
