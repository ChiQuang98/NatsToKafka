package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func worker(jobs chan string, worknumber int, result chan string,nc *nats.Conn) {
	for channel := range jobs {
		//fmt.Println("worker", worknumber, "finished job", channel)
		//for {
			if _, err := nc.QueueSubscribe(channel,"worker", func(m *nats.Msg) {
				//wg.Done()
				//fmt.Println("Nhan duoc rui",string(m.Data))
				result <- string(m.Data)
			}); err != nil {
				log.Fatal(err)
			}
		//}
	}
}
func ConnectNats()  {

}
func main() {
	nc, err := nats.Connect("10.4.200.61:4222")
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
	numberOfWorkers :=1
	for i := 0; i < numberOfWorkers; i++ {
		go worker(jobs, i, result,nc)
	}
	//for {
		numberOfJobs := 1
		for j := 0; j < numberOfJobs; j++ {
			go func(j int) {
				jobs <- "channels.1a56060a-b034-4b65-a3f5-59d12a97d5da"
			}(j)
		}
		for {
			select {
				case resultMsg := <-result:
					fmt.Println(resultMsg)
			}
		}
		//// chờ nhận đủ kết quả
		//for c := 0; c < numberOfJobs; c++ {
		//	fmt.Println(<-result)
		//}
		//fmt.Println("DONE SESSION")
		//time.Sleep(24 * time.Hour)
	//}
}



