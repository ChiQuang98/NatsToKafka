package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func worker(jobs chan string, worknumber int, result chan string,ec *nats.EncodedConn,ks chan bool) {
	for true{
		select {
		case job := <- jobs:
			if _, err := ec.QueueSubscribe(job,"worker", func(m *nats.Msg) {
				//wg.Done()
				fmt.Println("Nhan duoc rui",string(m.Data))
				result <- string(m.Data)
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

}
func main() {
	nc, err := nats.Connect("192.168.3.129:4222")
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	// queue of jobs
	jobs := make(chan string)

	// done channel lấy ra kết quả của jobs
	result := make(chan string)
	// số lượng worker trong pool
	//vi` moi worker lam viec khong ket thuc, phai lang nghe lien tuc nen so luong worker bang so luong channel
	for {
		killsignal := make(chan bool)
		fmt.Println("Start")
		numberOfWorkers :=1
		for i := 0; i < numberOfWorkers; i++ {
			go worker(jobs, i, result,ec,killsignal)
		}
		//for {
		numberOfJobs := 1
		for j := 0; j < numberOfJobs; j++ {
			go func(j int) {
				jobs <- "channels.86903597-b75e-4a4d-bbfd-17a304714b86"
			}(j)
		}
		go func() {
			for {
				select {
				case resultMsg := <-result:
					fmt.Println(resultMsg)
				}
			}
		}()
		time.Sleep(10 * time.Second)
		close(killsignal)
	}

		//// chờ nhận đủ kết quả
		//for c := 0; c < numberOfJobs; c++ {
		//	fmt.Println(<-result)
		//}
		//fmt.Println("DONE SESSION")

	//}
}



