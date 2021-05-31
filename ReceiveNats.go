package main

import (
	"NatsToKafka/models"
	"context"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func workerKafka(result chan models.MessageNats,ks chan bool){
	// to consume messages
	for{
		select{
		case msgNats:=<-result:
			//topic := "my-topic"

			topic := msgNats.IdChannel
			partition := 0
			conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
			if err != nil {
				log.Fatal("failed to dial leader:", err)
			}
			conn.SetWriteDeadline(time.Now().Add(10*time.Second))
			_, err = conn.WriteMessages(
				kafka.Message{Value: []byte(msgNats.MessageData)},
			)
			if err != nil {
				log.Fatal("failed to write messages:", err)
			}
			fmt.Println("Tranfered data to topic %s kafka",msgNats.IdChannel)

			if err := conn.Close(); err != nil {
				log.Fatal("failed to close writer:", err)
			}
		case <-ks:
			fmt.Println("Worker halted,: ")
			return
		}
	}
}
func workerNats(jobs chan string, worknumber int, result chan models.MessageNats,ec *nats.EncodedConn,ks chan bool) {
	for true{
		select {
		case job := <- jobs:
			if _, err := ec.Subscribe(job, func(m *nats.Msg) {
				//ec.Flush()
				//wg.Done()
				//fmt.Println("Dât",string(m.Data))
				var message models.Message
				err:=proto.Unmarshal(m.Data,&message)
				if err == nil{
					//fmt.Println("=>>>: ",string(message.Payload))
					result <- models.MessageNats{
						IdChannel:   job,
						MessageData: string(message.Payload),
					}
				} else {
					fmt.Println(err)
				}
				//m := interface{}(m.Data)
				fmt.Println("Received data from topic %s Nats",job)
			}); err != nil {
				log.Fatal(err)
			}
		case <-ks:
			fmt.Println("Worker halted,: ")
			return
		}
	}

}
func ConnectNats()  {
//https://www.digitalocean.com/community/tutorials/how-to-install-apache-kafka-on-ubuntu-20-04
}
func main() {
	//Nats
	nc, err := nats.Connect("192.168.3.129:4222")

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	// queue of jobs
	jobs := make(chan string)

	// done channel lấy ra kết quả của jobs
	result := make(chan models.MessageNats)
	// số lượng worker trong pool
	//vi` moi worker lam viec khong ket thuc, phai lang ng  he lien tuc nen so luong worker bang so luong channel
	for {
		killsignal := make(chan bool)
		fmt.Println("Start")
		numberOfWorkers :=1
		for i := 0; i < numberOfWorkers; i++ {
			go workerNats(jobs, i, result,ec,killsignal)
			go workerKafka(result,killsignal)
		}
		//for {
		numberOfJobs := 1
		for j := 0; j < numberOfJobs; j++ {
			go func(j int) {
				jobs <- "channels.86903597-b75e-4a4d-bbfd-17a304714b86"
			}(j)
		}
		//go func() {
		//	for {
		//		select {
		//		case resultMsg := <-result:
		//			//fmt.Println(resultMsg)
		//		}
		//	}
		//}()
		time.Sleep(10 * time.Hour)
		close(killsignal)
	}

}



