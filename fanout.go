// Package mestralive provides the team-facing Go SDK for in-process TLV fanout.
//
// Internet clients should attach via mestralive-live (JSON-WS). Hot-path server
// publishes use opaque []byte through this SDK → pkg/fanout → certified runtime.
package mestralive

import (
	"errors"
	"os"

	"github.com/mestralive/mestralive/pkg/fanout"
)

// Re-export core types so callers need one import.
type (
	Bus    = fanout.Bus
	Config = fanout.Config
	ConnID = fanout.ConnID
	Result = fanout.Result
)

var (
	ErrUnauthorized    = fanout.ErrUnauthorized
	ErrUnsafeListen    = fanout.ErrUnsafeListen
	ErrNotStarted      = fanout.ErrNotStarted
	ErrClosed          = fanout.ErrClosed
	ErrNotSupported    = fanout.ErrNotSupported
	ErrPayloadTooLarge = fanout.ErrPayloadTooLarge
	ErrInvalidTopic    = fanout.ErrInvalidTopic
	ErrDialUnsupported = errors.New("mestralive-go: dial mode not supported in v1 (use InProcess)")
)

// Mode selects how the SDK reaches the bus.
type Mode int

const (
	// ModeInProcess embeds pkg/fanout in this process (v1 default; lowest latency).
	ModeInProcess Mode = iota
	// ModeDial is reserved for a future private TCP/UDS path.
	ModeDial
)

// Open opens an in-process bus (ModeInProcess).
func Open(cfg Config) (*Bus, error) {
	return fanout.Open(cfg)
}

// OpenFromEnv loads MESTRALIVE_SERVICE_TOKEN (required) and optional listen/owners.
func OpenFromEnv() (*Bus, error) {
	tok := os.Getenv("MESTRALIVE_SERVICE_TOKEN")
	cfg := Config{
		ServiceToken:   tok,
		ListenAddress:  os.Getenv("MESTRALIVE_FANOUT_LISTEN"),
		AllowPublicListen: os.Getenv("MESTRALIVE_FANOUT_ALLOW_PUBLIC") == "1",
	}
	return Open(cfg)
}

// OpenDial always fails in v1 — kept so callers get an explicit error.
func OpenDial(_ Config) (*Bus, error) {
	return nil, ErrDialUnsupported
}
