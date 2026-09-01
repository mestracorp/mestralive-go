// Example: in-process TLV publish/subscribe via mestralive-go.
package main

import (
	"fmt"
	"log"
	"os"

	mestralive "github.com/mestracorp/mestralive-go"
)

func main() {
	tok := os.Getenv("MESTRALIVE_SERVICE_TOKEN")
	if tok == "" {
		tok = "dev-token"
	}
	bus, err := mestralive.Open(mestralive.Config{ServiceToken: tok, Owners: 1})
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Stop()
	if err := bus.Start(); err != nil {
		log.Fatal(err)
	}

	a, err := bus.Accept()
	if err != nil {
		log.Fatal(err)
	}
	b, err := bus.Accept()
	if err != nil {
		log.Fatal(err)
	}
	_ = bus.Subscribe(a, "demo")
	_ = bus.Subscribe(b, "demo")

	tlv := []byte{0x01, 0x02, 0x03, 0xff}
	res, err := bus.Publish("demo", 1, tlv)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("subscribers=%d accepted=%d rejected=%d\n", res.Subscribers, res.Accepted, res.Rejected)
}
