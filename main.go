package main

import (
	"NatsToKafka/models"
	"NatsToKafka/utils"
	"NatsToKafka/workers"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/nats-io/nats.go"
	"os"
	"sync"
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
	nc, err := nats.Connect("10.16.150.132:4222",nats.ErrorHandler(natsErrHandler),nats.PingInterval(20*time.Second), nats.MaxPingsOutstanding(5))

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
		go workers.WorkerNats(jobs, i, result,nc,killsignal,&wg)
		go workers.WorkerKafka(result,killsignal)
	}
	e:=nc.Flush()
	if e!=nil{
		glog.Error("Flush error")
	}
	//for {
	numberOfJobs := numberOfWorkers
	fmt.Println("LEN CHAN: ", numberOfJobs)
	//numberOfJobs := 1
	for j := 0; j < numberOfJobs; j++ {
		go func(j int) {
			jobs <- "channels."+channelsAtomic.channels[j].Channel_id
			//fmt.Println("channels."+channelsAtomic.channels[j].Id)
		}(j)
	}
	//go func() {
	//	for  {
	//		fmt.Println("=====Tong cong: ",total)
	//		time.Sleep(5*time.Second)
	//	}
	//}()
	wg.Wait()

}



