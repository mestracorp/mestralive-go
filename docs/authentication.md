# Authentication

Two different identities — do not mix them.

## 1. Service token (required by this SDK)

Proves **your Go process** may open a fanout bus.

| | |
|--|--|
| When | Once, at `Open` / `OpenFromEnv` |
| Where | Server secret manager / env |
| Env var | `MESTRALIVE_SERVICE_TOKEN` |
| Failure | `ErrUnauthorized` — bus does not open |

### Production tip

Prefer an allowlist so a random non-empty string is not enough:

```go
mestralive.Open(mestralive.Config{
	ServiceToken:          os.Getenv("MESTRALIVE_SERVICE_TOKEN"),
	ServiceTokenAllowlist: []string{os.Getenv("MESTRALIVE_SERVICE_TOKEN")},
})
```

In single-tenant lab mode, any non-empty token is accepted if the allowlist is empty.

### Rotation

1. Issue a new token.  
2. Deploy it to the secret store / env.  
3. Restart the service.  

v1 does not hot-reload tokens mid-process.

### Never

- Ship the service token to browsers or mobile apps  
- Log the raw token  
- Send the token on every `Publish`

## 2. End-user / partner auth (outside this SDK)

Proves **who the human or partner app is**.

| | |
|--|--|
| When | WebSocket / HTTP session to **mestralive-live** (or your TLS edge) |
| Mechanism | Bearer `key_id.secret`, identity tokens (`idt.`), etc. |
| After success | Your service maps the session to allowed topics (capability set) |

Then call:

```go
// sess already authenticated at the edge
if !sess.CanPublish(topic) {
	return ErrDenied
}
bus.Publish(topic, typ, tlv)
```

## Optional listen auth note

If you set `ListenAddress` (not recommended for v1 product paths), that opens a **runtime TCP listener**. It is still not the public JSON client API. Public binds require `AllowPublicListen` and are discouraged. Prefer `ListenAddress: ""` (in-process only).
