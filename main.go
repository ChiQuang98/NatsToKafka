package main

import (
	"NatsToKafka/models"
	"NatsToKafka/utils"
	"context"
	"flag"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/nats-io/nats.go"
	"github.com/segmentio/kafka-go"
	"os"
	"sync"
	"sync/atomic"
	"time"
)
func init() {
	os.Mkdir("./logs", 0777)
	flag.Lookup("stderrthreshold").Value.Set("[INFO|WARN|FATAL]")
	flag.Lookup("logtostderr").Value.Set("false")
	flag.Lookup("alsologtostderr").Value.Set("true")
	flag.Lookup("log_dir").Value.Set("./logs")
	glog.MaxSize = 1024 * 1024 * 256
	flag.Lookup("v").Value.Set(fmt.Sprintf("%d", 8))
	flag.Parse()
}
var total uint64

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
				glog.Error("failed to dial leader:", err)
			}
			//conn.SetWriteDeadline(time.Now().Add(10*time.Second))
			_, err = conn.WriteMessages(
				kafka.Message{Value: []byte(msgNats.MessageData)},
			)
			if err != nil {
				glog.Error("failed to write messages:", err)
			}
			glog.Info("Tranfered data to topic %s kafka",msgNats.IdChannel)

			if err := conn.Close(); err != nil {
				glog.Error("failed to close writer:", err)
			}
		case <-ks:
			fmt.Println("Worker halted,: ")
			return
		}
	}
}
func workerNats(jobs chan string, worknumber int, result chan models.MessageNats,nc *nats.Conn,ks chan bool,wg *sync.WaitGroup) {
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
					//result <- models.MessageNats{
					//	IdChannel:   job,
					//	MessageData: string(message.Payload),
					//}
					atomic.AddUint64(&total, 1)
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

func ConnectNats()  {
	//https://www.digitalocean.com/community/tutorials/how-to-install-apache-kafka-on-ubuntu-20-04
}
func natsErrHandler(nc *nats.Conn, sub *nats.Subscription, natsErr error) {
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

// Set the error handler when creating a connection.

var channelsAtomic struct {
	sync.Mutex
	channels []models.Channel
}
func main() {
	var channelsTmp []models.Channel
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		var err error
		channelsAtomic.Lock()
		channelsTmp, err= utils.GetAllChannels()
		channelsAtomic.channels = channelsTmp
		channelsAtomic.Unlock()
		if err !=nil{
			glog.Error("fail to get all chanels")
			return
		}
		wg.Done()
		//time.Sleep(60*time.Second)
	}()
	wg.Wait()
	fmt.Println(len(channelsAtomic.channels))
	//Nats
	nc, err := nats.Connect("10.16.150.132:4222",nats.ErrorHandler(natsErrHandler))

	//ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		glog.Error(err)
	} else {
		glog.Info(fmt.Sprintf("Connected to Nats Server at %s","10.16.150.132:4222"))
	}
	defer nc.Close()
	// queue of jobs
	jobs := make(chan string)

	// done channel lấy ra kết quả của jobs
	result := make(chan models.MessageNats)
	// số lượng worker trong pool
	//vi` moi worker lam viec khong ket thuc, phai lang ng  he lien tuc nen so luong worker bang so luong channel
	killsignal := make(chan bool)
	fmt.Println("Start")
	numberOfWorkers :=len(channelsAtomic.channels)
	//numberOfWorkers :=1
	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go workerNats(jobs, i, result,nc,killsignal,&wg)
		//go workerKafka(result,killsignal)
	}
	//e:=nc.Flush()
	//if e!=nil{
	//	glog.Error("Flush error")
	//}
	//for {
	numberOfJobs := numberOfWorkers
	fmt.Println("LEN CHAN: ", numberOfJobs)
	//numberOfJobs := 1
	for j := 0; j < numberOfJobs; j++ {
		go func(j int) {
			jobs <- "channels."+channelsAtomic.channels[j].Id
			//fmt.Println("channels."+channels[j].Id)
		}(j)
	}
	go func() {
		for  {
			fmt.Println("=====Tong cong: ",total)
			time.Sleep(5*time.Second)
		}
	}()
	wg.Wait()

}



