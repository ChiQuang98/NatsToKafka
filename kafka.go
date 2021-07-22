package main

import (
	"NatsToKafka/utils"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

//func produce()  {
//	// to produce messages
//	topic := "TestKafka"
//	partition := 2
//
//	conn, err := kafka.DialLeader(context.Background(), "tcp", "10.16.150.138:6667", topic, partition)
//	if err != nil {
//		log.Fatal("failed to dial leader:", err)
//	}
//
//	conn.SetWriteDeadline(time.Now().Add(10*time.Second))
//
//	_, err = conn.WriteMessages(
//		kafka.Message{Value: []byte("Quang1!")},
//		//kafka.Message{Value: []byte("Quang2!")},
//		//kafka.Message{Value: []byte("Quang3!")},
//	)
//	if err != nil {
//		log.Fatal("failed to write messages:", err)
//	}
//
//	if err := conn.Close(); err != nil {
//		log.Fatal("failed to close writer:", err)
//	}
//}
//func  consume()  {
//	// to consume messages
//	topic := "TestKafka"
//	partition := 2
//
//	conn, err := kafka.DialLeader(context.Background(), "tcp", "10.16.150.138:6667", topic, partition)
//	if err != nil {
//		log.Fatal("failed to dial leader:", err)
//	}
//
//	conn.SetReadDeadline(time.Now().Add(10*time.Second))
//	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
//
//	b := make([]byte, 10e3) // 10KB max per message
//	for {
//		_, err := batch.Read(b)
//		if err != nil {
//			break
//		}
//		fmt.Println(string(b))
//	}
//
//	if err := batch.Close(); err != nil {
//		log.Fatal("failed to close batch:", err)
//	}
//
//	if err := conn.Close(); err != nil {
//		log.Fatal("failed to close connection:", err)
//	}
//}
func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "10.16.150.138:6667"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	//topic := "myTopic"
	//for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "cliejjjjjjjjjjjjjjnt"} {
	//	p.Produce(&kafka.Message{
	//		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	//		Value:          []byte(word),
	//	}, nil)
	//}
	utils.PublishMessage("QUANGCHI1333323","QUANGtopic",p)
	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}