# mestralive-go

Go SDK for **fast, in-process pub/sub** on [mestralive](https://github.com/sahi7/mestralive).

Use this when your **Go service** (chat backend, bridge, worker) needs to publish **opaque binary payloads** (TLV, AWP, protobuf, raw bytes) without paying for JSON on every message.

```text
Internet users  →  mestralive-live  (JSON WebSocket / HTTP)   ← browsers & partners
Your Go service →  mestralive-go    (in-process Publish)      ← this SDK
```

**Module:** [`github.com/mestracorp/mestralive-go`](https://github.com/mestracorp/mestralive-go)

---

## What this SDK does

| Does | Does not |
|------|----------|
| Embed a mestralive **runtime bus** in your process | Replace mestralive-live for browser clients |
| `Publish(topic, typ, []byte)` with **opaque** payloads | Call `POST /app/v1/publish` or speak WS JSON |
| `Subscribe` / `Accept` for in-process fanout targets | Dial `:9000` on localhost (v1) |
| Authenticate the **service once** at `Open` | Re-auth with a JSON envelope per message |
| Report `Accepted` / `Rejected` / `Disconnected` | Promise that a remote device received the bytes |

`Accepted` means the runtime accepted the message into an outbound/mailbox buffer — **not** “the peer got a TCP packet.”

---

## Install

```bash
go get github.com/mestracorp/mestralive-go@latest
```

Until a tagged mestralive release that includes `pkg/fanout` is what your `go.mod` resolves to, local development may need:

```go
// go.mod
replace github.com/mestralive/mestralive => ../mestralive
```

(Point `../mestralive` at a checkout that contains `pkg/fanout`, e.g. the `wip` branch.)

---

## Quick start

```go
package main

import (
	"log"
	"os"

	mestralive "github.com/mestracorp/mestralive-go"
)

func main() {
	bus, err := mestralive.Open(mestralive.Config{
		ServiceToken: os.Getenv("MESTRALIVE_SERVICE_TOKEN"),
		Owners:       1, // fine for a single-process service
	})
	if err != nil {
		log.Fatal(err) // missing/invalid token → ErrUnauthorized
	}
	defer bus.Stop()

	if err := bus.Start(); err != nil {
		log.Fatal(err)
	}

	// Create in-process subscriber endpoints (synthetic connections).
	alice, err := bus.Accept()
	if err != nil {
		log.Fatal(err)
	}
	bob, err := bus.Accept()
	if err != nil {
		log.Fatal(err)
	}
	_ = bus.Subscribe(alice, "chat.room1")
	_ = bus.Subscribe(bob, "chat.room1")

	// Your codec owns the bytes — mestralive does not JSON-parse them.
	tlv := []byte{0x01, 0x02, 0xff}
	res, err := bus.Publish("chat.room1", 1 /* app message type */, tlv)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("subscribers=%d accepted=%d rejected=%d",
		res.Subscribers, res.Accepted, res.Rejected)
}
```

Runnable copy: [`examples/inprocess-tlv`](examples/inprocess-tlv).

```bash
export MESTRALIVE_SERVICE_TOKEN=dev-token
go run ./examples/inprocess-tlv
```

---

## Authentication

### Service (this SDK)

The **process** proves it may open a bus. That happens **once** at `Open`, not on every `Publish`.

| Step | Action |
|------|--------|
| 1 | Create a service secret (token) and store it in your secret manager |
| 2 | Set `MESTRALIVE_SERVICE_TOKEN` in the environment, **or** pass `Config.ServiceToken` |
| 3 | Optional: set `Config.ServiceTokenAllowlist` so only listed tokens are accepted |
| 4 | Call `Open` / `OpenFromEnv` — empty or wrong token → `ErrUnauthorized` |

```go
// From environment
bus, err := mestralive.OpenFromEnv()

// Or explicit
bus, err := mestralive.Open(mestralive.Config{
	ServiceToken:          "...",
	ServiceTokenAllowlist: []string{"..."}, // optional; recommended in production
})
```

**Rotation (v1):** put the new secret in the environment / secret store and **restart** the process.

Never put the service token in a browser, mobile app, or public repo.

### End users (not this SDK)

Human/partner clients authenticate to **mestralive-live** (Bearer `key_id.secret` or identity tokens) over TLS + JSON WebSocket/HTTP.

Your service should:

1. Trust the session only after your edge / live auth succeeds  
2. Enforce product rules (e.g. “may this user publish to `chat.room1`?”)  
3. Then call `Publish` with TLV bytes  

See the Agyraa-shaped spike: capability check **before** Publish —  
[`agyraa/spikes/mestralive-bridge`](https://github.com/sahi7/agyraa) (local checkout path may vary).

---

## API overview

| Method | Purpose |
|--------|---------|
| `Open` / `OpenFromEnv` | Create a bus (not started); validates service token |
| `Start` / `Stop` | Start/stop the embedded runtime (`Stop` is idempotent) |
| `Accept` | Admit an in-process connection id (subscriber endpoint) |
| `Subscribe` / `Unsubscribe` | Topic membership for a connection |
| `Publish` | Fan out opaque `[]byte` to current subscribers |
| `OpenDial` | **Unsupported in v1** — returns an error (no `:9000` dial) |

### `Publish` result

| Field | Meaning |
|-------|---------|
| `Subscribers` | Matched before fanout budget |
| `Attempted` | Tried this call |
| `Accepted` | Entered outbound path / mailbox (**not** peer TCP receipt) |
| `Rejected` | Queue full / policy — apply backpressure |
| `Disconnected` | Recipient gone / stale |
| `Deferred` | Skipped by fanout budget |

If buffers fill and nothing drains, later publishes may show `Rejected` or `Disconnected`. That is **intentional** bounded-memory behavior.

### Topics

Allowed characters: letters, digits, `.` `_` `-` `/`  
Case-sensitive. Max length 256 bytes. Invalid → `ErrInvalidTopic`.

### Config knobs (common)

| Field | Default / notes |
|-------|-----------------|
| `ServiceToken` | **Required** |
| `Owners` | Runtime owner threads (default 2) |
| `MaxConnections` | Default 1024 |
| `MaxMessageBytes` | 0 = runtime default; oversize → `ErrPayloadTooLarge` |
| `ListenAddress` | **Empty = in-process only (recommended)** |
| `AllowPublicListen` | Required for non-loopback listen; **not** an internet client plane |

---

## Testing

From this repo:

```bash
# Unit + integration
go test ./... -count=1

# Data race detector
go test ./... -race -count=1

# Example
MESTRALIVE_SERVICE_TOKEN=dev-token go run ./examples/inprocess-tlv

# Optional microbenches (lab only — not capacity SLOs)
go test -run='^$' -bench=BenchmarkSDK_ -benchmem -count=3 -benchtime=200ms
```

What the tests cover:

- Open without / with wrong token  
- Public listen denied by default  
- Opaque binary publish (including `0x00` bytes)  
- Multi-topic isolation, concurrent publishers, stop under load  
- Oversized payload and connection admission limits  

---

## Architecture (why it is fast)

```text
your service
    mestralive.Open / Start
         │
         ▼
    pkg/fanout  ──►  internal/runtime Router.Publish
         │                │
         │                ├─ one shared payload copy for mailbox fanout
         │                └─ 24-byte big-endian frames into connection buffers
         ▼
    Result{Accepted, Rejected, ...}
```

- **No** `json.Marshal` on the Publish path  
- **No** HTTP/WebSocket hop for this SDK  
- **No** dial to `:9000` in v1 (`OpenDial` fails closed)  

Browser JSON-WS stays on **mestralive-live**. Mixing those planes is a common mistake: use live for users, this SDK for server-side speed.

---

## Errors you will see

| Error | Typical cause |
|-------|----------------|
| `ErrUnauthorized` | Missing or allowlist-rejected service token |
| `ErrUnsafeListen` | Non-loopback `ListenAddress` without `AllowPublicListen` |
| `ErrNotStarted` / `ErrClosed` | Call before `Start` or after `Stop` |
| `ErrPayloadTooLarge` | Payload exceeds `MaxMessageBytes` |
| `ErrInvalidTopic` | Empty or illegal topic name |
| `ErrDialUnsupported` | `OpenDial` in v1 |

---

## Versioning & support

- **Go:** see `go.mod`  
- **License:** [MIT](LICENSE)  
- **Issues / source:** https://github.com/mestracorp/mestralive-go  
- Deeper engine contracts live in the mestralive repo under `docs/fanout/` (optional reading for platform engineers).

---

## FAQ

**Can my React/PWA use this SDK?**  
No. Use mestralive-live (JSON WebSocket). This module is for **Go services**.

**Is `Accepted` the same as “delivered to the phone”?**  
No. It means the runtime accepted outbound work. End-to-end delivery needs your product ACK/protocol.

**Should I dial the runtime TCP port for isolation?**  
Not with this SDK in v1. Prefer in-process embed. Localhost TCP is often much slower and easy to misconfigure (echo lab vs real bus).

**Where do I put user auth?**  
On your edge / mestralive-live. Then authorize topics in your service, then `Publish`.
