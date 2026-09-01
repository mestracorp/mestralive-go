package mestralive_test

import (
	"encoding/json"
	"testing"

	mestralive "github.com/mestracorp/mestralive-go"
)

// Phase 3.4: SDK Publish path should stay near pkg/fanout (no extra JSON).
func BenchmarkSDK_PublishTLV(b *testing.B) {
	bus, err := mestralive.Open(mestralive.Config{
		ServiceToken:        "tok",
		Owners:              1,
		MaxConnections:      8,
		CrossQueueCapacity:  4096,
		MaxWriteBufferBytes: 4 << 20,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer bus.Stop()
	if err := bus.Start(); err != nil {
		b.Fatal(err)
	}
	id, err := bus.Accept()
	if err != nil {
		b.Fatal(err)
	}
	if err := bus.Subscribe(id, "sdk"); err != nil {
		b.Fatal(err)
	}
	pl := make([]byte, 215)
	if _, err := bus.Publish("sdk", 1, pl); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bus.Publish("sdk", 1, pl)
	}
}

func BenchmarkSDK_JSONThenPublish(b *testing.B) {
	bus, err := mestralive.Open(mestralive.Config{
		ServiceToken:        "tok",
		Owners:              1,
		MaxConnections:      8,
		CrossQueueCapacity:  4096,
		MaxWriteBufferBytes: 4 << 20,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer bus.Stop()
	if err := bus.Start(); err != nil {
		b.Fatal(err)
	}
	id, err := bus.Accept()
	if err != nil {
		b.Fatal(err)
	}
	if err := bus.Subscribe(id, "sdk"); err != nil {
		b.Fatal(err)
	}
	type env struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}
	msg := env{Type: "publish", Data: string(make([]byte, 215))}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		raw, err := json.Marshal(msg)
		if err != nil {
			b.Fatal(err)
		}
		_, _ = bus.Publish("sdk", 1, raw)
	}
}
