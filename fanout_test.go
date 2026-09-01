package mestralive_test

import (
	"os"
	"testing"

	mestralive "github.com/mestracloud/mestralive-go"
)

func TestOpenPublish(t *testing.T) {
	b, err := mestralive.Open(mestralive.Config{ServiceToken: "tok", Owners: 1})
	if err != nil {
		t.Fatal(err)
	}
	defer b.Stop()
	if err := b.Start(); err != nil {
		t.Fatal(err)
	}
	id, err := b.Accept()
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Subscribe(id, "t"); err != nil {
		t.Fatal(err)
	}
	res, err := b.Publish("t", 1, []byte{0x01, 0x02})
	if err != nil {
		t.Fatal(err)
	}
	if res.Subscribers != 1 || res.Accepted < 1 {
		t.Fatalf("%+v", res)
	}
}

func TestOpenFromEnvRequiresToken(t *testing.T) {
	t.Setenv("MESTRALIVE_SERVICE_TOKEN", "")
	if _, err := mestralive.OpenFromEnv(); err == nil {
		t.Fatal("expected error")
	}
	t.Setenv("MESTRALIVE_SERVICE_TOKEN", "env-tok")
	b, err := mestralive.OpenFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	_ = b.Stop()
}

func TestDialUnsupported(t *testing.T) {
	if _, err := mestralive.OpenDial(mestralive.Config{ServiceToken: "x"}); err == nil {
		t.Fatal("expected dial unsupported")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
