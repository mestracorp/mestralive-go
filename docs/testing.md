# Testing

## Commands

```bash
# All tests
go test ./... -count=1

# Race detector (recommended in CI)
go test ./... -race -count=1

# Verbose
go test ./... -v -count=1

# Example smoke
MESTRALIVE_SERVICE_TOKEN=dev-token go run ./examples/inprocess-tlv
```

## What ships in this repo

| File | Focus |
|------|--------|
| `fanout_test.go` | Open, env token, Dial unsupported, basic publish |
| `integration_test.go` | Multi-topic, concurrency, stop-under-load, security limits |
| `bench_test.go` | Optional microbenches (lab only) |

## Writing tests in your product

Pattern used by the Agyraa spike:

1. `Open` with a test service token  
2. `Start`  
3. `Accept` + `Subscribe` for fake subscribers  
4. `Publish` your TLV  
5. Assert `Result.Accepted` / errors (`ErrUnauthorized`, capability denied in *your* layer)  
6. `Stop`  

```go
func TestMyPublish(t *testing.T) {
	bus, err := mestralive.Open(mestralive.Config{ServiceToken: "test", Owners: 1})
	if err != nil {
		t.Fatal(err)
	}
	defer bus.Stop()
	if err := bus.Start(); err != nil {
		t.Fatal(err)
	}

	id, err := bus.Accept()
	if err != nil {
		t.Fatal(err)
	}
	if err := bus.Subscribe(id, "conv.1"); err != nil {
		t.Fatal(err)
	}

	res, err := bus.Publish("conv.1", 1, []byte{0x01, 0x02})
	if err != nil {
		t.Fatal(err)
	}
	if res.Accepted < 1 {
		t.Fatalf("expected accept, got %+v", res)
	}
}
```

## CI tips

- Always run `-race` on PRs that touch the bus.  
- Do not treat microbench numbers as capacity SLOs.  
- If you use `replace` for local mestralive, document it in your repo README so CI checkouts stay reproducible.
