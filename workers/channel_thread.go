package workers

import "NatsToKafka/models"
var Jobs chan string
var Result chan models.MessageNats
var KillsignalKafka chan string
var KillsignalNats chan string
func init() {
	// queue of jobs
	Jobs = make(chan string)

	// done channel lấy ra kết quả của jobs
	Result = make(chan models.MessageNats)
	// số lượng worker trong pool
	//vi` moi worker lam viec khong ket thuc, phai lang ng  he lien tuc nen so luong worker bang so luong channel
	KillsignalKafka = make(chan string)
	KillsignalNats = make(chan string)
}
func PushJobToChannel(job string)  {
	Jobs <- job
}
func PushResultNatsKafka(result models.MessageNats)  {
	Result <- result
}
func PushKillSignalChannelKafka(topicKill string)  {
	KillsignalKafka <- topicKill
}
func PushKillSignalChannelNats(topicKill string)  {
	KillsignalNats <- topicKill
}
func PollKillSignalChannelKafka() string {
	return <- KillsignalKafka
}
func PollKillSignalChannelNats() string {
	return <- KillsignalNats
}
func PollJobChannelNatsKafka() string  {
	return <- Jobs
}
func GetJobsChannel() chan string {
	return Jobs
}
func GetResultChannel() chan models.MessageNats {
	return Result
}
func GetKillSignalChannelKafka() chan string {
	return KillsignalKafka
}
func GetKillSignalChannelNats() chan string {
	return KillsignalNats
}