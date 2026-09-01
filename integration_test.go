package mestralive_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	mestralive "github.com/mestracorp/mestralive-go"
)

func TestContract_OpenSubPubStop(t *testing.T) {
	bus, err := mestralive.Open(mestralive.Config{ServiceToken: "tok", Owners: 2, MaxConnections: 64})
	if err != nil {
		t.Fatal(err)
	}
	if err := bus.Start(); err != nil {
		t.Fatal(err)
	}
	a, err := bus.Accept()
	if err != nil {
		t.Fatal(err)
	}
	b, err := bus.Accept()
	if err != nil {
		t.Fatal(err)
	}
	if err := bus.Subscribe(a, "room.a"); err != nil {
		t.Fatal(err)
	}
	if err := bus.Subscribe(b, "room.a"); err != nil {
		t.Fatal(err)
	}
	res, err := bus.Publish("room.a", 1, []byte{0x00, 0xff})
	if err != nil {
		t.Fatal(err)
	}
	if res.Subscribers != 2 || res.Attempted != 2 || res.Accepted < 1 {
		t.Fatalf("%+v", res)
	}
	if err := bus.Stop(); err != nil {
		t.Fatal(err)
	}
	if _, err := bus.Publish("room.a", 1, []byte("x")); err == nil {
		t.Fatal("expected error after Stop")
	}
}

func TestContract_MultiTopicIsolation(t *testing.T) {
	bus, err := mestralive.Open(mestralive.Config{ServiceToken: "tok", Owners: 1, MaxConnections: 16})
	if err != nil {
		t.Fatal(err)
	}
	defer bus.Stop()
	if err := bus.Start(); err != nil {
		t.Fatal(err)
	}
	a, _ := bus.Accept()
	b, _ := bus.Accept()
	_ = bus.Subscribe(a, "topic.one")
	_ = bus.Subscribe(b, "topic.two")
	r1, err := bus.Publish("topic.one", 1, []byte("1"))
	if err != nil {
		t.Fatal(err)
	}
	r2, err := bus.Publish("topic.two", 1, []byte("2"))
	if err != nil {
		t.Fatal(err)
	}
	if r1.Subscribers != 1 || r2.Subscribers != 1 {
		t.Fatalf("isolation broken: %+v %+v", r1, r2)
	}
}

func TestContract_ConcurrentPublishers(t *testing.T) {
	bus, err := mestralive.Open(mestralive.Config{
		ServiceToken:        "tok",
		Owners:              2,
		MaxConnections:      64,
		CrossQueueCapacity:  1024,
		MaxWriteBufferBytes: 1 << 20,
		MaxFanoutPerPublish: 32,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer bus.Stop()
	if err := bus.Start(); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 8; i++ {
		id, err := bus.Accept()
		if err != nil {
			t.Fatal(err)
		}
		if err := bus.Subscribe(id, "conc"); err != nil {
			t.Fatal(err)
		}
	}
	var wg sync.WaitGroup
	var fails atomic.Uint64
	for p := 0; p < 8; p++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				if _, err := bus.Publish("conc", 1, []byte("x")); err != nil {
					fails.Add(1)
				}
			}
		}()
	}
	wg.Wait()
	if fails.Load() != 0 {
		t.Fatalf("publish errors: %d", fails.Load())
	}
}

func TestContract_StopUnderLoad(t *testing.T) {
	bus, err := mestralive.Open(mestralive.Config{ServiceToken: "tok", Owners: 1, MaxConnections: 16})
	if err != nil {
		t.Fatal(err)
	}
	if err := bus.Start(); err != nil {
		t.Fatal(err)
	}
	id, _ := bus.Accept()
	_ = bus.Subscribe(id, "load")
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				_, _ = bus.Publish("load", 1, []byte("y"))
			}
		}
	}()
	if err := bus.Stop(); err != nil {
		t.Fatal(err)
	}
	close(done)
}

func TestSecurity_StolenEmptyToken(t *testing.T) {
	if err := openErr(mestralive.Config{}); err == nil {
		t.Fatal("empty token")
	}
	err := openErr(mestralive.Config{ServiceToken: "stolen", ServiceTokenAllowlist: []string{"real"}})
	if !errors.Is(err, mestralive.ErrUnauthorized) {
		t.Fatalf("got %v", err)
	}
}

func TestSecurity_PublicListenDenied(t *testing.T) {
	err := openErr(mestralive.Config{ServiceToken: "tok", ListenAddress: "0.0.0.0:9201"})
	if !errors.Is(err, mestralive.ErrUnsafeListen) {
		t.Fatalf("got %v", err)
	}
}

func TestSecurity_OversizedPayload(t *testing.T) {
	bus, err := mestralive.Open(mestralive.Config{ServiceToken: "tok", Owners: 1, MaxMessageBytes: 4})
	if err != nil {
		t.Fatal(err)
	}
	defer bus.Stop()
	_ = bus.Start()
	_, err = bus.Publish("t", 1, []byte("12345"))
	if !errors.Is(err, mestralive.ErrPayloadTooLarge) {
		t.Fatalf("got %v", err)
	}
}

func TestSecurity_SubscribeStormLimited(t *testing.T) {
	bus, err := mestralive.Open(mestralive.Config{ServiceToken: "tok", Owners: 1, MaxConnections: 3})
	if err != nil {
		t.Fatal(err)
	}
	defer bus.Stop()
	_ = bus.Start()
	for i := 0; i < 3; i++ {
		if _, err := bus.Accept(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := bus.Accept(); err == nil {
		t.Fatal("expected admission limit")
	}
}

func openErr(cfg mestralive.Config) error {
	_, err := mestralive.Open(cfg)
	return err
}
