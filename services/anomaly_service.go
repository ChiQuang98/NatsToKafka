package services

import (
	"NatsToKafka/models"
	"NatsToKafka/workers"
)

func AnomalyChannel(anomaly_cfg *models.AnomalyChannel) (int, []byte) {
	if anomaly_cfg.Type == 1{
		workers.PushJobToChannel("channels."+anomaly_cfg.ChannelID)
	} else if anomaly_cfg.Type == 2{
		 go workers.PushKillSignalChannelKafka("channels."+anomaly_cfg.ChannelID)
		 go workers.PushKillSignalChannelNats("channels."+anomaly_cfg.ChannelID)
	}
	//fmt.Println(<-jobsChannel)
	return 200,[]byte("2")
}
