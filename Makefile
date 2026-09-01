.PHONY: test race example fmt

test:
	go test ./... -count=1

race:
	go test ./... -race -count=1

example:
	@test -n "$$MESTRALIVE_SERVICE_TOKEN" || (echo "export MESTRALIVE_SERVICE_TOKEN=dev-token" && false)
	go run ./examples/inprocess-tlv

fmt:
	go fmt ./...
