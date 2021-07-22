package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)
func natsErrHandler1(nc *nats.Conn, sub *nats.Subscription, natsErr error) {
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
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "10.16.150.138:6667",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"channels.7b0fa712-3301-44f3-a26d-d5e8dc1b9aea", "^aRegex.*[Tt]opic"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()

	//wg := sync.WaitGroup{}
	//wg.Add(10)
	//
	//nc, err := nats.Connect("10.16.150.132:4222",nats.ErrorHandler(natsErrHandler1))
	//if err != nil {
	//	glog.Error(err)
	//} else {
	//	glog.Info(fmt.Sprintf("Connected to Nats Server at %s","10.16.150.132:4222"))
	//}
	//defer nc.Close()
	////for{
	//	if _, err := nc.Subscribe("channels.7b0fa712-3301-44f3-a26d-d5e8dc1b9aea", func(m *nats.Msg) {
	//		wg.Done()
	//		//fmt.Println("Dât",string(m.Data))
	//		var message models.Message
	//		err:=proto.Unmarshal(m.Data,&message)
	//		if err == nil{
	//			//fmt.Println("=>>>: ",string(message.Payload))
	//			//result <- models.MessageNats{
	//			//	IdChannel:   job,
	//			//	MessageData: string(message.Payload),
	//			//}
	//			fmt.Println(string(message.Payload))
	//			//atomic.AddUint64(&total, 1)
	//		} else {
	//			fmt.Println(err)
	//		}
	//		//m := interface{}(m.Data)
	//		glog.Info(fmt.Sprintf("Received data from topic %s Nats","channels.7b0fa712-3301-44f3-a26d-d5e8dc1b9aea"))
	//		//fmt.Println(string(message.Payload))
	//	}); err != nil {
	//		glog.Error(err)
	//	}
	////}
	//wg.Wait()

}
