package workers

import (
	"NatsToKafka/models"
	"NatsToKafka/utils"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/nats-io/nats.go"
	"sync"
)
var total uint64
func WorkerKafka(result chan models.MessageNats,ks chan bool){
	// to consume messages
	p:=utils.GetProducer()
	for{
		select{
		case msgNats:=<-result:
			//topic := "my-topic"

			topic := msgNats.IdChannel
			//partition := 0
			//conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
			//if err != nil {
			//	glog.Error("failed to dial leader:", err)
			//}
			//conn.SetWriteDeadline(time.Now().Add(10*time.Second))
			utils.PublishMessage(msgNats.MessageData,topic,p)
			//if err != nil {
			//	glog.Error("failed to write messages:", err)
			//}
			glog.Info("Tranfered data to topic %s kafka",msgNats.IdChannel)

			//if err := conn.Close(); err != nil {
			//	glog.Error("failed to close writer:", err)
			//}
		case <-ks:
			fmt.Println("Worker halted,: ")
			return
		}
	}
}
func WorkerNats(jobs chan string, worknumber int, result chan models.MessageNats,nc *nats.Conn,ks chan bool,wg *sync.WaitGroup) {
	defer wg.Done()
	for true{
		select {
		case job := <- jobs:
			//fmt.Println(job)
			if _, err := nc.QueueSubscribe(job,job, func(m *nats.Msg) {
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
					//fmt.Println(string(message.Payload))
					//atomic.AddUint64(&total, 1)
				} else {
					fmt.Println(err)
				}
				//m := interface{}(m.Data)
				glog.Info(fmt.Sprintf("Received data from topic %s Nats",job))
				//fmt.Println(string(message.Payload))
			}); err != nil {
				glog.Error(err)
			}
		case <-ks:
			fmt.Println("Worker halted,: ")
			return
		}
	}

}
