# mestralive-go

Go SDK for **in-process TLV fanout** on the mestralive certified runtime.

| Plane | Use |
|-------|-----|
| **mestralive-live** (`/app/v1` JSON-WS) | Internet clients (PWA, partners) |
| **mestralive-go** → `pkg/fanout` | Team services publishing opaque `[]byte` after edge auth |

## Install (local replace until published)

```bash
go get github.com/mestracloud/mestralive-go@latest
# during development:
# go.mod: replace github.com/mestralive/mestralive => ../mestralive
```

## Quick start

```go
bus, err := mestralive.Open(mestralive.Config{ServiceToken: os.Getenv("MESTRALIVE_SERVICE_TOKEN")})
if err != nil { log.Fatal(err) }
defer bus.Stop()
_ = bus.Start()

id, _ := bus.Accept()
_ = bus.Subscribe(id, "chat.room")
res, _ := bus.Publish("chat.room", 1, tlvBytes)
_ = res.Accepted // runtime outbound accept, not TCP receipt
```

Or: `mestralive.OpenFromEnv()` with `MESTRALIVE_SERVICE_TOKEN` set.

## Safety

- Do **not** expose raw runtime TCP (`:9000`) on the public internet.
- Do **not** use live JSON-WS for the server hot TLV path.
- `Open` fails closed without a service token.
- Public listen requires explicit opt-in (`AllowPublicListen` / env) — discouraged for v1.
- Dial mode is **not supported** in v1 (`OpenDial` returns an error).

## Auth

1. Issue a service token once (secret manager).
2. Set `MESTRALIVE_SERVICE_TOKEN` for the process.
3. Authenticate **end users** on mestralive-live / your edge.
4. After the session is trusted, call `Publish(tlv)` — no per-message service re-auth.

## Module

`github.com/mestracloud/mestralive-go`
