package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
)

func natsErrHandler21(nc *nats.Conn, sub *nats.Subscription, natsErr error) {
	fmt.Printf("error: %v\n", natsErr)
	if natsErr == nats.ErrSlowConsumer {
		pendingMsgs, _, err := sub.Pending()
		if err != nil {
			fmt.Printf("couldn't get pending messages: %v", err)
			return
		}
		fmt.Printf("Falling behind with %d pending messages on subject %q.\n",
			pendingMsgs, sub.Subject)
		// Log error, notify operations...
	}
	// check for other errors
}
func main() {
	nc, err := nats.Connect("10.16.150.132:4222",nats.ErrorHandler(natsErrHandler21))
	if err != nil {
		glog.Error(err)
	} else {
		glog.Info(fmt.Sprintf("Connected to Nats Server at %s","10.16.150.132:4222"))
	}
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Subscribe
	if _, err := nc.Subscribe("channels.963f4b95-6a2c-49e4-a300-21d231c0f8c6", func(m *nats.Msg) {
		fmt.Println("Receiveed")
	}); err != nil {
		log.Fatal(err)
	}

	// Wait for a message to come in
	wg.Wait()
}
