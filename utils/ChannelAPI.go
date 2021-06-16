package utils

import (
	"NatsToKafka/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetAllChannels() ([]models.Channel,error) {
	client := &http.Client{}
	var channels []models.Channel

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8888/api-gw/v1/thing/getall", nil)
	if err !=nil{
		return channels,err
	}
	token,err:=GetToken()
	if err!=nil{
		return nil,nil
	}
	req.Header.Set("Authorization",token)
	req.Header.Set("Accept","application/json")
	resTTTM, err := client.Do(req)
	if err!=nil{
		return channels,err
	}
	data, _ := ioutil.ReadAll(resTTTM.Body)
	err = json.Unmarshal(data,&channels)
	if err!=nil{
		return channels,err
	}
	return channels,nil
}
