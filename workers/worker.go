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
func WorkerKafka(topikKillSig string,result chan models.MessageNats,ks chan string){
	fmt.Println("in kafka worker")
	topicKill:= topikKillSig
	// to consume messages
	p:=utils.GetProducer()
	for{
		select{
		case msgNats:=<-result:
			topic := msgNats.IdChannel
			topicKill = topic
			jsonSenMl:= models.JSONSenML{Valueksql: json.RawMessage(msgNats.MessageData)}
			s,_:=json.Marshal(jsonSenMl)
			println(string(s))
			//println(string(s))
			utils.PublishMessage(s,topic,p)
			glog.Info("Tranfered data to topic %s kafka ",msgNats.IdChannel)
		case signalKill:=<-ks:
			fmt.Println("in kafka worker----")
			if topicKill == signalKill {
				// can xu ly dong kafka neu dung luong nay
				fmt.Println(fmt.Sprintf("Worker %s killed Kafka: ",signalKill))
				return
			}
		}
	}
}
func WorkerNats(job  string, worknumber string, result chan models.MessageNats,nc *nats.Conn,ks chan string,wg *sync.WaitGroup) {
	topicKill := "init"
	fmt.Println(job + "In thread")
	topicKill = job
	subNats, err := nc.QueueSubscribe(job,job, func(m *nats.Msg){
		//wg.Done()
		//fmt.Println("Dât",string(m.Data))
		var message models.Message
		err:=proto.Unmarshal(m.Data,&message)
		if err == nil{
			//fmt.Println("=>>>: ",string(message.Payload))
			result <- models.MessageNats{
				IdChannel:   job,
				MessageData: message.Payload,
			}
		} else {
			fmt.Println(err)
		}
		println(string(message.Payload))
		glog.Info(fmt.Sprintf("Received data from topic %s Nats",job))
	});
	if(err!=nil){
		glog.Error(err)
	}
	defer subNats.Unsubscribe()
	defer wg.Done()
	for true{
		select {
		case signalKill:=<-ks:
			if signalKill == topicKill {
				fmt.Println(fmt.Sprintf("Worker %s killed Nats: ",signalKill))
				return
			}
		}
	}

}
