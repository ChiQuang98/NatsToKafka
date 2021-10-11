package workers

import (
	"NatsToKafka/models"
	"NatsToKafka/utils"
	"encoding/json"
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
			topic := msgNats.IdChannel
			jsonSenMl:= models.JSONSenML{Value: msgNats.MessageData}
			s,_:=json.Marshal(jsonSenMl)
			utils.PublishMessage(s,topic,p)
			glog.Info("Tranfered data to topic %s kafka",msgNats.IdChannel)
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
			fmt.Println(job)
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
				} else {
					fmt.Println(err)
				}
				println(string(message.Payload))
				glog.Info(fmt.Sprintf("Received data from topic %s Nats",job))
			}); err != nil {
				glog.Error(err)
			}
		case <-ks:
			fmt.Println("Worker halted,: ")
			return
		}
	}

}
