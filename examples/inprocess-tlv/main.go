// Example: publish opaque TLV bytes to two in-process subscribers.
//
//	MESTRALIVE_SERVICE_TOKEN=dev-token go run .
package main

import (
	"fmt"
	"log"
	"os"

	mestralive "github.com/mestracorp/mestralive-go"
)

func main() {
	token := os.Getenv("MESTRALIVE_SERVICE_TOKEN")
	if token == "" {
		log.Fatal("set MESTRALIVE_SERVICE_TOKEN")
	}

	bus, err := mestralive.Open(mestralive.Config{
		ServiceToken: token,
		Owners:       1,
	})
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer bus.Stop()

	if err := bus.Start(); err != nil {
		log.Fatalf("start: %v", err)
	}

	alice, err := bus.Accept()
	if err != nil {
		log.Fatalf("accept alice: %v", err)
	}
	bob, err := bus.Accept()
	if err != nil {
		log.Fatalf("accept bob: %v", err)
	}

	const topic = "demo.room"
	if err := bus.Subscribe(alice, topic); err != nil {
		log.Fatalf("subscribe alice: %v", err)
	}
	if err := bus.Subscribe(bob, topic); err != nil {
		log.Fatalf("subscribe bob: %v", err)
	}

	// Opaque application bytes — not JSON.
	tlv := []byte{0x01, 0x02, 0x03, 0xff}
	res, err := bus.Publish(topic, 1, tlv)
	if err != nil {
		log.Fatalf("publish: %v", err)
	}

	fmt.Printf("subscribers=%d accepted=%d rejected=%d disconnected=%d\n",
		res.Subscribers, res.Accepted, res.Rejected, res.Disconnected)
}
