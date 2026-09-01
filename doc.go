// Package mestralive is the Go SDK for in-process mestralive pub/sub.
//
// # Who this is for
//
// Backend Go services that need to publish opaque binary payloads (TLV, AWP,
// protobuf, raw bytes) through an embedded mestralive runtime — without JSON
// WebSocket or HTTP on the hot path.
//
// Browser and partner clients should use mestralive-live (JSON over TLS), not
// this package.
//
// # Minimal flow
//
//	bus, err := mestralive.Open(mestralive.Config{ServiceToken: token})
//	if err != nil { ... }
//	defer bus.Stop()
//	_ = bus.Start()
//	id, _ := bus.Accept()
//	_ = bus.Subscribe(id, "chat.room")
//	res, _ := bus.Publish("chat.room", 1, tlvBytes)
//
// Service authentication happens once at Open. End-user authentication happens
// at your edge; enforce product authorization before Publish.
//
// See the repository README for install, auth, testing, and FAQ.
package mestralive
