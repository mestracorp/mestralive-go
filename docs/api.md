# API reference (developer)

Full Go docs: `go doc github.com/mestracorp/mestralive-go`  
Source of types: this module re-exports `pkg/fanout` from mestralive.

## Open / lifecycle

```go
Open(Config) (*Bus, error)
OpenFromEnv() (*Bus, error)    // reads MESTRALIVE_SERVICE_TOKEN
OpenDial(Config) (*Bus, error) // always errors in v1

(*Bus) Start() error
(*Bus) Stop() error            // idempotent
```

## Pub/sub

```go
(*Bus) Accept() (ConnID, error)
(*Bus) Subscribe(conn ConnID, topic string) error
(*Bus) Unsubscribe(conn ConnID, topic string) error
(*Bus) Publish(topic string, typ uint16, payload []byte) (Result, error)
```

- `typ` is your application message type (uint16). The bus does not interpret it.  
- `payload` is opaque. Binary safe (including `0x00`).

## Result

| Field | Meaning |
|-------|---------|
| `Subscribers` | How many matched the topic before the fanout budget |
| `Attempted` | How many were tried this call |
| `Accepted` | Queued into runtime outbound path |
| `Rejected` | Refused (full buffer / policy) |
| `Disconnected` | Target connection gone |
| `Deferred` | Not tried (fanout budget) |
| `FanoutLimited` | Budget truncated the subscriber list |

## Config

| Field | Required | Notes |
|-------|----------|-------|
| `ServiceToken` | **yes** | Process credential |
| `ServiceTokenAllowlist` | no | If set, token must be listed |
| `Owners` | no | Default 2 |
| `MaxConnections` | no | Default 1024 |
| `MaxFanoutPerPublish` | no | 0 = runtime default |
| `MaxMessageBytes` | no | Enforced on Publish when set |
| `CrossQueueCapacity` | no | Mailbox depth |
| `MaxWriteBufferBytes` | no | Per-connection outbound bound |
| `ListenAddress` | no | Empty = in-process only |
| `AllowPublicListen` | no | Required for non-loopback listen |

## Errors

See README “Errors you will see”, or constants on the package:  
`ErrUnauthorized`, `ErrUnsafeListen`, `ErrNotStarted`, `ErrClosed`, `ErrPayloadTooLarge`, `ErrInvalidTopic`, `ErrDialUnsupported`.
