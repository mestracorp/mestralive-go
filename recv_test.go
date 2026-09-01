package mestralive_test

import (
	"bytes"
	"testing"

	mestralive "github.com/mestracorp/mestralive-go"
)

func TestSDK_RecvTwoSubscribers(t *testing.T) {
	bus, err := mestralive.Open(mestralive.Config{ServiceToken: "tok", Owners: 1})
	if err != nil {
		t.Fatal(err)
	}
	defer bus.Stop()
	if err := bus.Start(); err != nil {
		t.Fatal(err)
	}
	a, _ := bus.Accept()
	b, _ := bus.Accept()
	_ = bus.Subscribe(a, "room")
	_ = bus.Subscribe(b, "room")
	pl := bytes.Repeat([]byte{0xab}, 215)
	res, err := bus.Publish("room", 9, pl)
	if err != nil || res.Accepted != 2 {
		t.Fatalf("pub %+v err=%v", res, err)
	}
	buf := make([]byte, 215)
	n, typ, err := bus.Recv(a, buf)
	if err != nil || n != 215 || typ != 9 || !bytes.Equal(buf, pl) {
		t.Fatalf("a n=%d typ=%d err=%v", n, typ, err)
	}
	n, typ, err = bus.Recv(b, buf)
	if err != nil || n != 215 || typ != 9 || !bytes.Equal(buf, pl) {
		t.Fatalf("b n=%d typ=%d err=%v", n, typ, err)
	}
}
