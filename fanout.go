package mestralive

import (
	"errors"
	"os"

	"github.com/mestralive/mestralive/pkg/fanout"
)

// Re-exported types (stable names for application imports).
type (
	// Bus is an in-process pub/sub handle.
	Bus = fanout.Bus
	// Config configures Open.
	Config = fanout.Config
	// ConnID identifies an in-process connection created by Accept.
	ConnID = fanout.ConnID
	// Result is the outcome of one Publish.
	Result = fanout.Result
)

// Re-exported sentinel errors.
var (
	ErrUnauthorized    = fanout.ErrUnauthorized
	ErrUnsafeListen    = fanout.ErrUnsafeListen
	ErrNotStarted      = fanout.ErrNotStarted
	ErrClosed          = fanout.ErrClosed
	ErrNotSupported    = fanout.ErrNotSupported
	ErrPayloadTooLarge = fanout.ErrPayloadTooLarge
	ErrInvalidTopic    = fanout.ErrInvalidTopic
	ErrShortBuffer     = fanout.ErrShortBuffer
	ErrUnknownConn     = fanout.ErrUnknownConn
	// ErrDialUnsupported is returned by OpenDial in v1.
	ErrDialUnsupported = errors.New("mestralive-go: dial mode not supported in v1 (use in-process Open)")
)

// Mode selects how the SDK reaches the bus.
type Mode int

const (
	// ModeInProcess embeds the bus in this process (v1 default).
	ModeInProcess Mode = iota
	// ModeDial is reserved. OpenDial returns ErrDialUnsupported in v1.
	ModeDial
)

// Open creates a bus using ModeInProcess. The bus is not started until Start.
// ServiceToken is required (see docs/authentication.md).
func Open(cfg Config) (*Bus, error) {
	return fanout.Open(cfg)
}

// OpenFromEnv is Open with ServiceToken from MESTRALIVE_SERVICE_TOKEN.
//
// Optional:
//
//	MESTRALIVE_FANOUT_LISTEN         — runtime ListenAddress (prefer empty)
//	MESTRALIVE_FANOUT_ALLOW_PUBLIC=1 — allow non-loopback listen (discouraged)
func OpenFromEnv() (*Bus, error) {
	cfg := Config{
		ServiceToken:      os.Getenv("MESTRALIVE_SERVICE_TOKEN"),
		ListenAddress:     os.Getenv("MESTRALIVE_FANOUT_LISTEN"),
		AllowPublicListen: os.Getenv("MESTRALIVE_FANOUT_ALLOW_PUBLIC") == "1",
	}
	return Open(cfg)
}

// OpenDial is not supported in v1. Prefer Open (in-process).
func OpenDial(_ Config) (*Bus, error) {
	return nil, ErrDialUnsupported
}
